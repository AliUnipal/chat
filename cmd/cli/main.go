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
	firstName := userFlags.String("firstName", "", "First name (Required)")
	lastName := userFlags.String("lastName", "", "Last name")
	username := userFlags.String("username", "", "Username (Required)")
	imageURL := userFlags.String("imageUrl", "", "Image URL (Required)")

	err := userFlags.Parse(args)
	if err != nil {
		return userInput, err
	}

	if *firstName == "" {
		return userInput, errors.New("missing -firstName")
	}
	if *username == "" {
		return userInput, errors.New("missing -username")
	}
	if *imageURL == "" {
		return userInput, errors.New("missing -imageUrl")
	}

	userInput = usersvc.CreateUserInput{
		FirstName: *firstName,
		LastName:  *lastName,
		Username:  *username,
		ImageURL:  *imageURL,
	}

	return userInput, nil
}

func getGetUserInput(args []string) (uuid.UUID, error) {
	getUserFlags := flag.NewFlagSet("getUser", flag.ExitOnError)
	strID := getUserFlags.String("id", "", "User ID")

	err := getUserFlags.Parse(args)
	if err != nil {
		return uuid.Nil, err
	}

	if *strID == "" {
		return uuid.Nil, errors.New("missing -id")
	}

	id, err := uuid.Parse(*strID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func getCreateChatInput(args []string) ([2]uuid.UUID, error) {
	createChatFlags := flag.NewFlagSet("createChat", flag.ExitOnError)
	currID := createChatFlags.String("currentUserID", "", "Current User ID (Required)")
	otherID := createChatFlags.String("otherUserID", "", "Other User ID (Required)")

	err := createChatFlags.Parse(args)
	if err != nil {
		return [2]uuid.UUID{}, err
	}

	if *currID == "" {
		return [2]uuid.UUID{}, errors.New("missing -currentUserID")
	}
	if *otherID == "" {
		return [2]uuid.UUID{}, errors.New("missing -otherUserID")
	}
	pCurrID, err := uuid.Parse(*currID)
	if err != nil {
		return [2]uuid.UUID{}, err
	}
	pOtherID, err := uuid.Parse(*otherID)
	if err != nil {
		return [2]uuid.UUID{}, err
	}

	return [2]uuid.UUID{pCurrID, pOtherID}, nil
}

func getUserID(args []string) (uuid.UUID, error) {
	userIDFlags := flag.NewFlagSet("getUserID", flag.ExitOnError)
	strID := userIDFlags.String("id", "", "User ID (Required)")

	err := userIDFlags.Parse(args)
	if err != nil {
		return uuid.Nil, err
	}

	if *strID == "" {
		return uuid.Nil, errors.New("missing -id")
	}
	id, err := uuid.Parse(*strID)
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

	err := msgFlags.Parse(args)
	if err != nil {
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
		fmt.Println(`Commands:
- create-user
- get-user
- create-chat
- get-chats-by-user
- send-message
- get-messages`)
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
		ids, err := getCreateChatInput(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		currID, otherID := ids[0], ids[1]

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
