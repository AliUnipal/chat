package main

import (
	"errors"
	"flag"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/AliUnipal/chat/internal/service/usersvc"
	"github.com/google/uuid"
)

func getCreateUserInput(args []string) (usersvc.CreateUserInput, error) {
	var userInput usersvc.CreateUserInput

	userFlags := flag.NewFlagSet("createUser", flag.ExitOnError)
	userFlags.StringVar(&userInput.FirstName, "firstName", "", "First name (Required)")
	userFlags.StringVar(&userInput.LastName, "lastName", "potato", "Last name")
	userFlags.StringVar(&userInput.Username, "username", "", "Username (Required)")
	userFlags.StringVar(&userInput.ImageURL, "imageUrl", "", "Image URL (Required)")

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
