package clock

import (
	"sync"
	"testing"
	"time"
)

func TestRealClock_Now(t *testing.T) {
	c := New()
	before := time.Now()
	now := c.Now()
	after := time.Now()

	if now.Before(before) || now.After(after) {
		t.Errorf("Clock.Now() returned time outside expected range")
	}
}

func TestRealClock_Interface(t *testing.T) {
	var _ Clock = RealClock{}
	var _ Clock = New()
}

func TestMockClock_Now(t *testing.T) {
	fixedTime := time.Date(2025, 12, 8, 12, 0, 0, 0, time.UTC)
	m := NewMock(fixedTime)

	now := m.Now()
	if !now.Equal(fixedTime) {
		t.Errorf("MockClock.Now() = %v, want %v", now, fixedTime)
	}
}

func TestMockClock_Set(t *testing.T) {
	initialTime := time.Date(2025, 12, 8, 12, 0, 0, 0, time.UTC)
	newTime := time.Date(2025, 12, 9, 14, 30, 0, 0, time.UTC)

	m := NewMock(initialTime)
	m.Set(newTime)

	now := m.Now()
	if !now.Equal(newTime) {
		t.Errorf("After Set(), MockClock.Now() = %v, want %v", now, newTime)
	}
}

func TestMockClock_Add(t *testing.T) {
	initialTime := time.Date(2025, 12, 8, 12, 0, 0, 0, time.UTC)
	m := NewMock(initialTime)

	duration := 2 * time.Hour
	m.Add(duration)

	expected := initialTime.Add(duration)
	now := m.Now()
	if !now.Equal(expected) {
		t.Errorf("After Add(%v), MockClock.Now() = %v, want %v", duration, now, expected)
	}
}

func TestMockClock_AddMultiple(t *testing.T) {
	initialTime := time.Date(2025, 12, 8, 12, 0, 0, 0, time.UTC)
	m := NewMock(initialTime)

	m.Add(1 * time.Hour)
	m.Add(30 * time.Minute)
	m.Add(45 * time.Second)

	expected := initialTime.Add(1*time.Hour + 30*time.Minute + 45*time.Second)
	now := m.Now()
	if !now.Equal(expected) {
		t.Errorf("After multiple Add(), MockClock.Now() = %v, want %v", now, expected)
	}
}

func TestMockClock_Interface(t *testing.T) {
	var _ Clock = &MockClock{}
	var _ Clock = NewMock(time.Now())
}

func TestMockClock_Concurrent(t *testing.T) {
	initialTime := time.Date(2025, 12, 8, 12, 0, 0, 0, time.UTC)
	m := NewMock(initialTime)

	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Concurrent readers
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = m.Now()
			}
		}()
	}

	// Concurrent writers
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				m.Add(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// Verify the final time is after the initial time
	final := m.Now()
	if !final.After(initialTime) {
		t.Errorf("After concurrent operations, time should have advanced from %v, got %v", initialTime, final)
	}
}
