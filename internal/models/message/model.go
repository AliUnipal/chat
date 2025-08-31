package message

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID          uuid.UUID
	SenderID    uuid.UUID
	ChatID      uuid.UUID
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
