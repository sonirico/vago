package cqrs

import (
	"context"
	"fmt"
)

type SagaHandlerFunc func(ctx context.Context, event Event, commander Commander) error

type SagaHandlerOpts struct {
	CommitOnError bool
}

type BaseSagaHandler struct {
	version  string
	resource string
	action   string

	handler SagaHandlerFunc
	opts    SagaHandlerOpts
}

func (h BaseSagaHandler) Version() string  { return h.version }
func (h BaseSagaHandler) Resource() string { return h.resource }
func (h BaseSagaHandler) Action() string   { return h.action }

func (h BaseSagaHandler) Handle(ctx context.Context, event Event, commander Commander) error {
	if err := h.handler(ctx, event, commander); err != nil {
		return fmt.Errorf("saga Handler %s failed: %w", hashKey(h), err)
	}
	return nil
}

func NewSagaHandler(v, r, a string, fn SagaHandlerFunc, opts *SagaHandlerOpts) SagaHandler {
	if opts == nil {
		opts = &SagaHandlerOpts{}
	}

	return &BaseSagaHandler{
		version:  v,
		resource: r,
		action:   a,
		handler:  fn,
		opts:     *opts,
	}
}
