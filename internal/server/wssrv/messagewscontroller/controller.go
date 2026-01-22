package messagewscontroller

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (uuid.UUID, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type messageWsController struct {
	messageSvc messageService
}

func New(messageSvc messageService) *messageWsController {
	return &messageWsController{messageSvc}
}

type (
	Testing struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	}

	CreateMessageRequest struct {
		SenderID    string              `json:"senderID"`
		ChatID      string              `json:"chatID"`
		Content     string              `json:"content"`
		ContentType message.ContentType `json:"contentType"`
	}
	validationErrors           map[string]string
	parsedCreateMessageRequest struct {
		senderID    uuid.UUID
		chatID      uuid.UUID
		content     []byte
		contentType message.ContentType
	}
	CreateMessageResponse struct {
		MessageID uuid.UUID `json:"messageID"`
	}
)

func (m *messageWsController) HandleMessages(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	ctx := context.Background()

	c, err := websocket.Accept(w, r, nil)

	if err != nil {
		slog.ErrorContext(ctx, "failed to accept websocket connection", "error", err)
		return
	}
	defer func() {
		_ = c.Close(websocket.StatusNormalClosure, "closing")
	}()

	c.SetReadLimit(1 << 20)

	for {
		var msg Testing
		if err := wsjson.Read(ctx, c, &msg); err != nil {
			if isNormalWSTermination(err) {
				return
			}
			slog.WarnContext(ctx, "failed to read websocket message", "error", err)
			return
		}

		var resp Testing
		switch msg.Type {
		case "ping":
			resp = Testing{Type: "pong", Text: msg.Text}
		default:
			// Generic echo
			resp = Testing{Type: "reply", Text: msg.Text}
		}

		writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = wsjson.Write(writeCtx, c, resp)
		cancel()
		if err != nil {
			if isNormalWSTermination(err) {
				return
			}
			slog.WarnContext(ctx, "failed to write websocket message", "error", err)
			return
		}
	}
}

func isNormalWSTermination(err error) bool {
	switch websocket.CloseStatus(err) {
	case websocket.StatusNormalClosure, websocket.StatusGoingAway:
		return true
	}

	if errors.Is(err, net.ErrClosed) {
		return true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}

	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}
