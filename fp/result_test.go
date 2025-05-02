package fp

import (
	"errors"
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
