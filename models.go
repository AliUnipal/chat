package main

import (
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	ImageURL    string
	FirstName   string
	LastName    string
	PhoneNumber string
}

type Chat struct {
	ID      uuid.UUID
	Type    ChatType
	Admin   uuid.UUID
	Name    string
	Image   string
	Members []uuid.UUID
}
type ChatType string

const (
	Direct ChatType = "DIRECT"
	Group  ChatType = "GROUP"
)

type Message struct {
	ID       uuid.UUID
	SenderID uuid.UUID
	ChatID   uuid.UUID
	Content  []byte
	ContentType
}
type ContentType int

const (
	String ContentType = iota
	Image  ContentType = iota
)
