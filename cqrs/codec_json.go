package cqrs

import (
	"encoding/json"
)

type JsonCodec struct{}

func (c JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c JsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func NewJson() JsonCodec {
	return JsonCodec{}
}
