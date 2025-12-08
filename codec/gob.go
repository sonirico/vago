package codec

import (
	"bytes"
	"encoding/gob"
)

// GobCodec implements Codec using Go's encoding/gob.
type GobCodec struct{}

// Encode marshals a value to gob format.
func (c GobCodec) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode unmarshals gob data into a value.
func (c GobCodec) Decode(data []byte, v any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}

// NewGob creates a new GobCodec.
func NewGob() GobCodec {
	return GobCodec{}
}
