package opts

// Configurator is an interface for applying configuration to a value.
// This is the functional options pattern interface.
type Configurator[T any] interface {
	Apply(*T)
}

// Fn is a function that implements Configurator.
// It allows using functions as configurators in the functional options pattern.
type Fn[T any] func(*T)

// Apply applies the configuration function to the target.
func (fn Fn[T]) Apply(x *T) {
	fn(x)
}

// ApplyAll applies multiple configurators to a target value.
func ApplyAll[T any](target *T, configurators ...Configurator[T]) {
	for _, configurator := range configurators {
		configurator.Apply(target)
	}
}
