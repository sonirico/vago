package cqrs

import (
	"encoding/json"
	"time"

	"github.com/sonirico/vago/fp"

	"github.com/google/uuid"
)

//go:generate easyjson

type (
	CommandPayload struct {
		sendMsg
	}

	EventPayload struct {
		sendMsg
	}
)

//easyjson:json
type sendMsg struct {
	I           string            `json:"id"`
	AffinityKey *string           `json:"key"`
	V           string            `json:"version"`
	R           string            `json:"resource"`
	A           string            `json:"action"`
	T           time.Time         `json:"time"`
	P           any               `json:"payload"`
	UserID      fp.Option[string] `json:"user_id"`
}

func (h sendMsg) Version() string         { return h.V }
func (h sendMsg) Resource() string        { return h.R }
func (h sendMsg) Action() string          { return h.A }
func (h sendMsg) Payload() any            { return h.P }
func (h sendMsg) Key() *string            { return h.AffinityKey }
func (h sendMsg) ID() string              { return h.I }
func (h sendMsg) User() fp.Option[string] { return h.UserID }

func (h sendMsg) String() string {
	b, _ := json.Marshal(h)
	return string(b)
}

func NewSimpleCommand(
	version string,
	resource string,
	action string,
	payload any,
	key *string,
) CommandPayload {
	return NewCommand(version, resource, action, payload, key, time.Now().UTC())
}

func NewCommand(
	version string,
	resource string,
	action string,
	payload any,
	key *string,
	time time.Time,
) CommandPayload {
	return CommandPayload{
		sendMsg: sendMsg{
			V:           version,
			R:           resource,
			A:           action,
			P:           payload,
			AffinityKey: key,
			I:           uuid.NewString(),
			T:           time,
		},
	}
}

func NewUserCommand(
	userID string,
	version string,
	resource string,
	action string,
	payload any,
	key *string,
	time time.Time,
) CommandPayload {
	return CommandPayload{
		sendMsg: sendMsg{
			V:           version,
			R:           resource,
			A:           action,
			P:           payload,
			AffinityKey: key,
			I:           uuid.NewString(),
			T:           time,
			UserID:      fp.Some(userID),
		},
	}
}

func NewSimpleEvent(
	version string,
	resource string,
	action string,
	payload any,
	key *string,
) EventPayload {
	return NewEvent(version, resource, action, payload, key, time.Now().UTC())
}

func NewEvent(
	version string,
	resource string,
	action string,
	payload any,
	key *string,
	time time.Time,
) EventPayload {
	return EventPayload{
		sendMsg: sendMsg{
			V:           version,
			R:           resource,
			A:           action,
			P:           payload,
			AffinityKey: key,
			I:           uuid.NewString(),
			T:           time,
		},
	}
}

func NewUserEvent(
	userID string,
	version string,
	resource string,
	action string,
	payload any,
	key *string,
	time time.Time,
) EventPayload {
	return EventPayload{
		sendMsg: sendMsg{
			V:           version,
			R:           resource,
			A:           action,
			P:           payload,
			AffinityKey: key,
			I:           uuid.NewString(),
			T:           time,
			UserID:      fp.Some(userID),
		},
	}
}
