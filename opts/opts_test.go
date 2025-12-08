package opts

import (
	"testing"
)

type config struct {
	name  string
	value int
	debug bool
}

func TestFn_Apply(t *testing.T) {
	c := &config{}

	fn := Fn[config](func(cfg *config) {
		cfg.name = "test"
		cfg.value = 42
	})

	fn.Apply(c)

	if c.name != "test" {
		t.Errorf("expected name test, got %s", c.name)
	}
	if c.value != 42 {
		t.Errorf("expected value 42, got %d", c.value)
	}
}

func TestApplyAll(t *testing.T) {
	c := &config{}

	opts := []Configurator[config]{
		Fn[config](func(cfg *config) { cfg.name = "applied" }),
		Fn[config](func(cfg *config) { cfg.value = 100 }),
		Fn[config](func(cfg *config) { cfg.debug = true }),
	}

	ApplyAll(c, opts...)

	if c.name != "applied" {
		t.Errorf("expected name applied, got %s", c.name)
	}
	if c.value != 100 {
		t.Errorf("expected value 100, got %d", c.value)
	}
	if !c.debug {
		t.Errorf("expected debug true, got false")
	}
}

func TestApplyAll_Empty(t *testing.T) {
	c := &config{name: "default", value: 1}

	ApplyAll(c)

	if c.name != "default" || c.value != 1 {
		t.Error("ApplyAll with no opts should not modify config")
	}
}

// Example functions that return configurators
func WithName(name string) Configurator[config] {
	return Fn[config](func(cfg *config) {
		cfg.name = name
	})
}

func WithValue(value int) Configurator[config] {
	return Fn[config](func(cfg *config) {
		cfg.value = value
	})
}

func WithDebug(debug bool) Configurator[config] {
	return Fn[config](func(cfg *config) {
		cfg.debug = debug
	})
}

func TestFunctionalOptions(t *testing.T) {
	c := &config{}

	ApplyAll(c,
		WithName("functional"),
		WithValue(999),
		WithDebug(true),
	)

	if c.name != "functional" {
		t.Errorf("expected name functional, got %s", c.name)
	}
	if c.value != 999 {
		t.Errorf("expected value 999, got %d", c.value)
	}
	if !c.debug {
		t.Errorf("expected debug true, got false")
	}
}
