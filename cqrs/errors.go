package cqrs

import "errors"

var (
	ErrBusNotFound             = errors.New("bus not found")
	ErrHandleCommand           = errors.New("error handling command")
	ErrHandleEvent             = errors.New("error handling event")
	ErrSubscribeNonRecoverable = errors.New(
		"unrecoverable subscribe error",
	) // E.g, Client was closed, a new client should be spawned

	ErrPublish = errors.New("unable to publish msg")
)
