package cond

import (
	"fmt"
	"testing"
)

func TestIf(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		result := If(true, "yes", "no")
		if result != "yes" {
			t.Errorf("If(true, yes, no) = %s, want yes", result)
		}

		result = If(false, "yes", "no")
		if result != "no" {
			t.Errorf("If(false, yes, no) = %s, want no", result)
		}
	})

	t.Run("integers", func(t *testing.T) {
		result := If(true, 42, 0)
		if result != 42 {
			t.Errorf("If(true, 42, 0) = %d, want 42", result)
		}

		result = If(false, 42, 0)
		if result != 0 {
			t.Errorf("If(false, 42, 0) = %d, want 0", result)
		}
	})

	t.Run("floats", func(t *testing.T) {
		result := If(true, 3.14, 2.71)
		if result != 3.14 {
			t.Errorf("If(true, 3.14, 2.71) = %f, want 3.14", result)
		}
	})

	t.Run("booleans", func(t *testing.T) {
		result := If(true, true, false)
		if result != true {
			t.Errorf("If(true, true, false) = %v, want true", result)
		}
	})
}

func TestIfElse(t *testing.T) {
	result := IfElse(true, "first", "second")
	if result != "first" {
		t.Errorf("IfElse(true, first, second) = %s, want first", result)
	}

	result = IfElse(false, "first", "second")
	if result != "second" {
		t.Errorf("IfElse(false, first, second) = %s, want second", result)
	}
}

func TestIfFunc(t *testing.T) {
	t.Run("lazy evaluation true", func(t *testing.T) {
		called := false
		result := IfFunc(true,
			func() string { return "true branch" },
			func() string {
				called = true
				return "false branch"
			},
		)

		if result != "true branch" {
			t.Errorf("IfFunc(true, ...) = %s, want true branch", result)
		}
		if called {
			t.Error("false branch should not have been called")
		}
	})

	t.Run("lazy evaluation false", func(t *testing.T) {
		called := false
		result := IfFunc(false,
			func() string {
				called = true
				return "true branch"
			},
			func() string { return "false branch" },
		)

		if result != "false branch" {
			t.Errorf("IfFunc(false, ...) = %s, want false branch", result)
		}
		if called {
			t.Error("true branch should not have been called")
		}
	})
}

func TestIfPtr(t *testing.T) {
	valueA := "A"
	valueB := "B"

	t.Run("condition true", func(t *testing.T) {
		result := IfPtr(true, &valueA, &valueB)
		if result == nil || *result != "A" {
			t.Errorf("IfPtr(true, &A, &B) = %v, want &A", result)
		}
	})

	t.Run("condition false", func(t *testing.T) {
		result := IfPtr(false, &valueA, &valueB)
		if result == nil || *result != "B" {
			t.Errorf("IfPtr(false, &A, &B) = %v, want &B", result)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		result := IfPtr(true, nil, &valueB)
		if result != nil {
			t.Errorf("IfPtr(true, nil, &B) = %v, want nil", result)
		}
	})
}

func TestOrDefault(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		result := OrDefault(true, "value")
		if result != "value" {
			t.Errorf("OrDefault(true, value) = %s, want value", result)
		}

		result = OrDefault(false, "value")
		if result != "" {
			t.Errorf("OrDefault(false, value) = %s, want empty string", result)
		}
	})

	t.Run("integers", func(t *testing.T) {
		result := OrDefault(true, 42)
		if result != 42 {
			t.Errorf("OrDefault(true, 42) = %d, want 42", result)
		}

		result = OrDefault(false, 42)
		if result != 0 {
			t.Errorf("OrDefault(false, 42) = %d, want 0", result)
		}
	})
}

// Example demonstrates basic usage of If with strings
func ExampleIf() {
	result := If(true, "yes", "no")
	fmt.Println(result)
	// Output: yes
}

// Example demonstrates If with different condition
func ExampleIf_false() {
	result := If(false, "yes", "no")
	fmt.Println(result)
	// Output: no
}

// Example demonstrates If with integers
func ExampleIf_integers() {
	age := 25
	category := If(age >= 18, "adult", "minor")
	fmt.Println(category)
	// Output: adult
}

// Example demonstrates lazy evaluation with IfFunc
func ExampleIfFunc() {
	// This avoids expensive computation if not needed
	result := IfFunc(true,
		func() string { return "computed true" },
		func() string { return "expensive computation" },
	)
	fmt.Println(result)
	// Output: computed true
}

// Example demonstrates IfPtr with pointers
func ExampleIfPtr() {
	valueA := "option A"
	valueB := "option B"

	result := IfPtr(true, &valueA, &valueB)
	fmt.Println(*result)
	// Output: option A
}

// Example demonstrates OrDefault
func ExampleOrDefault() {
	// Get value if condition is true, otherwise zero value
	result := OrDefault(true, "has value")
	fmt.Println(result)

	result = OrDefault(false, "has value")
	fmt.Println(result)
	// Output:
	// has value
	//
}
