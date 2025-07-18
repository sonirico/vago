package ent

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		setEnv   bool
		expected string
	}{
		{
			name:     "existing env var",
			key:      "TEST_GET_EXISTING",
			fallback: "default",
			envValue: "value_from_env",
			setEnv:   true,
			expected: "value_from_env",
		},
		{
			name:     "missing env var returns fallback",
			key:      "TEST_GET_MISSING",
			fallback: "default_value",
			envValue: "",
			setEnv:   false,
			expected: "default_value",
		},
		{
			name:     "empty env var returns fallback",
			key:      "TEST_GET_EMPTY",
			fallback: "default_value",
			envValue: "",
			setEnv:   true,
			expected: "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			result := Get(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleGet demonstrates how to retrieve environment variables with fallback values.
func ExampleGet() {
	// Set an environment variable
	os.Setenv("APP_NAME", "MyApplication")
	defer os.Unsetenv("APP_NAME")

	// Get existing env var
	appName := Get("APP_NAME", "DefaultApp")
	fmt.Println(appName)

	// Get non-existing env var with fallback
	dbHost := Get("DB_HOST", "localhost")
	fmt.Println(dbHost)
	// Output:
	//
	// MyApplication
	// localhost
}

func TestStr(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		setEnv   bool
		expected string
	}{
		{
			name:     "existing env var",
			key:      "TEST_STR_EXISTING",
			fallback: "default",
			envValue: "string_value",
			setEnv:   true,
			expected: "string_value",
		},
		{
			name:     "missing env var returns fallback",
			key:      "TEST_STR_MISSING",
			fallback: "fallback_string",
			envValue: "",
			setEnv:   false,
			expected: "fallback_string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			result := Str(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleStr demonstrates how to retrieve string environment variables with fallback values.
func ExampleStr() {
	// Set environment variable
	os.Setenv("USER_NAME", "john_doe")
	defer os.Unsetenv("USER_NAME")

	username := Str("USER_NAME", "anonymous")
	fmt.Println(username)
	// Output:
	//
	// john_doe
}

func TestCondStrOrPanic(t *testing.T) {
	tests := []struct {
		name        string
		condition   bool
		key         string
		envValue    string
		setEnv      bool
		expected    string
		shouldPanic bool
	}{
		{
			name:        "condition true with existing env var",
			condition:   true,
			key:         "TEST_COND_TRUE_EXISTING",
			envValue:    "conditional_value",
			setEnv:      true,
			expected:    "conditional_value",
			shouldPanic: false,
		},
		{
			name:        "condition false",
			condition:   false,
			key:         "TEST_COND_FALSE",
			envValue:    "",
			setEnv:      false,
			expected:    "",
			shouldPanic: false,
		},
		{
			name:        "condition true with missing env var should panic",
			condition:   true,
			key:         "TEST_COND_MISSING",
			envValue:    "",
			setEnv:      false,
			expected:    "",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					CondStrOrPanic(tt.condition, tt.key)
				})
			} else {
				result := CondStrOrPanic(tt.condition, tt.key)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleCondStrOrPanic demonstrates how to conditionally retrieve environment variables with panic on missing values.
func ExampleCondStrOrPanic() {
	// Set environment variable
	os.Setenv("DEBUG_MODE", "true")
	defer os.Unsetenv("DEBUG_MODE")

	// Get value only if condition is true
	debugMode := CondStrOrPanic(true, "DEBUG_MODE")
	fmt.Println(debugMode)

	// Returns empty string if condition is false
	emptyValue := CondStrOrPanic(false, "DEBUG_MODE")
	fmt.Println(emptyValue)
	// Output:
	//
	// true
	//
}

func TestFixedStrOrPanic(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		length      int
		envValue    string
		setEnv      bool
		shouldPanic bool
	}{
		{
			name:        "correct length",
			key:         "TEST_FIXED_CORRECT",
			length:      5,
			envValue:    "12345",
			setEnv:      true,
			shouldPanic: false,
		},
		{
			name:        "incorrect length should panic",
			key:         "TEST_FIXED_WRONG",
			length:      5,
			envValue:    "123",
			setEnv:      true,
			shouldPanic: true,
		},
		{
			name:        "missing env var should panic",
			key:         "TEST_FIXED_MISSING",
			length:      5,
			envValue:    "",
			setEnv:      false,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					FixedStrOrPanic(tt.key, tt.length)
				})
			} else {
				result := FixedStrOrPanic(tt.key, tt.length)
				assert.Equal(t, tt.envValue, result)
				assert.Len(t, result, tt.length)
			}
		})
	}
}

// ExampleFixedStrOrPanic demonstrates how to retrieve environment variables with length validation.
func ExampleFixedStrOrPanic() {
	// Set environment variable with exact length
	os.Setenv("API_KEY", "abc123")
	defer os.Unsetenv("API_KEY")

	// Get value with length validation
	apiKey := FixedStrOrPanic("API_KEY", 6)
	fmt.Println(apiKey)
	// Output:
	//
	// abc123
}

func TestStrOrPanic(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		envValue    string
		setEnv      bool
		shouldPanic bool
	}{
		{
			name:        "existing env var",
			key:         "TEST_STR_PANIC_EXISTING",
			envValue:    "panic_value",
			setEnv:      true,
			shouldPanic: false,
		},
		{
			name:        "missing env var should panic",
			key:         "TEST_STR_PANIC_MISSING",
			envValue:    "",
			setEnv:      false,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					StrOrPanic(tt.key)
				})
			} else {
				result := StrOrPanic(tt.key)
				assert.Equal(t, tt.envValue, result)
			}
		})
	}
}

// ExampleStrOrPanic demonstrates how to retrieve required environment variables that panic if missing.
func ExampleStrOrPanic() {
	// Set required environment variable
	os.Setenv("REQUIRED_CONFIG", "important_value")
	defer os.Unsetenv("REQUIRED_CONFIG")

	// Get required value (panics if missing)
	config := StrOrPanic("REQUIRED_CONFIG")
	fmt.Println(config)
	// Output:
	//
	// important_value
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		fallback    float64
		envValue    string
		setEnv      bool
		expected    float64
		shouldPanic bool
	}{
		{
			name:        "valid float",
			key:         "TEST_FLOAT64_VALID",
			fallback:    1.0,
			envValue:    "3.14159",
			setEnv:      true,
			expected:    3.14159,
			shouldPanic: false,
		},
		{
			name:        "missing env var returns fallback",
			key:         "TEST_FLOAT64_MISSING",
			fallback:    2.5,
			envValue:    "",
			setEnv:      false,
			expected:    2.5,
			shouldPanic: false,
		},
		{
			name:        "invalid float should panic",
			key:         "TEST_FLOAT64_INVALID",
			fallback:    1.0,
			envValue:    "not_a_float",
			setEnv:      true,
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					Float64(tt.key, tt.fallback)
				})
			} else {
				result := Float64(tt.key, tt.fallback)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleFloat64 demonstrates how to retrieve float64 environment variables with fallback values.
func ExampleFloat64() {
	// Set environment variable
	os.Setenv("PRICE", "19.99")
	defer os.Unsetenv("PRICE")

	price := Float64("PRICE", 0.0)
	fmt.Println(price)
	// Output:
	//
	// 19.99
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		fallback    int64
		envValue    string
		setEnv      bool
		expected    int64
		shouldPanic bool
	}{
		{
			name:        "valid int64",
			key:         "TEST_INT64_VALID",
			fallback:    10,
			envValue:    "12345",
			setEnv:      true,
			expected:    12345,
			shouldPanic: false,
		},
		{
			name:        "missing env var returns fallback",
			key:         "TEST_INT64_MISSING",
			fallback:    42,
			envValue:    "",
			setEnv:      false,
			expected:    42,
			shouldPanic: false,
		},
		{
			name:        "invalid int64 should panic",
			key:         "TEST_INT64_INVALID",
			fallback:    10,
			envValue:    "not_an_int",
			setEnv:      true,
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					Int64(tt.key, tt.fallback)
				})
			} else {
				result := Int64(tt.key, tt.fallback)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleInt64 demonstrates how to retrieve int64 environment variables with fallback values.
func ExampleInt64() {
	// Set environment variable
	os.Setenv("MAX_CONNECTIONS", "100")
	defer os.Unsetenv("MAX_CONNECTIONS")

	maxConn := Int64("MAX_CONNECTIONS", 50)
	fmt.Println(maxConn)
	// Output:
	//
	// 100
}

func TestInt64OrPanic(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		envValue    string
		setEnv      bool
		expected    int64
		shouldPanic bool
	}{
		{
			name:        "valid int64",
			key:         "TEST_INT64_PANIC_VALID",
			envValue:    "98765",
			setEnv:      true,
			expected:    98765,
			shouldPanic: false,
		},
		{
			name:        "missing env var should panic",
			key:         "TEST_INT64_PANIC_MISSING",
			envValue:    "",
			setEnv:      false,
			expected:    0,
			shouldPanic: true,
		},
		{
			name:        "invalid int64 should panic",
			key:         "TEST_INT64_PANIC_INVALID",
			envValue:    "invalid",
			setEnv:      true,
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					Int64OrPanic(tt.key)
				})
			} else {
				result := Int64OrPanic(tt.key)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleInt64OrPanic demonstrates how to retrieve required int64 environment variables that panic if missing or invalid.
func ExampleInt64OrPanic() {
	// Set required environment variable
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	port := Int64OrPanic("PORT")
	fmt.Println(port)
	// Output:
	//
	// 8080
}

func TestInt(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		fallback    int
		envValue    string
		setEnv      bool
		expected    int
		shouldPanic bool
	}{
		{
			name:        "valid int",
			key:         "TEST_INT_VALID",
			fallback:    5,
			envValue:    "123",
			setEnv:      true,
			expected:    123,
			shouldPanic: false,
		},
		{
			name:        "missing env var returns fallback",
			key:         "TEST_INT_MISSING",
			fallback:    99,
			envValue:    "",
			setEnv:      false,
			expected:    99,
			shouldPanic: false,
		},
		{
			name:        "invalid int should panic",
			key:         "TEST_INT_INVALID",
			fallback:    5,
			envValue:    "not_int",
			setEnv:      true,
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					Int(tt.key, tt.fallback)
				})
			} else {
				result := Int(tt.key, tt.fallback)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleInt demonstrates how to retrieve integer environment variables with fallback values.
func ExampleInt() {
	// Set environment variable
	os.Setenv("WORKER_COUNT", "4")
	defer os.Unsetenv("WORKER_COUNT")

	workers := Int("WORKER_COUNT", 1)
	fmt.Println(workers)
	// Output:
	//
	// 4
}

func TestIntOrPanic(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		envValue    string
		setEnv      bool
		expected    int
		shouldPanic bool
	}{
		{
			name:        "valid int",
			key:         "TEST_INT_PANIC_VALID",
			envValue:    "456",
			setEnv:      true,
			expected:    456,
			shouldPanic: false,
		},
		{
			name:        "missing env var should panic",
			key:         "TEST_INT_PANIC_MISSING",
			envValue:    "",
			setEnv:      false,
			expected:    0,
			shouldPanic: true,
		},
		{
			name:        "invalid int should panic",
			key:         "TEST_INT_PANIC_INVALID",
			envValue:    "invalid",
			setEnv:      true,
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					IntOrPanic(tt.key)
				})
			} else {
				result := IntOrPanic(tt.key)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleIntOrPanic demonstrates how to retrieve required integer environment variables that panic if missing or invalid.
func ExampleIntOrPanic() {
	// Set required environment variable
	os.Setenv("TIMEOUT", "30")
	defer os.Unsetenv("TIMEOUT")

	timeout := IntOrPanic("TIMEOUT")
	fmt.Println(timeout)
	// Output:
	//
	// 30
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		fallback    time.Duration
		envValue    string
		setEnv      bool
		expected    time.Duration
		shouldPanic bool
	}{
		{
			name:        "valid duration",
			key:         "TEST_DURATION_VALID",
			fallback:    time.Second,
			envValue:    "5m30s",
			setEnv:      true,
			expected:    5*time.Minute + 30*time.Second,
			shouldPanic: false,
		},
		{
			name:        "missing env var returns fallback",
			key:         "TEST_DURATION_MISSING",
			fallback:    time.Hour,
			envValue:    "",
			setEnv:      false,
			expected:    time.Hour,
			shouldPanic: false,
		},
		{
			name:        "invalid duration should panic",
			key:         "TEST_DURATION_INVALID",
			fallback:    time.Second,
			envValue:    "invalid_duration",
			setEnv:      true,
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					Duration(tt.key, tt.fallback)
				})
			} else {
				result := Duration(tt.key, tt.fallback)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleDuration demonstrates how to retrieve time.Duration environment variables with fallback values.
func ExampleDuration() {
	// Set environment variable
	os.Setenv("REQUEST_TIMEOUT", "30s")
	defer os.Unsetenv("REQUEST_TIMEOUT")

	timeout := Duration("REQUEST_TIMEOUT", 10*time.Second)
	fmt.Println(timeout)
	// Output:
	//
	// 30s
}

func TestSliceStr(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback []string
		envValue string
		setEnv   bool
		expected []string
	}{
		{
			name:     "valid string slice",
			key:      "TEST_SLICE_STR_VALID",
			fallback: []string{"default"},
			envValue: "apple,banana,cherry",
			setEnv:   true,
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "missing env var returns fallback",
			key:      "TEST_SLICE_STR_MISSING",
			fallback: []string{"fallback1", "fallback2"},
			envValue: "",
			setEnv:   false,
			expected: []string{"fallback1", "fallback2"},
		},
		{
			name:     "single value",
			key:      "TEST_SLICE_STR_SINGLE",
			fallback: []string{"default"},
			envValue: "single",
			setEnv:   true,
			expected: []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			result := SliceStr(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleSliceStr demonstrates how to retrieve string slice environment variables from comma-separated values.
func ExampleSliceStr() {
	// Set environment variable
	os.Setenv("ALLOWED_HOSTS", "localhost,127.0.0.1,example.com")
	defer os.Unsetenv("ALLOWED_HOSTS")

	hosts := SliceStr("ALLOWED_HOSTS", []string{"localhost"})
	fmt.Println(len(hosts))
	// Output:
	//
	// 3
}

func TestSliceInt(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		fallback    []int
		envValue    string
		setEnv      bool
		expected    []int
		shouldPanic bool
	}{
		{
			name:        "valid int slice",
			key:         "TEST_SLICE_INT_VALID",
			fallback:    []int{1},
			envValue:    "10,20,30",
			setEnv:      true,
			expected:    []int{10, 20, 30},
			shouldPanic: false,
		},
		{
			name:        "missing env var returns fallback",
			key:         "TEST_SLICE_INT_MISSING",
			fallback:    []int{1, 2, 3},
			envValue:    "",
			setEnv:      false,
			expected:    []int{1, 2, 3},
			shouldPanic: false,
		},
		{
			name:        "invalid int should panic",
			key:         "TEST_SLICE_INT_INVALID",
			fallback:    []int{1},
			envValue:    "10,invalid,30",
			setEnv:      true,
			expected:    nil,
			shouldPanic: true,
		},
		{
			name:        "with spaces",
			key:         "TEST_SLICE_INT_SPACES",
			fallback:    []int{1},
			envValue:    "10, 20 , 30",
			setEnv:      true,
			expected:    []int{10, 20, 30},
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					SliceInt(tt.key, tt.fallback)
				})
			} else {
				result := SliceInt(tt.key, tt.fallback)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleSliceInt demonstrates how to retrieve integer slice environment variables from comma-separated values.
func ExampleSliceInt() {
	// Set environment variable
	os.Setenv("PORTS", "8080,8081,8082")
	defer os.Unsetenv("PORTS")

	ports := SliceInt("PORTS", []int{3000})
	fmt.Println(len(ports))
	// Output:
	//
	// 3
}

func TestBool(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback bool
		envValue string
		setEnv   bool
		expected bool
	}{
		{
			name:     "true value",
			key:      "TEST_BOOL_TRUE",
			fallback: false,
			envValue: "true",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "TRUE value",
			key:      "TEST_BOOL_TRUE_UPPER",
			fallback: false,
			envValue: "TRUE",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "y value",
			key:      "TEST_BOOL_Y",
			fallback: false,
			envValue: "y",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "yes value",
			key:      "TEST_BOOL_YES",
			fallback: false,
			envValue: "yes",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "YES value",
			key:      "TEST_BOOL_YES_UPPER",
			fallback: false,
			envValue: "YES",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "1 value",
			key:      "TEST_BOOL_1",
			fallback: false,
			envValue: "1",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "false value",
			key:      "TEST_BOOL_FALSE",
			fallback: true,
			envValue: "false",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "0 value",
			key:      "TEST_BOOL_0",
			fallback: true,
			envValue: "0",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "missing env var returns fallback",
			key:      "TEST_BOOL_MISSING",
			fallback: true,
			envValue: "",
			setEnv:   false,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			result := Bool(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleBool demonstrates how to retrieve boolean environment variables with various true/false representations.
func ExampleBool() {
	// Set environment variable
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	debug := Bool("DEBUG", false)
	fmt.Println(debug)
	// Output:
	//
	// true
}

func TestJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name        string
		key         string
		fallback    string
		envValue    string
		setEnv      bool
		expected    TestStruct
		expectError bool
	}{
		{
			name:        "valid JSON",
			key:         "TEST_JSON_VALID",
			fallback:    `{"name":"default","age":0}`,
			envValue:    `{"name":"John","age":30}`,
			setEnv:      true,
			expected:    TestStruct{Name: "John", Age: 30},
			expectError: false,
		},
		{
			name:        "missing env var uses fallback",
			key:         "TEST_JSON_MISSING",
			fallback:    `{"name":"fallback","age":25}`,
			envValue:    "",
			setEnv:      false,
			expected:    TestStruct{Name: "fallback", Age: 25},
			expectError: false,
		},
		{
			name:        "invalid JSON returns error",
			key:         "TEST_JSON_INVALID",
			fallback:    `{"name":"default","age":0}`,
			envValue:    `invalid json`,
			setEnv:      true,
			expected:    TestStruct{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			result, err := JSON[TestStruct](tt.key, tt.fallback)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleJSON demonstrates how to parse JSON environment variables into Go structs with type safety.
func ExampleJSON() {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	// Set environment variable
	os.Setenv("DB_CONFIG", `{"host":"localhost","port":5432}`)
	defer os.Unsetenv("DB_CONFIG")

	config, err := JSON[Config]("DB_CONFIG", `{"host":"127.0.0.1","port":3306}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Host)
	// Output:
	//
	// localhost
}

func TestEnum(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		values   []string
		envValue string
		setEnv   bool
		expected string
	}{
		{
			name:     "valid enum value",
			key:      "TEST_ENUM_VALID",
			fallback: "dev",
			values:   []string{"dev", "staging", "prod"},
			envValue: "staging",
			setEnv:   true,
			expected: "staging",
		},
		{
			name:     "invalid enum value returns fallback",
			key:      "TEST_ENUM_INVALID",
			fallback: "dev",
			values:   []string{"dev", "staging", "prod"},
			envValue: "invalid",
			setEnv:   true,
			expected: "dev",
		},
		{
			name:     "missing env var returns fallback",
			key:      "TEST_ENUM_MISSING",
			fallback: "dev",
			values:   []string{"dev", "staging", "prod"},
			envValue: "",
			setEnv:   false,
			expected: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			result := Enum(tt.key, tt.fallback, tt.values...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleEnum demonstrates how to retrieve environment variables with validation against allowed values.
func ExampleEnum() {
	// Set environment variable
	os.Setenv("ENV", "staging")
	defer os.Unsetenv("ENV")

	env := Enum("ENV", "dev", "dev", "staging", "prod")
	fmt.Println(env)
	// Output:
	//
	// staging
}

func TestEnumOrPanic(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		values      []string
		envValue    string
		setEnv      bool
		expected    string
		shouldPanic bool
	}{
		{
			name:        "valid enum value",
			key:         "TEST_ENUM_PANIC_VALID",
			values:      []string{"small", "medium", "large"},
			envValue:    "medium",
			setEnv:      true,
			expected:    "medium",
			shouldPanic: false,
		},
		{
			name:        "invalid enum value should panic",
			key:         "TEST_ENUM_PANIC_INVALID",
			values:      []string{"small", "medium", "large"},
			envValue:    "invalid",
			setEnv:      true,
			expected:    "",
			shouldPanic: true,
		},
		{
			name:        "missing env var should panic",
			key:         "TEST_ENUM_PANIC_MISSING",
			values:      []string{"small", "medium", "large"},
			envValue:    "",
			setEnv:      false,
			expected:    "",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					EnumOrPanic(tt.key, tt.values...)
				})
			} else {
				result := EnumOrPanic(tt.key, tt.values...)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ExampleEnumOrPanic demonstrates how to retrieve required environment variables with validation against allowed values, panicking if invalid.
func ExampleEnumOrPanic() {
	// Set environment variable
	os.Setenv("LOG_LEVEL", "info")
	defer os.Unsetenv("LOG_LEVEL")

	logLevel := EnumOrPanic("LOG_LEVEL", "debug", "info", "warn", "error")
	fmt.Println(logLevel)
	// Output:
	//
	// info
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		samples  []string
		value    string
		expected bool
	}{
		{
			name:     "value exists",
			samples:  []string{"apple", "banana", "cherry"},
			value:    "banana",
			expected: true,
		},
		{
			name:     "value does not exist",
			samples:  []string{"apple", "banana", "cherry"},
			value:    "orange",
			expected: false,
		},
		{
			name:     "empty slice",
			samples:  []string{},
			value:    "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.samples, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ExampleContains demonstrates how to check if a value exists in a slice of strings.
func ExampleContains() {
	fruits := []string{"apple", "banana", "cherry"}

	hasBanana := contains(fruits, "banana")
	fmt.Println(hasBanana)

	hasOrange := contains(fruits, "orange")
	fmt.Println(hasOrange)
	// Output:
	//
	// true
	// false
}
