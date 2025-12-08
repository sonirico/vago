package cqrs

import (
	"context"
	"fmt"
)

// containerWrapper it's just a wrapper intended to capture an error to be able to handle it later, and thus allowing
// to provide both Commander and Eventer interfaces so that they do not return errors, which could stop consumers from
// consuming since and error returned by ConsumerHandler
type containerWrapper struct {
	*Container

	err error
}

func (c *containerWrapper) Command(ctx context.Context, busID string, cmd CommandPayload) {
	c.err = c.Container.Command(ctx, busID, cmd)
	if c.err != nil {
		c.err = fmt.Errorf("error emitting command: %w", c.err)
	}
}

func (c *containerWrapper) Event(ctx context.Context, busID string, e EventPayload) {
	c.err = c.Container.Event(ctx, busID, e)
	if c.err != nil {
		c.err = fmt.Errorf("error emitting event: %w", c.err)
	}
}
