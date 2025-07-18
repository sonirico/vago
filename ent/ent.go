// Package ent provides utilities for managing environment variables in a type-safe manner.
package ent

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Get retrieves the value of an environment variable.
// If the environment variable is not set or is empty, it returns the fallback value.
func Get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Str is an alias for Get. It retrieves the value of an environment variable
// with a fallback value if the variable is not set or is empty.
func Str(key, fallback string) string {
	return Get(key, fallback)
}

// CondStrOrPanic conditionally retrieves an environment variable.
// If condition is true, it calls StrOrPanic to get the value (which panics if not found).
// If condition is false, it returns an empty string.
func CondStrOrPanic(condition bool, key string) string {
	if condition {
		return StrOrPanic(key)
	}
	return ""
}

// FixedStrOrPanic retrieves an environment variable and validates its length.
// It panics if the environment variable is not set or if its length doesn't match the expected length.
func FixedStrOrPanic(key string, length int) string {
	data := StrOrPanic(key)

	if len(data) != length {
		panic(fmt.Sprintf("expected env %s to have length = %d", key, length))
	}

	return data
}

// StrOrPanic retrieves an environment variable and panics if it's not set.
// Use this for required configuration values that must be present.
func StrOrPanic(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("env var is not defined: " + key)
}

// Float64 retrieves an environment variable and parses it as a float64.
// If the environment variable is not set, returns the fallback value.
// Panics if the value cannot be parsed as a float64.
func Float64(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		fval, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse float64: '%s', '%s'", key, value))
		}
		return fval
	}
	return fallback
}

// Int64 retrieves an environment variable and parses it as an int64.
// If the environment variable is not set, returns the fallback value.
// Panics if the value cannot be parsed as an int64.
func Int64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		fval, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse int64: '%s', '%s'", key, value))
		}
		return fval
	}

	return fallback
}

// Int64OrPanic retrieves an environment variable and parses it as an int64.
// Panics if the environment variable is not set or cannot be parsed as an int64.
func Int64OrPanic(key string) int64 {
	if value := os.Getenv(key); value != "" {
		fval, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse int64: '%s', '%s'", key, value))
		}
		return fval
	}

	panic("env var is not defined: " + key)
}

// Int retrieves an environment variable and parses it as an int.
// If the environment variable is not set, returns the fallback value.
// Panics if the value cannot be parsed as an int.
func Int(key string, fallback int) int {
	return int(Int64(key, int64(fallback)))
}

// IntOrPanic retrieves an environment variable and parses it as an int.
// Panics if the environment variable is not set or cannot be parsed as an int.
func IntOrPanic(key string) int {
	if value := os.Getenv(key); value != "" {
		fval, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse int64: '%s', '%s'", key, value))
		}
		return int(fval)
	}

	panic("env var is not defined: " + key)
}

// Duration retrieves an environment variable and parses it as a time.Duration.
// If the environment variable is not set, returns the fallback value.
// Panics if the value cannot be parsed as a duration.
func Duration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if val, err := time.ParseDuration(value); err == nil {
			return val
		}
		panic(fmt.Sprintf("Cannot parse duration: %s %s", key, value))
	}
	return fallback
}

// SliceStr retrieves an environment variable and parses it as a slice of strings.
// The value should be comma-separated. If the environment variable is not set, returns the fallback slice.
func SliceStr(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return fallback
}

// SliceInt retrieves an environment variable and parses it as a slice of integers.
// The value should be comma-separated integers. If the environment variable is not set, returns the fallback slice.
// Panics if any value cannot be parsed as an integer.
func SliceInt(key string, fallback []int) []int {
	if value := os.Getenv(key); value != "" {
		var res []int
		for _, in := range strings.Split(value, ",") {
			if val, err := strconv.Atoi(strings.TrimSpace(in)); err == nil {
				res = append(res, val)
			} else {
				panic(fmt.Sprintf("Cannot parse int: %s %s", key, value))
			}
		}
		return res
	}
	return fallback
}

// Bool retrieves an environment variable and parses it as a boolean.
// Accepts "true", "y", "yes", "YES", "1", "TRUE" as true values.
// All other values are considered false. If the environment variable is not set, returns the fallback value.
func Bool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		switch strings.ToLower(value) {
		case "true", "y", "yes", "YES", "1", "TRUE":
			return true
		default:
			return false
		}
	}

	return fallback
}

// JSON retrieves an environment variable and unmarshals it as JSON into type T.
// If the environment variable is not set, uses the fallback string.
// Returns the unmarshaled value and any JSON parsing error.
func JSON[T any](key string, fallback string) (T, error) {
	val := Get(key, fallback)
	var x T
	err := json.Unmarshal([]byte(val), &x)
	return x, err
}

// Enum retrieves an environment variable and validates it against allowed values.
// If the value is not in the allowed values list or if the environment variable is not set,
// returns the fallback value.
func Enum(key string, fallback string, values ...string) string {
	value := Get(key, fallback)
	if !contains(values, value) {
		return fallback
	}
	return value
}

// EnumOrPanic retrieves an environment variable and validates it against allowed values.
// Panics if the environment variable is not set or if the value is not in the allowed values list.
func EnumOrPanic(key string, values ...string) string {
	value := StrOrPanic(key)
	if !contains(values, value) {
		panic(fmt.Sprintf("value '%s' not in enum: '%v'", value, values))
	}
	return value
}

// contains checks if a value exists in a slice of strings.
// Returns true if the value is found, false otherwise.
func contains(samples []string, value string) bool {
	for _, sample := range samples {
		if sample == value {
			return true
		}
	}
	return false
}
