package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/AliUnipal/chat/internal/models/chat"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/chatsvc"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo/inmemchatrepo"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
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

func getCreateUserInput(args []string) (usersvc.CreateUserInput, error) {
	var userInput usersvc.CreateUserInput

	userFlags := flag.NewFlagSet("createUser", flag.ExitOnError)
	userFlags.StringVar(&userInput.FirstName, "firstName", "", "First name (Required)")
	userFlags.StringVar(&userInput.LastName, "lastName", "potato", "Last name")
	userFlags.StringVar(&userInput.FirstName, "username", "", "Username (Required)")
	userFlags.StringVar(&userInput.FirstName, "imageUrl", "", "Image URL (Required)")

	if err := userFlags.Parse(args); err != nil {
		return userInput, err
	}

	if userInput.FirstName == "" {
		return userInput, errors.New("missing -firstName")
	}
	if userInput.Username == "" {
		return userInput, errors.New("missing -username")
	}
	if userInput.ImageURL == "" {
		return userInput, errors.New("missing -imageUrl")
	}

	return userInput, nil
}

func getGetUserInput(args []string) (uuid.UUID, error) {
	var strID string

	getUserFlags := flag.NewFlagSet("getUser", flag.ExitOnError)
	getUserFlags.StringVar(&strID, "id", "", "User ID")

	if err := getUserFlags.Parse(args); err != nil {
		return uuid.Nil, err
	}

	if strID == "" {
		return uuid.Nil, errors.New("missing -id")
	}

	id, err := uuid.Parse(strID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func getCreateChatInput(args []string) (uuid.UUID, uuid.UUID, error) {
	createChatFlags := flag.NewFlagSet("createChat", flag.ExitOnError)
	currID := createChatFlags.String("currentUserID", "", "Current User ID (Required)")
	otherID := createChatFlags.String("otherUserID", "", "Other User ID (Required)")

	if err := createChatFlags.Parse(args); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	if *currID == "" {
		return uuid.Nil, uuid.Nil, errors.New("missing -currentUserID")
	}
	if *otherID == "" {
		return uuid.Nil, uuid.Nil, errors.New("missing -otherUserID")
	}
	pCurrID, err := uuid.Parse(*currID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	pOtherID, err := uuid.Parse(*otherID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return pCurrID, pOtherID, nil
}

func getUserID(args []string) (uuid.UUID, error) {
	var strID string

	userIDFlags := flag.NewFlagSet("getUserID", flag.ExitOnError)
	userIDFlags.StringVar(&strID, "id", "", "User ID (Required)")

	if err := userIDFlags.Parse(args); err != nil {
		return uuid.Nil, err
	}

	if strID == "" {
		return uuid.Nil, errors.New("missing -id")
	}
	id, err := uuid.Parse(strID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func getSendMsgInput(args []string) (msgsvc.MessageInput, error) {
	var msgInput msgsvc.MessageInput

	msgFlags := flag.NewFlagSet("sendMsg", flag.ExitOnError)
	strSenderID := msgFlags.String("senderID", "", "Sender ID (Required)")
	strChatID := msgFlags.String("chatID", "", "ChatID (Required)")
	cont := msgFlags.String("content", "", "Content (Required)")

	if err := msgFlags.Parse(args); err != nil {
		return msgInput, err
	}

	if *strSenderID == "" {
		return msgInput, errors.New("missing -senderID")
	}
	if *strChatID == "" {
		return msgInput, errors.New("missing -chatID")
	}
	if *cont == "" {
		return msgInput, errors.New("missing -content")
	}

	senderID, err := uuid.Parse(*strSenderID)
	if err != nil {
		return msgInput, err
	}
	chatID, err := uuid.Parse(*strChatID)
	if err != nil {
		return msgInput, err
	}

	msgInput = msgsvc.MessageInput{
		SenderID:    senderID,
		ChatID:      chatID,
		Content:     []byte(*cont),
		ContentType: message.TextContentType,
	}

	return msgInput, nil
}

func getChatIDInput(args []string) (uuid.UUID, error) {
	chatFlags := flag.NewFlagSet("getChatID", flag.ExitOnError)
	strChatID := chatFlags.String("chatID", "", "ChatID (Required)")

	err := chatFlags.Parse(args)
	if err != nil {
		return uuid.Nil, err
	}

	if *strChatID == "" {
		return uuid.Nil, errors.New("missing -chatID")
	}
	chatID, err := uuid.Parse(*strChatID)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}

func main() {
	ctx := context.Background()

	userRepo := inmemuserrepo.New(ctx)
	userSvc := usersvc.NewService(userRepo)
	chatRepo := inmemchatrepo.New(ctx, userRepo)
	chatSvc := chatsvc.NewService(chatRepo)
	messageRepo := inmemmessagerepo.New(ctx, chatRepo)
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

	switch cmd {
	case "help":
		fmt.Println("Commands:")
		fmt.Println("- create-user")
		fmt.Println("- get-user")
		fmt.Println("- create-chat")
		fmt.Println("- get-chats-by-user")
		fmt.Println("- send-message")
		fmt.Println("- get-messages")
		os.Exit(0)
	case "create-user":
		userInput, err := getCreateUserInput(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		userID, err := app.userSvc.CreateUser(ctx, userInput)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		err = userRepo.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("User created:%+v\n", user.User{
			ID:        userID,
			ImageURL:  userInput.ImageURL,
			FirstName: userInput.FirstName,
			LastName:  userInput.LastName,
			Username:  userInput.Username,
		})
		os.Exit(0)
	case "get-user":
		id, err := getGetUserInput(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		u, err := app.userSvc.GetUser(ctx, id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"User Details:\n\nName: %s %s\nImage URL: %s\nUsername: %s\n",
			u.FirstName, u.LastName, u.ImageURL, u.Username,
		)
	case "create-chat":
		currID, otherID, err := getCreateChatInput(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		chatID, err := app.chatSvc.CreateChat(ctx, currID, otherID)
		if err != nil {
			log.Fatal(err)
		}

		err = chatRepo.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Created chat ID:", chatID)
		os.Exit(0)
	case "get-chats-by-user":
		userID, err := getUserID(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		chats, err := app.chatSvc.GetChats(ctx, userID)
		if err != nil {
			log.Fatal(err)
		}

		for i, c := range chats {
			fmt.Println("------")
			fmt.Println("Chat No.:", i)
			fmt.Printf("Chat ID: %s\n", c.ID)
			fmt.Printf("User one: %+v\n", c.CurrentUser)
			fmt.Printf("User two: %+v\n", c.OtherUser)
		}

		os.Exit(0)
	case "send-message":
		msg, err := getSendMsgInput(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		msgID, err := app.messageSvc.CreateMessage(ctx, msg)
		if err != nil {
			log.Fatal(err)
		}

		err = messageRepo.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Created message ID:", msgID)
		os.Exit(0)
	case "get-messages":
		chatID, err := getChatIDInput(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		msgs, err := app.messageSvc.GetMessages(ctx, chatID)
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
