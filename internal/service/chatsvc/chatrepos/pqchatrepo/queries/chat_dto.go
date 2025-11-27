package queries

import "github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"

func (c GetChatsByUserRow) ToRepoChat() *chatrepos.Chat {
	return &chatrepos.Chat{
		ID: c.Chat.ID,
		CurrentUser: chatrepos.User{
			ID:        c.User.ID,
			ImageURL:  c.User.ImageUrl.String,
			FirstName: c.User.FirstName,
			LastName:  c.User.LastName.String,
			Username:  c.User.Username,
		},
		OtherUser: chatrepos.User{
			ID:        c.User_2.ID,
			ImageURL:  c.User_2.ImageUrl.String,
			FirstName: c.User_2.FirstName,
			LastName:  c.User_2.LastName.String,
			Username:  c.User_2.Username,
		},
	}
}
