package repo

import "github.com/google/uuid"

type CreateUserInput struct {
	ID        uuid.UUID
	ImageURL  string
	FirstName string
	LastName  string
	Username  string
}
