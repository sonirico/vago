package cqrs

import (
	"context"
	"fmt"
)

type EventHandlerFunc func(ctx context.Context, event Event) error

type BaseEventHandler struct {
	version  string
	resource string
	action   string

	handler EventHandlerFunc
}

func (h BaseEventHandler) Version() string  { return h.version }
func (h BaseEventHandler) Resource() string { return h.resource }
func (h BaseEventHandler) Action() string   { return h.action }

func (h BaseEventHandler) Handle(ctx context.Context, event Event) error {
	if err := h.handler(ctx, event); err != nil {
		return fmt.Errorf("event Handler %s failed: %w", hashKey(h), err)
	}
	return nil
}

func NewEventHandler(v, r, a string, fn EventHandlerFunc) EventHandler {
	return &BaseEventHandler{
		version:  v,
		resource: r,
		action:   a,
		handler:  fn,
	}
}
