// Package model provides the data structures for the chat application.
package models

import (
	"github.com/google/uuid"
)

// TODO: Separate each model into its respective folder.

type User struct {
	ID          uuid.UUID
	ImageURL    string
	FirstName   string
	LastName    string
	PhoneNumber string
}

type DirectChat struct {
	ID          uuid.UUID
	Admin       uuid.UUID
	Name        string
	ImageURL    string
	Participant User
}

type Message struct {
	ID          uuid.UUID
	SenderID    uuid.UUID
	ChatID      uuid.UUID
	Content     []byte
	ContentType ContentType
}

type ContentType int

const (
	TextContentType ContentType = iota
	ImageContentType
	FileContentType
)
