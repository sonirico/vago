package codec

import (
	"testing"
)

type mockCodec struct {
	encodeFunc func(any) ([]byte, error)
	decodeFunc func([]byte, any) error
}

func (m mockCodec) Encode(v any) ([]byte, error) {
	return m.encodeFunc(v)
}

func (m mockCodec) Decode(data []byte, v any) error {
	return m.decodeFunc(data, v)
}

func TestEncode(t *testing.T) {
	codec := mockCodec{
		encodeFunc: func(v any) ([]byte, error) {
			return []byte("test"), nil
		},
	}

	result, err := Encode(codec, "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(result) != "test" {
		t.Errorf("expected test, got %s", string(result))
	}
}

func TestDecode(t *testing.T) {
	codec := mockCodec{
		decodeFunc: func(data []byte, v any) error {
			ptr := v.(*string)
			*ptr = "decoded"
			return nil
		},
	}

	result, err := Decode[string](codec, []byte("data"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "decoded" {
		t.Errorf("expected decoded, got %s", result)
	}
}
