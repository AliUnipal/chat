package repo

import (
	"github.com/google/uuid"
	"time"
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
	Messages    []Message
}

type Message struct {
	ID        uuid.UUID
	SenderID  uuid.UUID
	Content   []byte
	Timestamp time.Time
}

type CreateChatInput struct {
	ID          uuid.UUID
	CurrentUser User
	OtherUser   User
}
