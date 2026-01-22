package message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID
	ChatID      uuid.UUID
	SenderID    uuid.UUID
	Content     []byte
	ContentType ContentType
	Timestamp   time.Time
}

type ContentType int

const (
	TextContentType ContentType = iota
	ImageContentType
	FileContentType
)
