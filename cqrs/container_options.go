package cqrs

func ContainerMustProcessOrFail() Configurator[Container] {
	return configureFn[Container](func(c *Container) {
		c.mustProcessOrFail = true
	})
}

func ContainerWarnUnprocessed() Configurator[Container] {
	return configureFn[Container](func(c *Container) {
		c.warnUnprocessed = true
	})
}

func ContainerMustRestartOnError() Configurator[Container] {
	return configureFn[Container](func(c *Container) {
		c.restartOnError = true
	})
}

func ContainerDisableErrorCapture() Configurator[Container] {
	return configureFn[Container](func(c *Container) {
		c.errorCaptureDisabled = true
	})
}

func ContainerDisableAPM() Configurator[Container] {
	return configureFn[Container](func(c *Container) {
		c.apmDisabled = true
	})
}
