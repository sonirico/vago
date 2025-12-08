package codec

import (
	"testing"
)

func TestGobCodec_Encode(t *testing.T) {
	codec := NewGob()
	input := testStruct{Name: "gob", Value: 789}

	data, err := codec.Encode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestGobCodec_Decode(t *testing.T) {
	codec := NewGob()
	original := testStruct{Name: "gob", Value: 789}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	var decoded testStruct
	err = codec.Decode(data, &decoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("expected name %s, got %s", original.Name, decoded.Name)
	}
	if decoded.Value != original.Value {
		t.Errorf("expected value %d, got %d", original.Value, decoded.Value)
	}
}

func TestGobCodec_RoundTrip(t *testing.T) {
	codec := NewGob()
	original := testStruct{Name: "roundtrip", Value: 999}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	var decoded testStruct
	err = codec.Decode(data, &decoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if decoded.Name != original.Name || decoded.Value != original.Value {
		t.Errorf("roundtrip failed: expected %+v, got %+v", original, decoded)
	}
}
