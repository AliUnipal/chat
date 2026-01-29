package messagemanager

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

var (
	writeTimeout = 10 * time.Second

	pongWait     = 15 * time.Second
	pingInterval = (pongWait * 9) / 10
)

func NewClient(conn *websocket.Conn, mng *messageManager, userID uuid.UUID, chatID uuid.UUID) *client {
	return &client{
		conn:   conn,
		mng:    mng,
		userID: userID,
		chatID: chatID,
		egress: make(chan event),
	}
}

type client struct {
	conn   *websocket.Conn
	mng    *messageManager
	userID uuid.UUID
	chatID uuid.UUID
	egress chan event
}

func (c *client) readEvents(ctx context.Context) {
	defer func() {
		c.mng.removeClient(ctx, c.chatID, c)
	}()

	for {
		_, data, err := c.conn.Reader(ctx)
		if err != nil {
			status := websocket.CloseStatus(err)
			if status == websocket.StatusAbnormalClosure || status == websocket.StatusGoingAway {
				slog.ErrorContext(ctx, "websocket connection closed", "error", err)
			}
			return
		}

		var req event
		if err := json.NewDecoder(data).Decode(&req); err != nil {
			slog.ErrorContext(ctx, "failed to decode websocket message", "error", err)
			return
		}

		if err := c.mng.routeEvent(ctx, req, c); err != nil {
			slog.ErrorContext(ctx, "failed to route event", "error", err)
		}
	}
}

func (c *client) writeEvents(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.mng.removeClient(ctx, c.chatID, c)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.heartbeat(ctx); err != nil {
				slog.ErrorContext(ctx, "failed to ping websocket connection", "error", err)
				return
			}
			slog.InfoContext(ctx, "pinged websocket connection")
		case msg, ok := <-c.egress:
			if !ok {
				if err := c.conn.Close(websocket.StatusNormalClosure, ""); err != nil {
					slog.ErrorContext(ctx, "failed to close websocket connection", "error", err)
				}
				return
			}

			if err := c.writeJSON(ctx, msg); err != nil {
				slog.ErrorContext(ctx, "failed to write websocket message", "error", err)
				return
			}
		}
	}
}

func (c *client) writeJSON(ctx context.Context, e event) error {
	wctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	w, err := c.conn.Writer(wctx, websocket.MessageText)
	if err != nil {
		return err
	}

	encErr := json.NewEncoder(w).Encode(e)
	if encErr != nil {
		_ = w.Close()
		return encErr
	}

	return w.Close()
}

// NOTE: Probably this is not working properly it might need further testing.
func (c *client) heartbeat(ctx context.Context) error {
	pctx, cancel := context.WithTimeout(ctx, pongWait)
	defer cancel()

	return c.conn.Ping(pctx)
}
