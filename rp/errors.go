package rp

import "errors"

var (
	ErrConfig                 = errors.New("missing required config")
	ErrLoggerRequired         = errors.New("logger is required")
	ErrConsumerClosed         = errors.New("consumer is closed")
	ErrConsumerAlreadyCreated = errors.New("consumer is already created")
	ErrTopicsRequired         = errors.New("topics are required for consuming")
	ErrPingFailed             = errors.New("ping failed")
)
