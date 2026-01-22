package messagemanager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

var (
	ReadBufferSize int64 = 1 << 20
)

type chatID = uuid.UUID
type clientsList map[chatID]map[*client]bool

func New(msgSvc messageService) *messageManager {
	m := &messageManager{
		msgSvc:   msgSvc,
		clients:  make(clientsList),
		handlers: make(map[eventType]eventHandler),
	}

	m.setupEventHandlers()

	return m
}

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (message.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type messageManager struct {
	msgSvc  messageService
	clients clientsList
	sync.RWMutex
	handlers map[eventType]eventHandler
}

func (m *messageManager) setupEventHandlers() {
	m.handlers[createMessageEventType] = m.createMessageHandler
}

func (m *messageManager) routeEvent(ctx context.Context, e event, c *client) error {
	if handler, ok := m.handlers[e.Type]; ok {
		return handler(ctx, e, c)
	} else {
		return errors.New("event handler not found")
	}
}

type (
	ServeMessageRequest struct {
		ChatID     string `http_path:"id"`
		parsedUUID uuid.UUID
	}
)

func (r *ServeMessageRequest) validate(ctx context.Context) error {
	if r.ChatID == "" {
		slog.ErrorContext(ctx, "chat id is required", "error", "required")
		return errors.New("chat id is required")
	}

	chatId, err := uuid.Parse(r.ChatID)
	if err != nil {
		slog.ErrorContext(ctx, "invalid uuid", "error", err)
		return errors.New("invalid chat id")
	}

	r.parsedUUID = chatId
	return nil
}

func (m *messageManager) ServeWS(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := ServeMessageRequest{
		ChatID: r.PathValue("id"),
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:3000"},
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to accept websocket connection", "error", err)
		return
	}

	if err := req.validate(ctx); err != nil {
		// TODO: respond with error and etc...
		_ = c.Write(ctx, websocket.MessageText, []byte(err.Error()))
		_ = c.Close(websocket.StatusBadGateway, "")
		return
	}

	c.SetReadLimit(ReadBufferSize)
	client := NewClient(c, m, req.parsedUUID)

	m.addClient(ctx, req.parsedUUID, client)
	go client.readEvents(ctx)
	go client.writeEvents(ctx)
}

func (m *messageManager) addClient(_ context.Context, chatID uuid.UUID, c *client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[chatID]; !ok {
		m.clients[chatID] = make(map[*client]bool)
	}
	m.clients[chatID][c] = true
}

func (m *messageManager) removeClient(ctx context.Context, chatID uuid.UUID, c *client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[chatID]; !ok {
		return
	}
	if _, ok := m.clients[chatID][c]; !ok {
		return
	}

	delete(m.clients[chatID], c)
	if len(m.clients[chatID]) == 0 {
		delete(m.clients, chatID)
	}
	err := c.conn.Close(websocket.StatusNormalClosure, "closing")
	slog.InfoContext(ctx, "closed websocket connection")
	if err == nil {
		return
	}

	status := websocket.CloseStatus(err)
	if status == -1 || status == websocket.StatusNormalClosure {
		return
	}

	slog.ErrorContext(ctx, "failed to close websocket connection", "error", err)
}

// ======================= Handlers =======================

type (
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
		ID          uuid.UUID           `json:"id"`
		SenderID    uuid.UUID           `json:"senderID"`
		ChatID      uuid.UUID           `json:"chatID"`
		Content     []byte              `json:"content"`
		ContentType message.ContentType `json:"contentType"`
		Timestamp   time.Time           `json:"timestamp"`
	}
	CreateMessageEventResponse struct {
		Type eventType             `json:"type"`
		Data CreateMessageResponse `json:"data"`
	}
)

func (v validationErrors) Error() string {
	return "validation errors"
}

func (c *CreateMessageRequest) validate(ctx context.Context) (parsedCreateMessageRequest, error) {
	var p parsedCreateMessageRequest
	errs := validationErrors{}

	if c.SenderID == "" {
		errs["senderID"] = "required"
	}

	if c.ChatID == "" {
		errs["chatID"] = "required"
	}

	if len(c.Content) == 0 {
		errs["content"] = "required"
	}

	// TODO contentType & content validation

	senderId, err := uuid.Parse(c.SenderID)
	if err != nil {
		errs["senderID"] = "invalid uuid"
	}
	chatId, err := uuid.Parse(c.ChatID)
	if err != nil {
		errs["chatID"] = "invalid uuid"
	}
	content, err := base64.StdEncoding.DecodeString(c.Content)
	if err != nil {
		errs["content"] = "invalid base64"
	}

	if len(errs) > 0 {
		slog.ErrorContext(ctx, "validation errors", "errors", errs)
		return p, errs
	}

	p.senderID = senderId
	p.chatID = chatId
	p.content = content

	return p, nil
}

func (m *messageManager) createMessageHandler(ctx context.Context, e event, c *client) error {
	req := CreateMessageRequest{}
	if err := json.Unmarshal(e.Data, &req); err != nil {
		// TODO: Error handling
		return err
	}

	pReq, err := req.validate(ctx)
	if err != nil {
		return err
	}

	msg, err := m.msgSvc.CreateMessage(ctx, msgsvc.MessageInput{
		SenderID:    pReq.senderID,
		ChatID:      pReq.chatID,
		Content:     pReq.content,
		ContentType: pReq.contentType,
	})
	if err != nil {
		return err
	}

	if chat, ok := m.clients[pReq.chatID]; ok {
		broadMsg := CreateMessageResponse{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			ChatID:      msg.ChatID,
			Content:     msg.Content,
			ContentType: msg.ContentType,
			Timestamp:   msg.Timestamp,
		}
		data, err := json.Marshal(broadMsg)
		if err != nil {
			return err
		}

		out := event{
			Type: receiveMessageEventType,
			Data: data,
		}

		for client := range chat {
			client.egress <- out
		}
	} else {
		return errors.New("chat does not exist")
	}

	return nil
}
