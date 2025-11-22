package usercontroller

import (
	"context"

	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc"
	"github.com/google/uuid"
)

type userService interface {
	CreateUser(ctx context.Context, in usersvc.CreateUserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type userController struct {
	userService userService
}

func New(userService userService) *userController {
	return &userController{userService}
}
