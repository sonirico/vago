package cqrs

import (
	"encoding/json"
	"time"

	"github.com/sonirico/vago/fp"
)

//go:generate easyjson

type (
	Command struct {
		recvMsg
	}

	Event struct {
		recvMsg
	}
)

const (
	headerUserID = "user_id"
)

//easyjson:json
type recvMsg struct {
	I      string           `json:"id"`
	V      string           `json:"version"`
	R      string           `json:"resource"`
	A      string           `json:"action"`
	T      time.Time        `json:"time"`
	P      json.RawMessage  `json:"payload"`
	UserID fp.Option[string] `json:"user_id"`

	recordKey       []byte
	recordPartition int32
	recordTs        time.Time
}

func (m recvMsg) User() fp.Option[string] {
	return m.UserID
}

func (m recvMsg) Version() string  { return m.V }
func (m recvMsg) Resource() string { return m.R }
func (m recvMsg) Action() string   { return m.A }
func (m recvMsg) Payload() []byte  { return m.P }
func (m recvMsg) ID() string       { return m.I }

func (m recvMsg) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func (m recvMsg) Command() Command {
	return Command{recvMsg: m}
}

func (m recvMsg) Event() Event {
	return Event{recvMsg: m}
}

func newRecvMsg(
	id string,
	version string,
	resource string,
	action string,
	time time.Time,
	payload json.RawMessage,
	userID fp.Option[string],
	recordKey []byte,
	recordPartition int32,
	recordTs time.Time,
) recvMsg {
	return recvMsg{
		I:               id,
		V:               version,
		R:               resource,
		A:               action,
		T:               time,
		P:               payload,
		UserID:          userID,
		recordKey:       recordKey,
		recordPartition: recordPartition,
		recordTs:        recordTs,
	}
}

func NewRecvEvent(
	id string,
	version string,
	resource string,
	action string,
	time time.Time,
	payload json.RawMessage,
	userID fp.Option[string],
	recordKey []byte,
	recordPartition int32,
	recordTs time.Time,
) Event {
	return Event{
		recvMsg: newRecvMsg(
			id,
			version,
			resource,
			action,
			time,
			payload,
			userID,
			recordKey,
			recordPartition,
			recordTs,
		),
	}
}

func NewRecvCommand(
	id string,
	version string,
	resource string,
	action string,
	time time.Time,
	payload json.RawMessage,
	userID fp.Option[string],
	recordKey []byte,
	recordPartition int32,
	recordTs time.Time,
) Command {
	return Command{
		recvMsg: newRecvMsg(
			id,
			version,
			resource,
			action,
			time,
			payload,
			userID,
			recordKey,
			recordPartition,
			recordTs,
		),
	}
}
