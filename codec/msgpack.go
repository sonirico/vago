package codec

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MsgpackCodec implements Codec using MessagePack serialization.
type MsgpackCodec struct{}

// Encode marshals a value to MessagePack format.
func (c MsgpackCodec) Encode(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decode unmarshals MessagePack data into a value.
func (c MsgpackCodec) Decode(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}

// NewMsgpack creates a new MsgpackCodec.
func NewMsgpack() MsgpackCodec {
	return MsgpackCodec{}
}
