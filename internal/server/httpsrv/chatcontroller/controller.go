package chatcontroller

import (
	"context"
	"encoding/json"
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

type userAccessor interface {
	CurrentUserID(ctx context.Context) (uuid.UUID, bool)
}

type chatController struct {
	chatSvc chatService
	userAcc userAccessor
}

func New(chatSvc chatService, userAcc userAccessor) *chatController {
	return &chatController{chatSvc, userAcc}
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
		parsedUUID uuid.UUID
	}
	GetChatResponse struct {
		ID       string `json:"id"`
		UserID   string `json:"userID"`
		Name     string `json:"name"`
		ImageUrl string `json:"imageUrl"`
	}
	GetChatsResponse struct {
		Chats []GetChatResponse `json:"data"`
	}
)

func (c *chatController) GetChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usrId, ok := c.userAcc.CurrentUserID(ctx)
	if !ok {
		httpsrv.RespondWithAuthenticationError(ctx, w, errcodes.Unauthenticated)
		return
	}

	req := GetChatsRequest{
		parsedUUID: usrId,
	}

	chats, err := c.chatSvc.GetChats(r.Context(), req.parsedUUID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}
	resp := GetChatsResponse{
		Chats: []GetChatResponse{},
	}

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
