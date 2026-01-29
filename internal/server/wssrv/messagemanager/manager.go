package messagemanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

var (
	ReadBufferSize int64 = 1 << 20
)

type chatID = uuid.UUID
type clientsList map[chatID]map[*client]bool

func New(msgSvc messageService, userAcc userAccessor) *messageManager {
	m := &messageManager{
		msgSvc:   msgSvc,
		clients:  make(clientsList),
		userAcc:  userAcc,
		handlers: make(map[eventType]eventHandler),
	}

	m.setupEventHandlers()

	return m
}

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (message.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type userAccessor interface {
	CurrentUserID(ctx context.Context) (uuid.UUID, bool)
}

type messageManager struct {
	msgSvc  messageService
	clients clientsList
	userAcc userAccessor
	sync.RWMutex
	handlers map[eventType]eventHandler
}

func (m *messageManager) setupEventHandlers() {
	m.handlers[createMessageEventType] = m.createMessageHandler
	m.handlers[getMessagesEventType] = m.getMessagesHandler
}

func (m *messageManager) routeEvent(ctx context.Context, e event, c *client) error {
	if handler, ok := m.handlers[e.Type]; ok {
		return handler(ctx, e, c)
	} else {
		return errors.New("event handler not found")
	}
}

type (
	ServeMessageRequest struct {
		ChatID     string `http_path:"id"`
		parsedUUID uuid.UUID
	}
)

func (r *ServeMessageRequest) validate(ctx context.Context) error {
	if r.ChatID == "" {
		slog.ErrorContext(ctx, "chat id is required", "error", "required")
		return errors.New("chat id is required")
	}

	chatId, err := uuid.Parse(r.ChatID)
	if err != nil {
		slog.ErrorContext(ctx, "invalid uuid", "error", err)
		return errors.New("invalid chat id")
	}

	r.parsedUUID = chatId
	return nil
}

func (m *messageManager) ServeWS(w http.ResponseWriter, r *http.Request) {
	usrID, ok := m.userAcc.CurrentUserID(r.Context())
	slog.InfoContext(r.Context(), "user connected", "userID", usrID)
	if !ok {
		return
	}

	ctx := context.Background()
	req := ServeMessageRequest{
		ChatID: r.PathValue("id"),
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:3000", "127.0.0.1:3000"},
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to accept websocket connection", "error", err)
		return
	}

	if err := req.validate(ctx); err != nil {
		// TODO: respond with error and etc...
		_ = c.Write(ctx, websocket.MessageText, []byte(err.Error()))
		_ = c.Close(websocket.StatusBadGateway, "")
		return
	}

	c.SetReadLimit(ReadBufferSize)
	client := NewClient(c, m, usrID, req.parsedUUID)

	m.addClient(ctx, req.parsedUUID, client)
	go client.readEvents(ctx)
	go client.writeEvents(ctx)
}

func (m *messageManager) addClient(_ context.Context, chatID uuid.UUID, c *client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[chatID]; !ok {
		m.clients[chatID] = make(map[*client]bool)
	}
	m.clients[chatID][c] = true
}

func (m *messageManager) removeClient(ctx context.Context, chatID uuid.UUID, c *client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[chatID]; !ok {
		return
	}
	if _, ok := m.clients[chatID][c]; !ok {
		return
	}

	delete(m.clients[chatID], c)
	if len(m.clients[chatID]) == 0 {
		delete(m.clients, chatID)
	}
	err := c.conn.Close(websocket.StatusNormalClosure, "closing")
	slog.InfoContext(ctx, "closed websocket connection")
	if err == nil {
		return
	}

	status := websocket.CloseStatus(err)
	if status == -1 || status == websocket.StatusNormalClosure {
		return
	}

	slog.ErrorContext(ctx, "failed to close websocket connection", "error", err)
}
