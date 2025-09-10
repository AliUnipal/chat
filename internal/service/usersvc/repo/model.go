package repo

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
}

type CreateUserInput struct {
	ID        uuid.UUID
	ImageURL  string
	FirstName string
	LastName  string
	Username  string
}
