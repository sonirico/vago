package codec

import (
	"encoding/json"
)

// JsonCodec implements Codec using the standard library's encoding/json.
type JsonCodec struct{}

// Encode marshals a value to JSON.
func (c JsonCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Decode unmarshals JSON data into a value.
func (c JsonCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// NewJson creates a new JsonCodec.
func NewJson() JsonCodec {
	return JsonCodec{}
}
