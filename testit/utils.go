package testit

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func RunSafe(
	m *testing.M,
	log Logger,
	pool *DockerResourcesPool,
) {
	var v int

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panic: %v", err)
		}

		snaps.Clean(m)
		pool.Down()
		os.Exit(v)
	}()

	v = m.Run()
}
