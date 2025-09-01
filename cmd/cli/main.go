package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/AliUnipal/chat/internal/models/chat"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/chatsvc"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo/inmemchatrepo"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/AliUnipal/chat/internal/service/msgsvc/repo"
	"github.com/AliUnipal/chat/internal/service/msgsvc/repo/inmemmessagerepo"
	"github.com/AliUnipal/chat/internal/service/usersvc"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo/inmemuserrepo"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

type userService interface {
	CreateUser(ctx context.Context, in usersvc.CreateUserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type chatService interface {
	CreateChat(ctx context.Context, currentUserID, otherUserID uuid.UUID) (uuid.UUID, error)
	GetChats(ctx context.Context, userID uuid.UUID) ([]chat.Chat, error)
}

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (uuid.UUID, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type application struct {
	chatSvc    chatService
	userSvc    userService
	messageSvc messageService
}

func main() {
	userRepo := inmemuserrepo.New()
	userSvc := usersvc.NewService(userRepo)
	chatRepo := inmemchatrepo.New(userRepo)
	chatSvc := chatsvc.NewService(chatRepo)
	messageRepo := inmemmessagerepo.New(chatRepo, make(map[uuid.UUID][]repo.Message))
	msgSvc := msgsvc.NewService(messageRepo)

	app := &application{
		chatSvc:    chatSvc,
		userSvc:    userSvc,
		messageSvc: msgSvc,
	}

	if len(os.Args) < 2 {
		fmt.Println("arguments must be above one")
		os.Exit(1)
	}
	cmd := os.Args[1]
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	if cmd == "" {
		log.Fatal("invalid command")
	}
	reader := bufio.NewReader(os.Stdin)

	switch cmd {
	case "help":
		fmt.Println(`Commands:
- create-user
- get-user
- create-chat
- get-chats-by-user
- send-message
- get-messages`)
		os.Exit(0)
	case "create-user":
		var userInput usersvc.CreateUserInput
		v, err := read(reader, "first name", true)
		if err != nil {
			log.Fatal(err)
		}
		userInput.FirstName = v

		v, err = read(reader, "last name", false)
		if err != nil {
			log.Fatal(err)
		}
		userInput.LastName = v

		v, err = read(reader, "username", true)
		if err != nil {
			log.Fatal(err)
		}
		userInput.Username = v

		v, err = read(reader, "image URL", false)
		if err != nil {
			log.Fatal(err)
		}
		userInput.ImageURL = v

		userID, err := app.userSvc.CreateUser(context.Background(), userInput)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Println("User created:", userID)
		os.Exit(0)
	case "get-user":
		strID, err := read(reader, "user ID", true)
		if err != nil {
			log.Fatal(err)
		}
		id, err := uuid.Parse(strID)
		if err != nil {
			log.Fatal(err)
		}

		u, err := app.userSvc.GetUser(context.Background(), id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"User Details:\n Name: %s %s\nImage URL: %s\nUsername: %s\n",
			u.FirstName, u.LastName, u.ImageURL, u.Username,
		)
	case "create-chat":
		firstID, err := read(reader, "first user ID", true)
		if err != nil {
			log.Fatal(err)
		}
		fUserID, err := uuid.Parse(firstID)
		if err != nil {
			log.Fatal(err)
		}
		secondID, err := read(reader, "second user ID", true)
		if err != nil {
			log.Fatal(err)
		}
		oUserID, err := uuid.Parse(secondID)
		if err != nil {
			log.Fatal(err)
		}
		chatID, err := app.chatSvc.CreateChat(context.Background(), fUserID, oUserID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Created chat ID:", chatID)
		os.Exit(0)
	case "get-chats-by-user":
		in, err := read(reader, "user ID", true)
		if err != nil {
			log.Fatal(err)
		}
		userID, err := uuid.Parse(in)
		if err != nil {
			log.Fatal(err)
		}
		chats, err := app.chatSvc.GetChats(context.Background(), userID)
		if err != nil {
			log.Fatal(err)
		}

		for i, c := range chats {
			fmt.Println("------")
			fmt.Println("Chat No.: ", i)
			fmt.Printf("Chat ID: %s\n", c.ID)
			fmt.Printf("User one ID: %s\n", c.CurrentUser)
			fmt.Printf("User two ID: %s\n", c.OtherUser)
		}

		os.Exit(0)
	case "send-message":
		senderIDIn, err := read(reader, "sender ID", true)
		if err != nil {
			log.Fatal(err)
		}
		senderID, err := uuid.Parse(senderIDIn)
		if err != nil {
			log.Fatal(err)
		}
		chatIDIn, err := read(reader, "chat ID", true)
		if err != nil {
			log.Fatal(err)
		}
		chatID, err := uuid.Parse(chatIDIn)
		if err != nil {
			log.Fatal(err)
		}

		cont, err := read(reader, "content", false)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: will make it only text for now, later we do the rest.
		//fmt.Println("Content Types:")
		//contType, err := read(reader, "type", false)
		//if err != nil {
		//	log.Fatal(err)
		//}

		msgID, err := app.messageSvc.CreateMessage(context.Background(), msgsvc.MessageInput{
			SenderID:    senderID,
			ChatID:      chatID,
			Content:     []byte(cont),
			ContentType: 0,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Created message ID:", msgID)
		os.Exit(0)
	case "get-messages":
		chatIDIn, err := read(reader, "chat ID", true)
		if err != nil {
			log.Fatal(err)
		}
		chatID, err := uuid.Parse(chatIDIn)
		if err != nil {
			log.Fatal(err)
		}

		msgs, err := app.messageSvc.GetMessages(context.Background(), chatID)
		if err != nil {
			log.Fatal(err)
		}
		for i, m := range msgs {
			fmt.Println("------")
			fmt.Println("Message No.:", i)
			fmt.Printf("Message ID: %s\n", m.ID)
			fmt.Printf("SenderID: %s\n", m.SenderID)
			fmt.Printf("Content: %s\n", string(m.Content))
			fmt.Printf("Content Type: %v\n", m.ContentType)
			fmt.Printf("Timestamp: %v\n", m.Timestamp)
		}
		os.Exit(0)
	default:
		log.Fatal("unknown command")
	}

}

func read(reader *bufio.Reader, fieldName string, required bool) (string, error) {
	fmt.Printf("Enter your %s:\n", fieldName)
	in, err := reader.ReadString('\n')
	in = strings.TrimSpace(in)
	if err != nil {
		return "", err
	}
	if in == "" && required {
		return "", fmt.Errorf("missing required field: %s", fieldName)
	}

	return in, nil
}
