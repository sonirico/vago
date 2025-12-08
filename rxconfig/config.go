package rxconfig

import (
	"encoding/json"
)

type Config[T any] interface {
	Get() T
	Clone() Config[T]
}

type Factory[T any] func(bts []byte) (Config[T], error)

type config[T any] struct {
	Value T
}

func (c *config[T]) Get() T {
	return c.Value
}

func (c *config[T]) Clone() Config[T] {
	return &config[T]{
		Value: c.Value,
	}
}

func newJSONConfig[T any](bts []byte) (Config[T], error) {
	var res T
	if err := json.Unmarshal(bts, &res); err != nil {
		return nil, err
	}

	return &config[T]{
		Value: res,
	}, nil
}

func NewJSONConfigFactory[T any]() Factory[T] {
	return func(bts []byte) (Config[T], error) {
		return newJSONConfig[T](bts)
	}
}

func Wrap[T any](x T) Config[T] {
	return &config[T]{Value: x}
}
