package codec

// Encoder is the interface for encoding data to bytes.
type Encoder interface {
	Encode(any) ([]byte, error)
}

// Decoder is the interface for decoding bytes to data.
type Decoder interface {
	Decode([]byte, any) error
}

// Codec is the interface that combines Encoder and Decoder.
type Codec interface {
	Encoder
	Decoder
}

// Decode is a generic helper function that decodes bytes into a value of type T.
func Decode[T any](decoder Decoder, data []byte) (T, error) {
	var x T
	err := decoder.Decode(data, &x)
	return x, err
}

// Encode is a generic helper function that encodes a value of type T into bytes.
func Encode[T any](encoder Encoder, x T) ([]byte, error) {
	return encoder.Encode(x)
}
