package rp

import "time"

var apmTxType = "rp"

type (
	Msg struct {
		Topic     string
		Key       []byte
		Value     []byte
		Ts        time.Time
		Partition int32
	}
)
