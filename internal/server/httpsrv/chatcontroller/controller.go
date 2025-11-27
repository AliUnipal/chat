package chatcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AliUnipal/chat/internal/models/chat"
	"github.com/AliUnipal/chat/internal/server"
	"github.com/AliUnipal/chat/internal/server/httpsrv"
	"github.com/AliUnipal/chat/pkg/errcodes"
	"github.com/google/uuid"
)

type chatService interface {
	CreateChat(ctx context.Context, currentUserID, otherUserID uuid.UUID) (uuid.UUID, error)
	GetChats(ctx context.Context, userID uuid.UUID) ([]chat.Chat, error)
}

type chatController struct {
	chatSvc chatService
}

func New(chatSvc chatService) *chatController {
	return &chatController{chatSvc}
}

type (
	CreateChatRequest struct {
		CurrentUserID uuid.UUID `json:"currentUserID"`
		OtherUserID   uuid.UUID `json:"otherUserID"`
	}
	CreateChatResponse struct {
		ChatID uuid.UUID `json:"chatID"`
	}
)

func (c *chatController) CreateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var in CreateChatRequest

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		slog.ErrorContext(ctx, "failed to decode request CreateChatRequest", "error", err)
		httpsrv.RespondWithError(ctx, w, "invalid_body", server.ErrInvalidRequestBody)
		return
	}

	id, err := c.chatSvc.CreateChat(ctx, in.CurrentUserID, in.OtherUserID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, CreateChatResponse{id})
}

type (
	GetChatsRequest struct {
		UserID     string `http_query:"userID"`
		parsedUUID uuid.UUID
	}
	GetChatResponse struct {
		ID       string `json:"id"`
		UserID   string `json:"userID"`
		Name     string `json:"name"`
		ImageUrl string `json:"imageUrl"`
	}
	GetChatsResponse struct {
		Chats []GetChatResponse `json:"chats"`
	}
)

func (r *GetChatsRequest) validate(ctx context.Context) error {
	if r.UserID == "" {
		slog.ErrorContext(ctx, "id is required", "error")
		return errors.New("id is required")
	}

	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "invalid uuid", "error", err)
		return server.ErrInvalidUUID
	}

	r.parsedUUID = userID
	return nil
}

func (c *chatController) GetChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetChatsRequest{
		UserID: r.URL.Query().Get("userID"),
	}

	if err := req.validate(ctx); err != nil {
		httpsrv.RespondWithBadRequestError(ctx, w, errcodes.InvalidUUID, err)
		return
	}

	chats, err := c.chatSvc.GetChats(r.Context(), req.parsedUUID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}
	var resp GetChatsResponse

	for _, c := range chats {
		resp.Chats = append(resp.Chats, GetChatResponse{
			ID:       c.ID.String(),
			UserID:   c.OtherUser.ID.String(),
			Name:     c.OtherUser.FirstName,
			ImageUrl: c.OtherUser.ImageURL,
		})
	}

	httpsrv.RespondWithJSON(ctx, w, resp)
}
