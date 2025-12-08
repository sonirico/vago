package rxconfig

import "time"

//go:generate easyjson

type (
	//easyjson:json
	pubSubMsg struct {
		Service string    `json:"service"`
		Node    string    `json:"node"`
		Channel string    `json:"channel"`
		Time    time.Time `json:"time"`
		// Payload
	}
)
