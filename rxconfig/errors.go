package rxconfig

import "errors"

var (
	ErrDoesNotExist = errors.New("does not exist")
)

var (
	errPubSubClosed = errors.New("pubsub closed")
)
