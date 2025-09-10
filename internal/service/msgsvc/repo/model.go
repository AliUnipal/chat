package repo

import (
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID          uuid.UUID           `json:"id"`
	SenderID    uuid.UUID           `json:"sender_id"`
	ChatID      uuid.UUID           `json:"chat_id"`
	Content     []byte              `json:"content"`
	ContentType message.ContentType `json:"content_type"`
	Timestamp   time.Time           `json:"timestamp"`
}

type CreateMessageInput struct {
	ID          uuid.UUID
	SenderID    uuid.UUID
	ChatID      uuid.UUID
	Content     []byte
	ContentType message.ContentType
	Timestamp   time.Time
}
