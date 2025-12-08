package cqrs

type (
	Encoder interface {
		Encode(any) ([]byte, error)
	}

	Decoder interface {
		Decode([]byte, any) error
	}

	Codec interface {
		Encoder
		Decoder
	}
)

func Decode[T any](decoder Decoder, data []byte) (T, error) {
	var x T
	err := decoder.Decode(data, &x)
	return x, err
}

func Encode[T any](encoder Encoder, x T) ([]byte, error) {
	return encoder.Encode(x)
}
