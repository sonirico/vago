package str

import "testing"

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", true},
		{"non-empty string", "hello", false},
		{"whitespace", " ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.s); got != tt.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsSet(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", false},
		{"non-empty string", "hello", true},
		{"whitespace", " ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSet(tt.s); got != tt.want {
				t.Errorf("IsSet(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
