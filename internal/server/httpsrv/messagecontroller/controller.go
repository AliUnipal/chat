package messagecontroller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/server/httpsrv"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/AliUnipal/chat/pkg/errcodes"
	"github.com/google/uuid"
)

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (message.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type messageController struct {
	messageSvc messageService
}

func New(messageSvc messageService) *messageController {
	return &messageController{messageSvc}
}

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

func (m *messageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "failed to decode request CreateMessageRequest", "error", err)
		httpsrv.RespondWithError(ctx, w, errcodes.InvalidInput, err)
		return
	}

	pReq, err := req.validate(ctx)
	if err != nil {
		if fieldErrs, ok := err.(validationErrors); ok {
			var kvs httpsrv.KVs

			for field, errMsg := range fieldErrs {
				kvs = append(kvs, httpsrv.KV(field, errMsg))
			}

			httpsrv.RespondWithValidationError(ctx, w, errcodes.InvalidInput, kvs...)
			return
		}

		httpsrv.RespondWithError(ctx, w, errcodes.InvalidInput, err)
		return
	}

	msg, err := m.messageSvc.CreateMessage(ctx, msgsvc.MessageInput{
		SenderID:    pReq.senderID,
		ChatID:      pReq.chatID,
		Content:     pReq.content,
		ContentType: pReq.contentType,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create message", "error", err)
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, CreateMessageResponse{
		ID:          msg.ID,
		SenderID:    msg.SenderID,
		ChatID:      msg.ChatID,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		Timestamp:   msg.Timestamp,
	})
}

type (
	GetMessagesRequest struct {
		ChatID     string `http_path:"id"`
		parsedUUID uuid.UUID
	}
	GetMessageResponse struct {
		ID        uuid.UUID `json:"id"`
		SenderID  uuid.UUID `json:"senderID"`
		ChatID    uuid.UUID `json:"chatID"`
		Content   string    `json:"content"`
		Timestamp int64     `json:"timestamp"`
	}
	GetMessagesResponse struct {
		Messages []GetMessageResponse `json:"data"`
	}
)

func (r *GetMessagesRequest) validate(ctx context.Context) error {
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

func (m *messageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetMessagesRequest{
		ChatID: r.PathValue("id"),
	}

	if err := req.validate(ctx); err != nil {
		httpsrv.RespondWithBadRequestError(ctx, w, errcodes.InvalidUUID, err)
		return
	}

	messages, err := m.messageSvc.GetMessages(ctx, req.parsedUUID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	resp := GetMessagesResponse{
		Messages: []GetMessageResponse{},
	}

	for _, m := range messages {
		resp.Messages = append(resp.Messages, GetMessageResponse{
			ID:        m.ID,
			SenderID:  m.SenderID,
			ChatID:    m.ChatID,
			Content:   string(m.Content),
			Timestamp: m.Timestamp.Unix(),
		})
	}

	httpsrv.RespondWithJSON(ctx, w, resp)
}
