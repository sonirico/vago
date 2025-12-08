package ptr

import (
	"fmt"
	"testing"
)

func TestPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		value := "hello"
		p := Ptr(value)

		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *p != value {
			t.Errorf("expected %q, got %q", value, *p)
		}

		// Verify it's a different address
		value = "world"
		if *p == value {
			t.Errorf("pointer should not be affected by original variable change")
		}
	})

	t.Run("integer", func(t *testing.T) {
		value := 42
		p := Ptr(value)

		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *p != value {
			t.Errorf("expected %d, got %d", value, *p)
		}
	})

	t.Run("literal", func(t *testing.T) {
		// Useful for getting pointer to literal
		p := Ptr("literal")

		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *p != "literal" {
			t.Errorf("expected %q, got %q", "literal", *p)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		value := person{name: "Alice", age: 30}
		p := Ptr(value)

		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if p.name != value.name || p.age != value.age {
			t.Errorf("expected %+v, got %+v", value, *p)
		}
	})

	t.Run("zero_value", func(t *testing.T) {
		p := Ptr(0)

		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *p != 0 {
			t.Errorf("expected 0, got %d", *p)
		}
	})
}

// ExamplePtr demonstrates creating a pointer to a value.
func ExamplePtr() {
	// Get pointer to a literal
	name := Ptr("Alice")
	age := Ptr(30)

	fmt.Println(*name, *age)
	// Output:
	// Alice 30
}

// ExamplePtr_struct demonstrates creating a pointer to a struct.
func ExamplePtr_struct() {
	type Config struct {
		Host string
		Port int
	}

	// Get pointer to struct literal
	config := Ptr(Config{
		Host: "localhost",
		Port: 8080,
	})

	fmt.Println(config.Host, config.Port)
	// Output:
	// localhost 8080
}
