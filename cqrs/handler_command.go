package cqrs

import (
	"context"
	"fmt"
)

type CommandHandlerFunc func(ctx context.Context, cmd Command, eventer Eventer) error

type BaseCommandHandler struct {
	version  string
	resource string
	action   string

	handler CommandHandlerFunc
}

func (h BaseCommandHandler) Version() string  { return h.version }
func (h BaseCommandHandler) Resource() string { return h.resource }
func (h BaseCommandHandler) Action() string   { return h.action }

func (h BaseCommandHandler) Handle(ctx context.Context, cmd Command, eventer Eventer) error {
	if err := h.handler(ctx, cmd, eventer); err != nil {
		return fmt.Errorf("event Handler %s failed: %w", hashKey(h), err)
	}
	return nil
}

func NewCommandHandler(
	v, r, a string,
	fn CommandHandlerFunc,
) CommandHandler {
	return &BaseCommandHandler{
		version:  v,
		resource: r,
		action:   a,
		handler:  fn,
	}
}
