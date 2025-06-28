package fp

import (
	"errors"
	"fmt"
	"testing"
)

func TestResult(t *testing.T) {
	ok := Ok(1)
	if !ok.IsOk() {
		t.Errorf("unexpected IsOk result, want ok, have err: %s", ok.err)
	}
	if ok.IsErr() {
		t.Errorf("unexpected IsErr result, want no error, have err: %s", ok.err)
	}

	if !OkAny.IsOk() {
		t.Errorf("unexpected IsOk result, want ok, have err: %s", OkAny.err)
	}
	if OkAny.IsErr() {
		t.Errorf("unexpected IsErr result, want no error, have err: %s", OkAny.err)
	}

	fail := Err[struct{}](errors.New("cannot divide by zero"))

	if fail.IsOk() {
		t.Errorf("unexpected IsOk result, want err, have ok: %v", fail.value)
	}
	if !fail.IsErr() {
		t.Errorf("unexpected IsErr result, want error, have ok: %v", fail.value)
	}

	value, err := ok.Unwrap()
	if err != nil {
		t.Errorf("unexpected Unwrap result, want ok, have err: %s", err)
	}
	if value != 1 {
		t.Errorf("unexpected Unwrap value, want 1, have %d", value)
	}

	// UnwrapUnsafe
	_ = ok.UnwrapUnsafe()
}

func TestResult_Or(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.Or(Ok(2)).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = fail.Or(Ok(1)).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result on Err, want 1, have %d", value)
	}
}

func TestResult_OrElse(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.OrElse(func() Result[int] {
		return Ok(2)
	}).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = fail.OrElse(func() Result[int] {
		return Ok(1)
	}).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result on Err, want 1, have %d", value)
	}
}

func TestResult_UnwrapOr(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.UnwrapOr(3)

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = fail.UnwrapOr(1)

	if value != 1 {
		t.Errorf("unexpected result on Err, want 1, have %d", value)
	}
}

func TestResult_UnwrapOrElse(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.UnwrapOrElse(func() int { return 3 })

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = fail.UnwrapOrElse(func() int { return 1 })

	if value != 1 {
		t.Errorf("unexpected result on Err, want 1, have %d", value)
	}
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.UnwrapOrDefault()

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = fail.UnwrapOrDefault()

	if value != 0 {
		t.Errorf("unexpected result on Err, want 0, have %d", value)
	}
}

func TestResult_Match(t *testing.T) {
	ok := Ok(2)
	fail := Err[int](errors.New("cannot divide by zero"))

	value, err := ok.Match(
		func(x int) Result[int] {
			return Ok(x * x)
		},
		func(err error) Result[int] {
			return Ok(0)
		}).Unwrap()

	if err != nil {
		t.Errorf("unexpected error, want nil, have %s", err.Error())
	}

	if value != 4 {
		t.Errorf("unexpected result , want 4, have %d", value)
	}

	value, err = fail.Match(
		func(x int) Result[int] {
			return Ok(x * x)
		},
		func(err error) Result[int] {
			return Ok(0)
		}).Unwrap()

	if err != nil {
		t.Errorf("unexpected error, want nil, have %s", err.Error())
	}

	if value != 0 {
		t.Errorf("unexpected result on Err, want 0, have %d", value)
	}
}

func TestResult_And(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.And(Ok(2)).UnwrapUnsafe()

	if value != 2 {
		t.Errorf("unexpected result , want 2, have %d", value)
	}

	value, err := fail.And(Ok(1)).Unwrap()
	if err == nil {
		t.Errorf("unexpected result, want err but have none")
	}

	if value != 0 {
		t.Errorf("unexpected result on Err, want 0, have %d", value)
	}
}

func TestResult_AndThen(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.AndThen(func() int { return 2 }).UnwrapUnsafe()

	if value != 2 {
		t.Errorf("unexpected result , want 2, have %d", value)
	}

	value, err := fail.AndThen(func() int { return 1 }).Unwrap()
	if err == nil {
		t.Errorf("unexpected result, want err but have none")
	}

	if value != 0 {
		t.Errorf("unexpected result on Err, want 0, have %d", value)
	}
}

func TestResult_Map(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.Map(func(x int) int { return x + 1 }).UnwrapUnsafe()

	if value != 2 {
		t.Errorf("unexpected result , want 2, have %d", value)
	}

	value, err := fail.Map(func(x int) int { return x + 1 }).Unwrap()
	if err == nil {
		t.Errorf("unexpected result, want err but have none")
	}
}

func TestResult_MapOr(t *testing.T) {
	ok := Ok(1)
	fail := Err[int](errors.New("cannot divide by zero"))

	value := ok.MapOr(4, func(x int) int { return x + 1 }).UnwrapUnsafe()

	if value != 2 {
		t.Errorf("unexpected result, want 2, have %d", value)
	}

	value, err := fail.MapOr(1, func(x int) int { return x + 1 }).Unwrap()
	if err != nil {
		t.Errorf("unexpected result, want err but have none")
	}
	if value != 1 {
		t.Errorf("unexpected result, want 1, have %d", value)
	}
}

func TestResult_MapOrElse(t *testing.T) {
	ok := Ok(1)
	value := ok.MapOrElse(
		func(err error) int {
			return 1
		},
		func(x int) int {
			return x + 1
		},
	).UnwrapUnsafe()

	if value != 2 {
		t.Errorf("unexpected result, want 2, have %d", value)
	}

	fail := Err[int](errors.New("cannot divide by zero"))
	value, err := fail.MapOrElse(
		func(err error) int { return 1 },
		func(x int) int { return x + 1 },
	).Unwrap()

	if err != nil {
		t.Errorf("unexpected result, want err but have none")
	}
	if value != 1 {
		t.Errorf("unexpected result, want 1, have %d", value)
	}
}

// ExampleResult demonstrates basic usage of the Result type.
func ExampleResult() {
	// Create successful and error results
	success := Ok("Hello, World!")
	failure := Err[string](errors.New("something went wrong"))

	// Check if results are ok
	fmt.Printf("Success is ok: %v\n", success.IsOk())
	fmt.Printf("Failure is ok: %v\n", failure.IsOk())

	// Extract values safely
	if value, err := success.Unwrap(); err == nil {
		fmt.Printf("Success value: %s\n", value)
	}

	if _, err := failure.Unwrap(); err != nil {
		fmt.Printf("Failure error: %v\n", err)
	}

	// Output:
	// Success is ok: true
	// Failure is ok: false
	// Success value: Hello, World!
	// Failure error: something went wrong
}

// ExampleResult_Map demonstrates transforming values inside Result.
func ExampleResult_Map() {
	// Start with a result containing a number
	result := Ok(5)

	// Transform to its square
	squared := result.Map(func(x int) int { return x * x })

	// Transform an error result
	errorResult := Err[int](errors.New("invalid input"))
	errorSquared := errorResult.Map(func(x int) int { return x * x })

	fmt.Printf("Square of 5: %v\n", squared.UnwrapOr(0))
	fmt.Printf("Square of error: %v\n", errorSquared.UnwrapOr(-1))

	// Output:
	// Square of 5: 25
	// Square of error: -1
}

// ExampleResult_Match demonstrates pattern matching with Result.
func ExampleResult_Match() {
	// Helper function that may fail
	divide := func(x, y int) Result[int] {
		if y == 0 {
			return Err[int](errors.New("division by zero"))
		}
		return Ok(x / y)
	}

	// Pattern match on results
	success := divide(10, 2)
	failure := divide(10, 0)

	result1 := success.Match(
		func(value int) Result[int] {
			return Ok(value * 2)
		},
		func(err error) Result[int] {
			return Err[int](fmt.Errorf("handled: %w", err))
		},
	)

	result2 := failure.Match(
		func(value int) Result[int] {
			return Ok(value * 2)
		},
		func(err error) Result[int] {
			return Err[int](fmt.Errorf("handled: %w", err))
		},
	)

	fmt.Printf("Success result: %v\n", result1.UnwrapOr(-1))
	fmt.Printf("Failure handled: %v\n", result2.IsErr())

	// Output:
	// Success result: 10
	// Failure handled: true
}

// ExampleResult_Or demonstrates providing fallback results.
func ExampleResult_Or() {
	// Create primary and fallback results
	primary := Err[string](errors.New("primary failed"))
	fallback := Ok("fallback value")

	// Use fallback when primary fails
	result := primary.Or(fallback)

	fmt.Printf("Result: %s\n", result.UnwrapOr("default"))

	// Output:
	// Result: fallback value
}

// ExampleOk demonstrates creating a successful Result.
func ExampleOk() {
	// Create a successful result
	result := Ok("Success!")

	fmt.Printf("Is ok: %t\n", result.IsOk())
	fmt.Printf("Value: %s\n", result.UnwrapOr("default"))

	// Output:
	// Is ok: true
	// Value: Success!
}

// ExampleErr demonstrates creating an error Result.
func ExampleErr() {
	// Create an error result
	result := Err[string](errors.New("something failed"))

	fmt.Printf("Is error: %t\n", result.IsErr())
	fmt.Printf("Value: %s\n", result.UnwrapOr("default"))

	// Output:
	// Is error: true
	// Value: default
}

// ExampleOkZero demonstrates creating a Result with zero value.
func ExampleOkZero() {
	// Create a successful result with zero value
	result := OkZero[int]()

	fmt.Printf("Is ok: %t\n", result.IsOk())
	fmt.Printf("Value: %d\n", result.UnwrapOr(-1))

	// Output:
	// Is ok: true
	// Value: 0
}
