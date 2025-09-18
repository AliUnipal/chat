package repo

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
}

type Chat struct {
	ID          uuid.UUID `json:"id"`
	CurrentUser User      `json:"current_user"`
	OtherUser   User      `json:"other_user"`
}

type CreateChatInput struct {
	ID, CurrentUserID, OtherUserID uuid.UUID
}
