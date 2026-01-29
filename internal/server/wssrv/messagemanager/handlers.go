package messagemanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/google/uuid"
)

type (
	CreateMessageRequest struct {
		SenderID    string              `json:"senderID"`
		ChatID      string              `json:"chatID"`
		Content     string              `json:"content"`
		ContentType message.ContentType `json:"contentType"`
	}
	validationErrors           map[string]string
	parsedCreateMessageRequest struct {
		senderID    uuid.UUID
		chatID      uuid.UUID
		content     []byte
		contentType message.ContentType
	}
	CreateMessageResponse struct {
		ID          uuid.UUID           `json:"id"`
		SenderID    uuid.UUID           `json:"senderID"`
		ChatID      uuid.UUID           `json:"chatID"`
		Content     []byte              `json:"content"`
		ContentType message.ContentType `json:"contentType"`
		Timestamp   time.Time           `json:"timestamp"`
	}
	CreateMessageEventResponse struct {
		Type eventType             `json:"type"`
		Data CreateMessageResponse `json:"data"`
	}
)

func (v validationErrors) Error() string {
	return "validation errors"
}

func (c *CreateMessageRequest) validate(ctx context.Context) (parsedCreateMessageRequest, error) {
	fmt.Printf("Content raw: %q\n", c.Content)
	fmt.Printf("Content bytes: %#v\n", []byte(c.Content))

	var p parsedCreateMessageRequest
	errs := validationErrors{}

	if c.SenderID == "" {
		errs["senderID"] = "required"
	}

	if c.ChatID == "" {
		errs["chatID"] = "required"
	}

	if len(c.Content) == 0 {
		errs["content"] = "required"
	}

	// TODO contentType & content validation

	senderId, err := uuid.Parse(c.SenderID)
	if err != nil {
		errs["senderID"] = "invalid uuid"
	}
	chatId, err := uuid.Parse(c.ChatID)
	if err != nil {
		errs["chatID"] = "invalid uuid"
	}
	//content, err := base64.StdEncoding.DecodeString(c.Content)
	//if err != nil {
	//	errs["content"] = "invalid base64"
	//}
	//fmt.Printf("Content decoded: %v\n", content)
	if len(errs) > 0 {
		slog.ErrorContext(ctx, "validation errors", "errors", errs)
		return p, errs
	}

	p.senderID = senderId
	p.chatID = chatId
	p.content = []byte(c.Content)

	return p, nil
}

func (m *messageManager) createMessageHandler(ctx context.Context, e event, c *client) error {
	req := CreateMessageRequest{}
	if err := json.Unmarshal(e.Data, &req); err != nil {
		// TODO: Error handling
		return err
	}

	pReq, err := req.validate(ctx)
	if err != nil {
		return err
	}

	msg, err := m.msgSvc.CreateMessage(ctx, msgsvc.MessageInput{
		SenderID:    pReq.senderID,
		ChatID:      pReq.chatID,
		Content:     pReq.content,
		ContentType: pReq.contentType,
	})
	if err != nil {
		return err
	}

	if chat, ok := m.clients[pReq.chatID]; ok {
		broadMsg := CreateMessageResponse{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			ChatID:      msg.ChatID,
			Content:     msg.Content,
			ContentType: msg.ContentType,
			Timestamp:   msg.Timestamp,
		}
		data, err := json.Marshal(broadMsg)
		if err != nil {
			return err
		}

		out := event{
			Type: receiveMessageEventType,
			Data: data,
		}

		for client := range chat {
			client.egress <- out
		}
	} else {
		return errors.New("chat does not exist")
	}

	return nil
}

type (
	// NOTE: will be through the client chatID
	// TODO: add load more, or length handling etc...
	GetMessagesRequest struct {
	}
	GetMessageResponse struct {
		ID          uuid.UUID           `json:"id"`
		SenderID    uuid.UUID           `json:"senderID"`
		ChatID      uuid.UUID           `json:"chatID"`
		Content     []byte              `json:"content"`
		ContentType message.ContentType `json:"contentType"`
		Timestamp   time.Time           `json:"timestamp"`
	}
	GetMessagesEventResponse struct {
		Type eventType            `json:"type"`
		Data []GetMessageResponse `json:"data"`
	}
)

// NOTE: I think making this function as a loader to get the messages once the user connects would be better.
//
//	however, I think it will be a alright to keep it like this to make the client decide whether it's a new empty
//	chat or loaded!!
func (m *messageManager) getMessagesHandler(ctx context.Context, e event, c *client) error {
	messages, err := m.msgSvc.GetMessages(ctx, c.chatID)
	if err != nil {
		return err
	}

	outMsgs := make([]GetMessageResponse, 0)

	for _, msg := range messages {
		outMsgs = append(outMsgs, GetMessageResponse{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			ChatID:      msg.ChatID,
			Content:     msg.Content,
			ContentType: msg.ContentType,
			Timestamp:   msg.Timestamp,
		})
	}

	data, err := json.Marshal(outMsgs)
	if err != nil {
		return err
	}

	out := event{
		Type: getMessagesEventType,
		Data: data,
	}
	c.egress <- out

	return nil
}
