package cqrs

import "context"

type (
	Eventer interface {
		Event(ctx context.Context, bus string, e EventPayload)
	}

	Commander interface {
		Command(ctx context.Context, bus string, cmd CommandPayload)
	}
)

type (
	ContainerCommander interface {
		Command(ctx context.Context, busID string, cmd CommandPayload) error
	}

	ContainerEventer interface {
		Event(ctx context.Context, busID string, event EventPayload) error
	}

	ContainerOperator interface {
		ContainerCommander
		ContainerEventer
	}
)

type (
	Handler interface {
		Version() string
		Resource() string
		Action() string
	}

	EventHandler interface {
		Handler

		// Handle handles events. If errors are returned, message consumption is stopped.
		Handle(ctx context.Context, e Event) error
	}

	CommandHandler interface {
		Handler

		// Handle handles commands. If errors are returned, message consumption is stopped.
		Handle(ctx context.Context, cmd Command, eventer Eventer) error
	}

	SagaHandler interface {
		Handler

		// Handle handles events. If errors are returned, message consumption is stopped.
		Handle(ctx context.Context, e Event, commander Commander) error
	}
)
