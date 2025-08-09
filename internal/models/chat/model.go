package chat

import (
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/google/uuid"
)

type Chat struct {
	ID          uuid.UUID
	CurrentUser user.User
	OtherUser   user.User
	Messages    []message.Message
}
