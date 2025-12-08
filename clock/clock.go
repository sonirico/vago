package clock

import (
	"sync"
	"time"
)

// Clock is an interface for getting the current time.
// This abstraction allows for easier testing by allowing time to be mocked.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the actual system time.
type RealClock struct{}

// Now returns the current system time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// New creates a new Clock that returns the actual system time.
func New() Clock {
	return RealClock{}
}

// MockClock implements Clock with a configurable time for testing.
// It is safe for concurrent use.
type MockClock struct {
	mu  sync.RWMutex
	now time.Time
}

// NewMock creates a new MockClock with the given initial time.
func NewMock(t time.Time) *MockClock {
	return &MockClock{now: t}
}

// Now returns the current mock time.
func (m *MockClock) Now() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.now
}

// Set updates the mock time to the given value.
func (m *MockClock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = t
}

// Add advances the mock time by the given duration.
func (m *MockClock) Add(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = m.now.Add(d)
}
