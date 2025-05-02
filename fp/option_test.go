package fp

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestOption(t *testing.T) {
	some := Some(1)
	if !some.IsSome() {
		t.Error("unexpected result, want ok, have none")
	}
	if some.IsNone() {
		t.Errorf("unexpected result, want some, have none")
	}

	none := None[any]()

	if none.IsSome() {
		t.Error("unexpected result, want none, have some")
	}
	if !none.IsNone() {
		t.Error("unexpected result, want none, have some")
	}

	value, ok := some.Unwrap()
	if !ok {
		t.Errorf("unexpected result, want some, have none")
	}

	if value != 1 {
		t.Errorf("unexpected value, want 1, have %d", value)
	}

	// UnwrapUnsafe
	_ = some.UnwrapUnsafe()
}

func TestOption_Or(t *testing.T) {
	some := Some(1)
	none := None[int]()

	value := some.Or(Some(2)).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = none.Or(Some(2)).UnwrapUnsafe()

	if value != 2 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}
}

func TestOption_OrElse(t *testing.T) {
	some := Some(1)
	none := None[int]()

	value := some.OrElse(func() Option[int] {
		return Some(2)
	}).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = none.OrElse(func() Option[int] {
		return Some(1)
	}).UnwrapUnsafe()

	if value != 1 {
		t.Errorf("unexpected result, want 1, have %d", value)
	}
}

func TestOption_UnwrapOr(t *testing.T) {
	some := Some(1)
	none := None[int]()

	value := some.UnwrapOr(2)

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = none.UnwrapOr(2)

	if value != 2 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}
}

func TestOption_UnwrapOrElse(t *testing.T) {
	some := Some(1)
	none := None[int]()

	value := some.UnwrapOrElse(func() int { return 2 })

	if value != 1 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}

	value = none.UnwrapOrElse(func() int { return 2 })

	if value != 2 {
		t.Errorf("unexpected result , want 1, have %d", value)
	}
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.UnwrapOrDefault()

	if value != "TOMBOLA" {
		t.Errorf("unexpected result , want TOMBOLA, have %s", value)
	}

	value = none.UnwrapOrDefault()

	if value != "" {
		t.Errorf("unexpected result , want zero value, have %s", value)
	}
}

func TestOption_Map(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.Map(func(x string) string {
		return strings.ToLower(x)
	}).UnwrapUnsafe()

	if value != "tombola" {
		t.Errorf("unexpected result , want tombola, have %s", value)
	}

	isNone := none.Map(func(x string) string { return "que pasa" }).IsNone()

	if !isNone {
		t.Error("unexpected result , want none, have some")
	}
}

func TestOption_MapOr(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.MapOr("MAYONESA", func(x string) string {
		return strings.ToLower(x)
	})

	if value != "tombola" {
		t.Errorf("unexpected result , want tombola, have %s", value)
	}

	value = none.MapOr("ALIOLI", func(x string) string { return "que pasa" })

	if value != "ALIOLI" {
		t.Errorf("unexpected result , want ALIOLI, have %s", value)
	}
}

func TestOption_MapOrElse(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.MapOrElse(
		func() string { return "MAYONESA" },
		func(x string) string {
			return strings.ToLower(x)
		},
	)

	if value != "tombola" {
		t.Errorf("unexpected result , want tombola, have %s", value)
	}

	value = none.MapOrElse(
		func() string { return "ALIOLI" },
		func(x string) string { return "que pasa" },
	)

	if value != "ALIOLI" {
		t.Errorf("unexpected result , want ALIOLI, have %s", value)
	}
}

func TestOption_OkOr(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.OkOr(io.EOF).UnwrapUnsafe()

	if value != "TOMBOLA" {
		t.Errorf("unexpected result , want tombola, have %s", value)
	}

	_, err := none.OkOr(io.EOF).Unwrap()
	if err == nil || !errors.Is(io.EOF, err) {
		t.Errorf("unexpected err, want io.EOF, have %v", err)
	}
}

func TestOption_OkOrElse(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.OkOrElse(func() error { return io.EOF }).UnwrapUnsafe()

	if value != "TOMBOLA" {
		t.Errorf("unexpected result , want tombola, have %s", value)
	}

	_, err := none.OkOrElse(func() error { return io.EOF }).Unwrap()
	if err == nil || !errors.Is(io.EOF, err) {
		t.Errorf("unexpected err, want io.EOF, have %v", err)
	}
}

func TestOption_Match(t *testing.T) {
	some := Some("TOMBOLA")
	none := None[string]()

	value := some.Match(
		func(x string) Option[string] { return Some(x + "S") },
		func() Option[string] { return Some("NADA") },
	).UnwrapUnsafe()

	if value != "TOMBOLAS" {
		t.Errorf("unexpected result , want TOMBOLAS, have %s", value)
	}

	value = none.Match(
		func(x string) Option[string] { return Some(x + "S") },
		func() Option[string] { return Some("test") },
	).UnwrapUnsafe()

	if value != "test" {
		t.Errorf("unexpected result, want test, have %s", value)
	}
}

func TestOptionFromTuple(t *testing.T) {
	option := OptionFromTuple(42, true)
	if !option.IsSome() {
		t.Error("unexpected result, want some, have none")
	}
	value, ok := option.Unwrap()
	if !ok || value != 42 {
		t.Errorf("unexpected result, want 42, have %d", value)
	}

	option = OptionFromTuple(0, false)
	if !option.IsNone() {
		t.Error("unexpected result, want none, have some")
	}
}

func TestOptionFromPtr(t *testing.T) {
	value := 42
	option := OptionFromPtr(&value)
	if !option.IsSome() {
		t.Error("unexpected result, want some, have none")
	}
	unwrappedValue, ok := option.Unwrap()
	if !ok || unwrappedValue != 42 {
		t.Errorf("unexpected result, want 42, have %d", unwrappedValue)
	}

	option = OptionFromPtr[int](nil)
	if !option.IsNone() {
		t.Error("unexpected result, want none, have some")
	}
}

func TestOptionFromZero(t *testing.T) {
	option := OptionFromZero(42)
	if !option.IsSome() {
		t.Error("unexpected result, want some, have none")
	}
	value, ok := option.Unwrap()
	if !ok || value != 42 {
		t.Errorf("unexpected result, want 42, have %d", value)
	}

	option = OptionFromZero(0)
	if !option.IsNone() {
		t.Error("unexpected result, want none, have some")
	}

	optionStr := OptionFromZero("hello")
	if !optionStr.IsSome() {
		t.Error("unexpected result, want some, have none")
	}
	strValue, ok := optionStr.Unwrap()
	if !ok || strValue != "hello" {
		t.Errorf("unexpected result, want 'hello', have '%s'", strValue)
	}

	optionStr = OptionFromZero("")
	if !optionStr.IsNone() {
		t.Error("unexpected result, want none, have some")
	}
}
