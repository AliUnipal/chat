package messagemanager

import (
	"context"
	"encoding/json"
)

type event struct {
	Type eventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type eventHandler func(ctx context.Context, e event, c *client) error

type eventType string

const (
	createMessageEventType  eventType = "createMessage"
	receiveMessageEventType           = "receiveMessage"
)

type messageCreatedEvent struct{}
