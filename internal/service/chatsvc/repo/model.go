package repo

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	ImageURL  string
	FirstName string
	LastName  string
	Username  string
}

type Chat struct {
	ID          uuid.UUID
	CurrentUser User
	OtherUser   User
}

type CreateChatInput struct {
	ID, CurrentUserID, OtherUserID uuid.UUID
}
