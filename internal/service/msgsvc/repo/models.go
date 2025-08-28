package repo

import (
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID          uuid.UUID
	SenderID    uuid.UUID
	ChatID      uuid.UUID
	Content     []byte
	ContentType message.ContentType
	Timestamp   time.Time
}

type CreateMessageInput struct {
	ID          uuid.UUID
	SenderID    uuid.UUID
	ChatID      uuid.UUID
	Content     []byte
	ContentType message.ContentType
	Timestamp   time.Time
}
