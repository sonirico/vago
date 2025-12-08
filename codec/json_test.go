package codec

import (
	"testing"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestJsonCodec_Encode(t *testing.T) {
	codec := NewJson()
	input := testStruct{Name: "test", Value: 42}

	data, err := codec.Encode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"name":"test","value":42}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestJsonCodec_Decode(t *testing.T) {
	codec := NewJson()
	data := []byte(`{"name":"test","value":42}`)

	var output testStruct
	err := codec.Decode(data, &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Name != "test" {
		t.Errorf("expected name test, got %s", output.Name)
	}
	if output.Value != 42 {
		t.Errorf("expected value 42, got %d", output.Value)
	}
}

func TestJsonCodec_RoundTrip(t *testing.T) {
	codec := NewJson()
	original := testStruct{Name: "roundtrip", Value: 99}

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
