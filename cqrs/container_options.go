package cqrs

import "github.com/sonirico/vago/opts"

func ContainerMustProcessOrFail() opts.Configurator[Container] {
	return opts.Fn[Container](func(c *Container) {
		c.mustProcessOrFail = true
	})
}

func ContainerWarnUnprocessed() opts.Configurator[Container] {
	return opts.Fn[Container](func(c *Container) {
		c.warnUnprocessed = true
	})
}

func ContainerMustRestartOnError() opts.Configurator[Container] {
	return opts.Fn[Container](func(c *Container) {
		c.restartOnError = true
	})
}

func ContainerDisableErrorCapture() opts.Configurator[Container] {
	return opts.Fn[Container](func(c *Container) {
		c.errorCaptureDisabled = true
	})
}

func ContainerDisableAPM() opts.Configurator[Container] {
	return opts.Fn[Container](func(c *Container) {
		c.apmDisabled = true
	})
}
