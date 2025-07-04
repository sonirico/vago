
<div align="center">
  <img src="vago.png" alt="Visigoth" width="200"/>
  
  # vago
  
  The ultimate toolkit for vaGo developers. A comprehensive collection of functions, data structures, and utilities designed to enhance productivity and code quality with no learning curve and less effort.

  [![Go Report Card](https://goreportcard.com/badge/github.com/sonirico/vago)](https://goreportcard.com/report/github.com/sonirico/vago)
  [![Go Reference](https://pkg.go.dev/badge/github.com/sonirico/vago.svg)](https://pkg.go.dev/github.com/sonirico/vago)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Release](https://img.shields.io/github/v/release/sonirico/vago.svg)](https://github.com/sonirico/vago/releases)
</div>

📖 **[View full documentation and examples on pkg.go.dev →](https://pkg.go.dev/github.com/sonirico/vago)**

## ✨ Workspace Architecture

This project leverages Go workspaces to provide **isolated dependencies** for each module. This means:

- 🎯 **Lightweight imports**: When you import `fp` or `streams`, you won't download database drivers or logging dependencies
- 🔧 **Modular design**: Each module (`db`, `lol`, `num`) maintains its own `go.mod` with specific dependencies
- 📦 **Zero bloat**: Use only what you need without carrying unnecessary dependencies
- 🚀 **Fast builds**: Smaller dependency graphs lead to faster compilation and smaller binaries

**Example**: Importing `github.com/sonirico/vago/fp` will only pull functional programming utilities, not database connections or logging frameworks.

## Modules

## <a name="table-of-contents"></a>Table of Contents

- [🪄 Fp](#fp) - 15 functions
- [📝 Lol](#lol) - 6 functions
- [🗝️ Maps](#maps) - 8 functions
- [🔢 Num](#num) - 14 functions
- [⛓️ Slices](#slices) - 10 functions
- [🌊 Streams](#streams) - 26 functions
- [🔞 Zero](#zero) - 3 functions

## <a name="fp"></a>🪄 Fp

Functional programming utilities including Option and Result types.

### Functions

- [Err](#fp-err)
- [None](#fp-none)
- [Ok](#fp-ok)
- [OkZero](#fp-okzero)
- [Option](#fp-option)
- [OptionFromPtr](#fp-optionfromptr)
- [OptionFromTuple](#fp-optionfromtuple)
- [Option_Map](#fp-option_map)
- [Option_Match](#fp-option_match)
- [Option_Or](#fp-option_or)
- [Result](#fp-result)
- [Result_Map](#fp-result_map)
- [Result_Match](#fp-result_match)
- [Result_Or](#fp-result_or)
- [Some](#fp-some)

#### fp Err

ExampleErr demonstrates creating an error Result.


<details><summary>Code</summary>

```go
func ExampleErr() {
	// Create an error result
	result := Err[string](errors.New("something failed"))

	fmt.Printf("Is error: %t\n", result.IsErr())
	fmt.Printf("Value: %s\n", result.UnwrapOr("default"))

	// Output:
	// Is error: true
	// Value: default
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp None

ExampleNone demonstrates creating an empty Option.


<details><summary>Code</summary>

```go
func ExampleNone() {
	// Create an empty Option
	empty := None[string]()

	fmt.Printf("Has value: %t\n", empty.IsSome())
	fmt.Printf("Value: %s\n", empty.UnwrapOr("default"))

	// Output:
	// Has value: false
	// Value: default
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Ok

ExampleOk demonstrates creating a successful Result.


<details><summary>Code</summary>

```go
func ExampleOk() {
	// Create a successful result
	result := Ok("Success!")

	fmt.Printf("Is ok: %t\n", result.IsOk())
	fmt.Printf("Value: %s\n", result.UnwrapOr("default"))

	// Output:
	// Is ok: true
	// Value: Success!
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp OkZero

ExampleOkZero demonstrates creating a Result with zero value.


<details><summary>Code</summary>

```go
func ExampleOkZero() {
	// Create a successful result with zero value
	result := OkZero[int]()

	fmt.Printf("Is ok: %t\n", result.IsOk())
	fmt.Printf("Value: %d\n", result.UnwrapOr(-1))

	// Output:
	// Is ok: true
	// Value: 0
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Option

ExampleOption demonstrates basic usage of the Option type.


<details><summary>Code</summary>

```go
func ExampleOption() {
	// Create Some and None options
	someValue := Some(42)
	noneValue := None[int]()

	// Check if options have values
	fmt.Printf("Some has value: %v\n", someValue.IsSome())
	fmt.Printf("None has value: %v\n", noneValue.IsSome())

	// Extract values safely
	if value, ok := someValue.Unwrap(); ok {
		fmt.Printf("Value: %d\n", value)
	}

	// Output:
	// Some has value: true
	// None has value: false
	// Value: 42
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp OptionFromPtr

ExampleOptionFromPtr demonstrates creating Option from a pointer.


<details><summary>Code</summary>

```go
func ExampleOptionFromPtr() {
	// From valid pointer
	value := "hello"
	opt1 := OptionFromPtr(&value)

	// From nil pointer
	var nilPtr *string
	opt2 := OptionFromPtr(nilPtr)

	fmt.Printf("From pointer: %s\n", opt1.UnwrapOr("empty"))
	fmt.Printf("From nil: %s\n", opt2.UnwrapOr("empty"))

	// Output:
	// From pointer: hello
	// From nil: empty
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp OptionFromTuple

ExampleOptionFromTuple demonstrates creating Option from a tuple pattern.


<details><summary>Code</summary>

```go
func ExampleOptionFromTuple() {
	// Common Go pattern: value, ok
	getValue := func(key string) (string, bool) {
		data := map[string]string{"name": "Alice", "age": "25"}
		value, ok := data[key]
		return value, ok
	}

	// Convert to Option
	nameOpt := OptionFromTuple(getValue("name"))
	missingOpt := OptionFromTuple(getValue("missing"))

	fmt.Printf("Name: %s\n", nameOpt.UnwrapOr("unknown"))
	fmt.Printf("Missing: %s\n", missingOpt.UnwrapOr("unknown"))

	// Output:
	// Name: Alice
	// Missing: unknown
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Option_Map

ExampleOption_Map demonstrates transforming values inside Option.


<details><summary>Code</summary>

```go
func ExampleOption_Map() {
	// Start with an optional number
	maybeNumber := Some(5)

	// Transform it to its square
	maybeSquare := maybeNumber.Map(func(x int) int { return x * x })

	// Transform None value
	noneNumber := None[int]()
	noneSquare := noneNumber.Map(func(x int) int { return x * x })

	fmt.Printf("Square of 5: %v\n", maybeSquare.UnwrapOr(0))
	fmt.Printf("Square of None: %v\n", noneSquare.UnwrapOr(-1))

	// Output:
	// Square of 5: 25
	// Square of None: -1
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Option_Match

ExampleOption_Match demonstrates pattern matching with Option.


<details><summary>Code</summary>

```go
func ExampleOption_Match() {
	// Helper function that may return a value
	getValue := func(id int) Option[string] {
		if id > 0 {
			return Some(fmt.Sprintf("User_%d", id))
		}
		return None[string]()
	}

	// Pattern match on the result
	validUser := getValue(42)
	invalidUser := getValue(-1)

	result1 := validUser.Match(
		func(user string) Option[string] {
			return Some("Found: " + user)
		},
		func() Option[string] {
			return Some("No user found")
		},
	)

	result2 := invalidUser.Match(
		func(user string) Option[string] {
			return Some("Found: " + user)
		},
		func() Option[string] {
			return Some("No user found")
		},
	)

	fmt.Printf("Valid user: %s\n", result1.UnwrapOr(""))
	fmt.Printf("Invalid user: %s\n", result2.UnwrapOr(""))

	// Output:
	// Valid user: Found: User_42
	// Invalid user: No user found
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Option_Or

ExampleOption_Or demonstrates providing fallback values.


<details><summary>Code</summary>

```go
func ExampleOption_Or() {
	// Create some options
	primary := None[string]()
	secondary := Some("backup")
	tertiary := Some("fallback")

	// Chain fallbacks
	result := primary.Or(secondary).Or(tertiary)

	fmt.Printf("Result: %s\n", result.UnwrapOr("default"))

	// Output:
	// Result: backup
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Result

ExampleResult demonstrates basic usage of the Result type.


<details><summary>Code</summary>

```go
func ExampleResult() {
	// Create successful and error results
	success := Ok("Hello, World!")
	failure := Err[string](errors.New("something went wrong"))

	// Check if results are ok
	fmt.Printf("Success is ok: %v\n", success.IsOk())
	fmt.Printf("Failure is ok: %v\n", failure.IsOk())

	// Extract values safely
	if value, err := success.Unwrap(); err == nil {
		fmt.Printf("Success value: %s\n", value)
	}

	if _, err := failure.Unwrap(); err != nil {
		fmt.Printf("Failure error: %v\n", err)
	}

	// Output:
	// Success is ok: true
	// Failure is ok: false
	// Success value: Hello, World!
	// Failure error: something went wrong
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Result_Map

ExampleResult_Map demonstrates transforming values inside Result.


<details><summary>Code</summary>

```go
func ExampleResult_Map() {
	// Start with a result containing a number
	result := Ok(5)

	// Transform to its square
	squared := result.Map(func(x int) int { return x * x })

	// Transform an error result
	errorResult := Err[int](errors.New("invalid input"))
	errorSquared := errorResult.Map(func(x int) int { return x * x })

	fmt.Printf("Square of 5: %v\n", squared.UnwrapOr(0))
	fmt.Printf("Square of error: %v\n", errorSquared.UnwrapOr(-1))

	// Output:
	// Square of 5: 25
	// Square of error: -1
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Result_Match

ExampleResult_Match demonstrates pattern matching with Result.


<details><summary>Code</summary>

```go
func ExampleResult_Match() {
	// Helper function that may fail
	divide := func(x, y int) Result[int] {
		if y == 0 {
			return Err[int](errors.New("division by zero"))
		}
		return Ok(x / y)
	}

	// Pattern match on results
	success := divide(10, 2)
	failure := divide(10, 0)

	result1 := success.Match(
		func(value int) Result[int] {
			return Ok(value * 2)
		},
		func(err error) Result[int] {
			return Err[int](fmt.Errorf("handled: %w", err))
		},
	)

	result2 := failure.Match(
		func(value int) Result[int] {
			return Ok(value * 2)
		},
		func(err error) Result[int] {
			return Err[int](fmt.Errorf("handled: %w", err))
		},
	)

	fmt.Printf("Success result: %v\n", result1.UnwrapOr(-1))
	fmt.Printf("Failure handled: %v\n", result2.IsErr())

	// Output:
	// Success result: 10
	// Failure handled: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Result_Or

ExampleResult_Or demonstrates providing fallback results.


<details><summary>Code</summary>

```go
func ExampleResult_Or() {
	// Create primary and fallback results
	primary := Err[string](errors.New("primary failed"))
	fallback := Ok("fallback value")

	// Use fallback when primary fails
	result := primary.Or(fallback)

	fmt.Printf("Result: %s\n", result.UnwrapOr("default"))

	// Output:
	// Result: fallback value
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### fp Some

ExampleSome demonstrates creating an Option with a value.


<details><summary>Code</summary>

```go
func ExampleSome() {
	// Create an Option containing a string
	message := Some("Hello, World!")

	fmt.Printf("Has value: %t\n", message.IsSome())
	fmt.Printf("Value: %s\n", message.UnwrapOr("default"))

	// Output:
	// Has value: true
	// Value: Hello, World!
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="lol"></a>📝 Lol

Package lol (lots of logs) provides a unified logging interface with multiple backends.

This package offers a simple, structured logging interface that can be backed by
different logging implementations. Currently it supports zerolog as the primary backend.

Key features:
- Structured logging with fields
- Multiple log levels (trace, debug, info, warn, error, fatal, panic)
- APM trace context integration
- Environment-aware configuration
- Testing utilities

Basic usage:

	logger := lol.NewZerologLogger(
		lol.Fields{"service": "myapp"},
		"production",
		"info",
		os.Stdout,
		lol.APMConfig{Enabled: true},
	)

	logger.Info("Application started")
	logger.WithField("user_id", 123).Warn("User action")

For testing:

	testLogger := lol.NewTest()
	testLogger.Error("This won't be printed")


### Functions

- [Logger_LogLevels](#lol-logger_loglevels)
- [Logger_WithField](#lol-logger_withfield)
- [NewTest](#lol-newtest)
- [NewZerologLogger](#lol-newzerologlogger)
- [ParseEnv](#lol-parseenv)
- [ParseLevel](#lol-parselevel)

#### lol Logger_LogLevels

ExampleLogger_LogLevels demonstrates different log levels


<details><summary>Code</summary>

```go
func ExampleLogger_LogLevels() {
	var buf bytes.Buffer

	logger := NewZerologLogger(
		Fields{"component": "auth"},
		"development",
		"trace", // Set to trace level to see all messages
		&buf,
		APMConfig{Enabled: false},
	)

	// Log at different levels
	logger.Trace("Entering authentication function")
	logger.Debug("Validating user credentials")
	logger.Info("User authentication successful")
	logger.Warn("Rate limit approaching")
	logger.Error("Authentication failed")

	output := buf.String()
	fmt.Printf(
		"Contains trace message: %t\n",
		bytes.Contains([]byte(output), []byte("Entering authentication")),
	)
	fmt.Printf(
		"Contains debug message: %t\n",
		bytes.Contains([]byte(output), []byte("Validating user")),
	)
	fmt.Printf(
		"Contains info message: %t\n",
		bytes.Contains([]byte(output), []byte("authentication successful")),
	)
	fmt.Printf("Contains warn message: %t\n", bytes.Contains([]byte(output), []byte("Rate limit")))
	fmt.Printf(
		"Contains error message: %t\n",
		bytes.Contains([]byte(output), []byte("Authentication failed")),
	)

	// Output:
	// Contains trace message: true
	// Contains debug message: true
	// Contains info message: true
	// Contains warn message: true
	// Contains error message: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### lol Logger_WithField

ExampleLogger_WithField demonstrates adding contextual fields to log messages


<details><summary>Code</summary>

```go
func ExampleLogger_WithField() {
	var buf bytes.Buffer

	logger := NewZerologLogger(
		Fields{"app": "demo"},
		"development",
		"debug",
		&buf,
		APMConfig{Enabled: false},
	)

	// Chain multiple fields
	enrichedLogger := logger.WithField("request_id", "req-123").
		WithField("user_agent", "test-client")
	enrichedLogger.Info("Processing request")

	// Add more context
	enrichedLogger.WithField("duration_ms", 45).Info("Request completed")

	output := buf.String()
	fmt.Printf(
		"Output contains 'request_id': %t\n",
		bytes.Contains([]byte(output), []byte("request_id")),
	)
	fmt.Printf(
		"Output contains 'user_agent': %t\n",
		bytes.Contains([]byte(output), []byte("user_agent")),
	)
	fmt.Printf(
		"Output contains 'duration_ms': %t\n",
		bytes.Contains([]byte(output), []byte("duration_ms")),
	)

	// Output:
	// Output contains 'request_id': true
	// Output contains 'user_agent': true
	// Output contains 'duration_ms': true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### lol NewTest

ExampleNewTest demonstrates creating a test logger that doesn't output anything


<details><summary>Code</summary>

```go
func ExampleNewTest() {
	// Create a test logger for unit tests
	testLogger := NewTest()

	// These messages won't be printed to stdout/stderr
	testLogger.Info("This is a test message")
	testLogger.Error("This error won't be shown")
	testLogger.WithField("test_field", "test_value").Debug("Debug message")

	fmt.Println("Test logger created successfully")
	fmt.Println("Messages logged silently for testing")

	// Output:
	// Test logger created successfully
	// Messages logged silently for testing
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### lol NewZerologLogger

ExampleNewZerologLogger demonstrates creating a structured logger with zerolog backend


<details><summary>Code</summary>

```go
func ExampleNewZerologLogger() {
	// Create a logger with custom fields and configuration
	var buf bytes.Buffer

	logger := NewZerologLogger(
		Fields{"service": "example-app", "version": "1.0.0"},
		"production",
		"info",
		&buf,
		APMConfig{Enabled: false}, // Disable APM for this example
	)

	// Log some messages
	logger.Info("Application started successfully")
	logger.WithField("user_id", 123).WithField("action", "login").Info("User logged in")
	logger.Warn("This is a warning message")

	fmt.Printf(
		"Logged output contains 'Application started': %t\n",
		bytes.Contains(buf.Bytes(), []byte("Application started")),
	)
	fmt.Printf(
		"Logged output contains 'user_id': %t\n",
		bytes.Contains(buf.Bytes(), []byte("user_id")),
	)
	fmt.Printf(
		"Logged output contains 'service': %t\n",
		bytes.Contains(buf.Bytes(), []byte("service")),
	)

	// Output:
	// Logged output contains 'Application started': true
	// Logged output contains 'user_id': true
	// Logged output contains 'service': true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### lol ParseEnv

ExampleParseEnv demonstrates parsing environment strings


<details><summary>Code</summary>

```go
func ExampleParseEnv() {
	// Parse different environments
	envs := []string{"test", "local", "development", "staging", "production"}

	for _, envStr := range envs {
		env := ParseEnv(envStr)
		fmt.Printf("Environment '%s' parsed as: %d\n", envStr, env)
	}

	// Output:
	// Environment 'test' parsed as: 0
	// Environment 'local' parsed as: 1
	// Environment 'development' parsed as: 2
	// Environment 'staging' parsed as: 3
	// Environment 'production' parsed as: 4
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### lol ParseLevel

ExampleParseLevel demonstrates parsing log levels from strings


<details><summary>Code</summary>

```go
func ExampleParseLevel() {
	// Parse different log levels
	levels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}

	for _, levelStr := range levels {
		level := ParseLevel(levelStr)
		fmt.Printf("Level '%s' parsed as: %d\n", levelStr, level)
	}

	// Output:
	// Level 'trace' parsed as: 6
	// Level 'debug' parsed as: 5
	// Level 'info' parsed as: 4
	// Level 'warn' parsed as: 3
	// Level 'error' parsed as: 2
	// Level 'fatal' parsed as: 1
	// Level 'panic' parsed as: 0
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="maps"></a>🗝️ Maps

Package maps provides generic utility functions to work with Go maps.
It offers a functional approach to common map operations like filtering, mapping,
reducing, and comparing maps.


### Functions

- [Equals](#maps-equals)
- [Filter](#maps-filter)
- [FilterMap](#maps-filtermap)
- [FilterMapTuple](#maps-filtermaptuple)
- [Fold](#maps-fold)
- [Map](#maps-map)
- [Reduce](#maps-reduce)
- [Slice](#maps-slice)

#### maps Equals

ExampleEquals demonstrates comparing two maps for equality.


<details><summary>Code</summary>

```go
func ExampleEquals() {
	// Create two maps
	map1 := map[string]int{"a": 1, "b": 2, "c": 3}
	map2 := map[string]int{"a": 1, "b": 2, "c": 3}
	map3 := map[string]int{"a": 1, "b": 2, "c": 4}

	// Compare using equality function
	equal1 := Equals(map1, map2, func(x, y int) bool { return x == y })
	equal2 := Equals(map1, map3, func(x, y int) bool { return x == y })

	fmt.Printf("map1 == map2: %t\n", equal1)
	fmt.Printf("map1 == map3: %t\n", equal2)
	// Output:
	// map1 == map2: true
	// map1 == map3: false
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps Filter

ExampleFilter demonstrates filtering a map by key-value pairs.


<details><summary>Code</summary>

```go
func ExampleFilter() {
	// Create a map of products to prices
	prices := map[string]int{
		"apple":  100,
		"banana": 50,
		"cherry": 200,
		"date":   75,
	}

	// Keep only items that cost more than 75
	expensive := Filter(prices, func(product string, price int) bool {
		return price > 75
	})

	fmt.Printf("Expensive items count: %d\n", len(expensive))
	// Output: Expensive items count: 2
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps FilterMap

ExampleFilterMap demonstrates filtering and transforming in a single operation.


<details><summary>Code</summary>

```go
func ExampleFilterMap() {
	// Create a map of names to ages
	ages := map[string]int{
		"Alice": 25,
		"Bob":   17,
		"Carol": 30,
		"Dave":  16,
	}

	// Keep only adults and transform to ID format
	adults := FilterMap(ages, func(name string, age int) fp.Option[tuples.Tuple2[string, string]] {
		if age >= 18 {
			id := fmt.Sprintf("ID_%s_%d", name, age)
			return fp.Some(tuples.Tuple2[string, string]{V1: name, V2: id})
		}
		return fp.None[tuples.Tuple2[string, string]]()
	})

	fmt.Printf("Adult count: %d\n", len(adults))
	// Output: Adult count: 2
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps FilterMapTuple

ExampleFilterMapTuple demonstrates filtering and transforming using tuple returns.


<details><summary>Code</summary>

```go
func ExampleFilterMapTuple() {
	// Create a map of scores
	scores := map[string]int{
		"Alice": 85,
		"Bob":   70,
		"Carol": 95,
		"Dave":  60,
	}

	// Keep high scores and convert to grade format
	grades := FilterMapTuple(scores, func(name string, score int) (string, string, bool) {
		if score >= 80 {
			var grade string
			if score >= 90 {
				grade = "A"
			} else {
				grade = "B"
			}
			return name, grade, true
		}
		return "", "", false
	})

	fmt.Printf("High performers: %d\n", len(grades))
	fmt.Printf("Alice's grade: %s\n", grades["Alice"])
	// Output:
	// High performers: 2
	// Alice's grade: B
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps Fold

ExampleFold demonstrates folding a map with an initial value.


<details><summary>Code</summary>

```go
func ExampleFold() {
	// Create a map of item prices
	prices := map[string]float64{
		"apple":  1.20,
		"banana": 0.80,
		"cherry": 2.50,
	}

	// Calculate total with initial tax
	totalWithTax := Fold(prices, func(acc float64, item string, price float64) float64 {
		return acc + price*1.1 // Add 10% tax
	}, 5.0) // Start with 5.0 base fee

	fmt.Printf("Total with tax: %.2f\n", totalWithTax)
	// Output:
	// Total with tax: 9.95
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps Map

ExampleMap demonstrates transforming keys and values in a map.


<details><summary>Code</summary>

```go
func ExampleMap() {
	// Create a map of numbers to their names
	numbers := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	// Transform to string keys and uppercase values
	transformed := Map(numbers, func(key int, value string) (string, string) {
		return fmt.Sprintf("num_%d", key), strings.ToUpper(value)
	})

	fmt.Println(transformed["num_1"])
	// Output: ONE
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps Reduce

ExampleReduce demonstrates reducing a map to a single value.


<details><summary>Code</summary>

```go
func ExampleReduce() {
	// Create a map of item quantities
	inventory := map[string]int{
		"apples":  10,
		"bananas": 5,
		"oranges": 8,
	}

	// Calculate total items (Reduce starts with zero value)
	total := Reduce(inventory, func(acc int, key string, value int) int {
		return acc + value
	})

	fmt.Printf("Total items: %d\n", total)
	// Output: Total items: 23
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### maps Slice

ExampleSlice demonstrates converting a map to a slice.


<details><summary>Code</summary>

```go
func ExampleSlice() {
	// Create a map of user data
	users := map[int]string{
		1: "Alice",
		2: "Bob",
		3: "Carol",
	}

	// Convert to slice of formatted strings
	userList := Slice(users, func(id int, name string) string {
		return fmt.Sprintf("ID:%d Name:%s", id, name)
	})

	fmt.Printf("Users count: %d\n", len(userList))
	// Note: map iteration order is not guaranteed
	// Output:
	// Users count: 3
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="num"></a>🔢 Num

Numeric utilities including high-precision decimal operations.

### Functions

- [Abs](#num-abs)
- [Dec_Add](#num-dec_add)
- [Dec_Compare](#num-dec_compare)
- [Dec_Div](#num-dec_div)
- [Dec_IsZero](#num-dec_iszero)
- [Dec_Mul](#num-dec_mul)
- [Dec_Percent](#num-dec_percent)
- [Dec_Round](#num-dec_round)
- [Dec_Sub](#num-dec_sub)
- [MustDecFromAny](#num-mustdecfromany)
- [MustDecFromString](#num-mustdecfromstring)
- [NewDecFromFloat](#num-newdecfromfloat)
- [NewDecFromInt](#num-newdecfromint)
- [NewDecFromString](#num-newdecfromstring)

#### num Abs

ExampleAbs demonstrates absolute value calculation


<details><summary>Code</summary>

```go
func ExampleAbs() {
	// Calculate absolute values for different types
	intVal := -42
	floatVal := -3.14

	absInt := Abs(intVal)
	absFloat := Abs(floatVal)

	fmt.Printf("Original int: %d, Absolute: %d\n", intVal, absInt)
	fmt.Printf("Original float: %.2f, Absolute: %.2f\n", floatVal, absFloat)

	// Output:
	// Original int: -42, Absolute: 42
	// Original float: -3.14, Absolute: 3.14
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Add

ExampleDec_Add demonstrates decimal addition


<details><summary>Code</summary>

```go
func ExampleDec_Add() {
	// Perform decimal addition
	price1 := MustDecFromString("123.45")
	price2 := MustDecFromString("67.89")

	total := price1.Add(price2)

	fmt.Printf("Price 1: %s\n", price1.String())
	fmt.Printf("Price 2: %s\n", price2.String())
	fmt.Printf("Total: %s\n", total.String())

	// Output:
	// Price 1: 123.45
	// Price 2: 67.89
	// Total: 191.34
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Compare

ExampleDec_Compare demonstrates decimal comparison


<details><summary>Code</summary>

```go
func ExampleDec_Compare() {
	// Compare decimal values
	price1 := MustDecFromString("99.99")
	price2 := MustDecFromString("100.00")
	price3 := MustDecFromString("99.99")

	fmt.Printf("Price 1: %s\n", price1.String())
	fmt.Printf("Price 2: %s\n", price2.String())
	fmt.Printf("Price 3: %s\n", price3.String())

	fmt.Printf("Price 1 < Price 2: %t\n", price1.LessThan(price2))
	fmt.Printf("Price 1 > Price 2: %t\n", price1.GreaterThan(price2))
	fmt.Printf("Price 1 == Price 3: %t\n", price1.Equal(price3))
	fmt.Printf("Price 1 <= Price 2: %t\n", price1.LessThanOrEqual(price2))

	// Output:
	// Price 1: 99.99
	// Price 2: 100
	// Price 3: 99.99
	// Price 1 < Price 2: true
	// Price 1 > Price 2: false
	// Price 1 == Price 3: true
	// Price 1 <= Price 2: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Div

ExampleDec_Div demonstrates decimal division


<details><summary>Code</summary>

```go
func ExampleDec_Div() {
	// Calculate unit price
	totalCost := MustDecFromString("127.50")
	quantity := MustDecFromString("25")

	unitPrice := totalCost.Div(quantity)

	fmt.Printf("Total cost: %s\n", totalCost.String())
	fmt.Printf("Quantity: %s\n", quantity.String())
	fmt.Printf("Unit price: %s\n", unitPrice.String())

	// Output:
	// Total cost: 127.5
	// Quantity: 25
	// Unit price: 5.1
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_IsZero

ExampleDec_IsZero demonstrates zero checking


<details><summary>Code</summary>

```go
func ExampleDec_IsZero() {
	// Check if decimal is zero
	zero := Zero
	nonZero := MustDecFromString("0.01")
	alsoZero := MustDecFromString("0.00")

	fmt.Printf("Zero value: %s, IsZero: %t\n", zero.String(), zero.IsZero())
	fmt.Printf("Non-zero value: %s, IsZero: %t\n", nonZero.String(), nonZero.IsZero())
	fmt.Printf("Also zero: %s, IsZero: %t\n", alsoZero.String(), alsoZero.IsZero())

	// Output:
	// Zero value: 0, IsZero: true
	// Non-zero value: 0.01, IsZero: false
	// Also zero: 0, IsZero: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Mul

ExampleDec_Mul demonstrates decimal multiplication


<details><summary>Code</summary>

```go
func ExampleDec_Mul() {
	// Calculate area
	length := MustDecFromString("12.5")
	width := MustDecFromString("8.4")

	area := length.Mul(width)

	fmt.Printf("Length: %s\n", length.String())
	fmt.Printf("Width: %s\n", width.String())
	fmt.Printf("Area: %s\n", area.String())

	// Output:
	// Length: 12.5
	// Width: 8.4
	// Area: 105
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Percent

ExampleDec_Percent demonstrates percentage calculations


<details><summary>Code</summary>

```go
func ExampleDec_Percent() {
	// Calculate percentage
	value := MustDecFromString("850.00")
	percentage := MustDecFromString("15") // 15%

	percentValue := value.ApplyPercent(percentage)
	finalValue := value.AddPercent(percentage)

	fmt.Printf("Original value: %s\n", value.String())
	fmt.Printf("Percentage: %s%%\n", percentage.String())
	fmt.Printf("Percentage amount: %s\n", percentValue.String())
	fmt.Printf("Final value (with percentage): %s\n", finalValue.String())

	// Output:
	// Original value: 850
	// Percentage: 15%
	// Percentage amount: 127.5
	// Final value (with percentage): 977.5
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Round

ExampleDec_Round demonstrates decimal rounding


<details><summary>Code</summary>

```go
func ExampleDec_Round() {
	// Round to different precision levels
	value := MustDecFromString("123.456789")

	rounded2 := value.RoundTo(2)
	rounded4 := value.RoundTo(4)
	rounded0 := value.RoundTo(0)

	fmt.Printf("Original: %s\n", value.String())
	fmt.Printf("Rounded to 2 decimals: %s\n", rounded2.String())
	fmt.Printf("Rounded to 4 decimals: %s\n", rounded4.String())
	fmt.Printf("Rounded to 0 decimals: %s\n", rounded0.String())

	// Output:
	// Original: 123.456789
	// Rounded to 2 decimals: 123.46
	// Rounded to 4 decimals: 123.4568
	// Rounded to 0 decimals: 123
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num Dec_Sub

ExampleDec_Sub demonstrates decimal subtraction


<details><summary>Code</summary>

```go
func ExampleDec_Sub() {
	// Perform decimal subtraction
	balance := MustDecFromString("1000.00")
	withdrawal := MustDecFromString("750.25")

	remaining := balance.Sub(withdrawal)

	fmt.Printf("Initial balance: %s\n", balance.String())
	fmt.Printf("Withdrawal: %s\n", withdrawal.String())
	fmt.Printf("Remaining: %s\n", remaining.String())

	// Output:
	// Initial balance: 1000
	// Withdrawal: 750.25
	// Remaining: 249.75
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num MustDecFromAny

ExampleMustDecFromAny demonstrates creating a decimal from any supported type


<details><summary>Code</summary>

```go
func ExampleMustDecFromAny() {
	// Create decimals from different types
	fromInt := MustDecFromAny(100)
	fromFloat := MustDecFromAny(3.14159)
	fromString := MustDecFromAny("42.42")

	fmt.Printf("From int: %s\n", fromInt.String())
	fmt.Printf("From float: %s\n", fromFloat.String())
	fmt.Printf("From string: %s\n", fromString.String())

	// Output:
	// From int: 100
	// From float: 3.14159
	// From string: 42.42
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num MustDecFromString

ExampleMustDecFromString demonstrates creating a decimal from a string (panics on error)


<details><summary>Code</summary>

```go
func ExampleMustDecFromString() {
	// Create decimal from valid string
	price := MustDecFromString("999.99")

	fmt.Printf("Price: %s\n", price.String())
	if f, ok := price.Float64(); ok {
		fmt.Printf("As float: %.2f\n", f)
	}

	// Output:
	// Price: 999.99
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num NewDecFromFloat

ExampleNewDecFromFloat demonstrates creating a decimal from a float


<details><summary>Code</summary>

```go
func ExampleNewDecFromFloat() {
	// Create decimal from float
	temperature := NewDecFromFloat(36.5)

	fmt.Printf("Temperature: %s\n", temperature.String())
	if f, ok := temperature.Float64(); ok {
		fmt.Printf("As float: %.1f\n", f)
	}
	fmt.Printf("Number of decimals: %d\n", temperature.NumberOfDecimals())

	// Output:
	// Temperature: 36.5
	// As float: 36.5
	// Number of decimals: 1
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num NewDecFromInt

ExampleNewDecFromInt demonstrates creating a decimal from an integer


<details><summary>Code</summary>

```go
func ExampleNewDecFromInt() {
	// Create decimal from integer
	quantity := NewDecFromInt(42)

	fmt.Printf("Quantity: %s\n", quantity.String())
	fmt.Printf("As int: %d\n", quantity.IntPart())
	fmt.Printf("Is zero: %t\n", quantity.IsZero())

	// Output:
	// Quantity: 42
	// As int: 42
	// Is zero: false
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### num NewDecFromString

ExampleNewDecFromString demonstrates creating a decimal from a string


<details><summary>Code</summary>

```go
func ExampleNewDecFromString() {
	// Create decimal from string
	price, err := NewDecFromString("123.456")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Price: %s\n", price.String())
	fmt.Printf("Is set: %t\n", price.Isset())

	// Try with invalid string
	_, err = NewDecFromString("invalid")
	fmt.Printf("Invalid string error: %v\n", err != nil)

	// Output:
	// Price: 123.456
	// Is set: true
	// Invalid string error: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="slices"></a>⛓️ Slices

Package slices provides a comprehensive set of generic utility functions for working with slices.
It offers a functional approach to common slice operations such as transforming, filtering,
searching, and manipulating elements in a type-safe manner.


### Functions

- [All](#slices-all)
- [Contains](#slices-contains)
- [Filter](#slices-filter)
- [FilterMap](#slices-filtermap)
- [Find](#slices-find)
- [Fold](#slices-fold)
- [Map](#slices-map)
- [Reduce](#slices-reduce)
- [Some](#slices-some)
- [ToMap](#slices-tomap)

#### slices All

ExampleAll demonstrates checking if all elements satisfy a condition.


<details><summary>Code</summary>

```go
func ExampleAll() {
	// Create a slice of positive numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Check if all numbers are positive
	allPositive := All(numbers, func(n int) bool {
		return n > 0
	})

	// Check if all numbers are even
	allEven := All(numbers, func(n int) bool {
		return n%2 == 0
	})

	fmt.Printf("All positive: %t\n", allPositive)
	fmt.Printf("All even: %t\n", allEven)
	// Output:
	// All positive: true
	// All even: false
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Contains

ExampleContains demonstrates checking if any element satisfies a condition.


<details><summary>Code</summary>

```go
func ExampleContains() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Check if any number is greater than 3
	hasLarge := Contains(numbers, func(n int) bool {
		return n > 3
	})

	// Check if any number is negative
	hasNegative := Contains(numbers, func(n int) bool {
		return n < 0
	})

	fmt.Printf("Has number > 3: %t\n", hasLarge)
	fmt.Printf("Has negative: %t\n", hasNegative)
	// Output:
	// Has number > 3: true
	// Has negative: false
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Filter

ExampleFilter demonstrates filtering a slice to keep only elements that satisfy a condition.


<details><summary>Code</summary>

```go
func ExampleFilter() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter to keep only even numbers
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	fmt.Println(evenNumbers)
	// Output: [2 4 6 8 10]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices FilterMap

ExampleFilterMap demonstrates filtering and transforming in a single operation.


<details><summary>Code</summary>

```go
func ExampleFilterMap() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5, 6}

	// Keep only even numbers and square them
	evenSquares := FilterMap(numbers, func(n int) fp.Option[int] {
		if n%2 == 0 {
			return fp.Some(n * n)
		}
		return fp.None[int]()
	})

	fmt.Println(evenSquares)
	// Output: [4 16 36]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Find

ExampleFind demonstrates finding the first element that matches a condition.


<details><summary>Code</summary>

```go
func ExampleFind() {
	// Create a slice of names
	names := []string{"Alice", "Bob", "Charlie", "David"}

	// Find the first name that starts with 'C'
	result, found := Find(names, func(name string) bool {
		return len(name) > 0 && name[0] == 'C'
	})

	fmt.Printf("Found: %t, Name: %s\n", found, result)
	// Output: Found: true, Name: Charlie
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Fold

ExampleFold demonstrates folding a slice with an initial value.


<details><summary>Code</summary>

```go
func ExampleFold() {
	// Create a slice of strings
	words := []string{"Hello", "World", "from", "Go"}

	// Join with custom separator and prefix
	result := Fold(words, func(acc, word string) string {
		if acc == "" {
			return "Greeting: " + word
		}
		return acc + " " + word
	}, "")

	fmt.Printf("Result: %s\n", result)
	// Output:
	// Result: Greeting: Hello World from Go
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Map

ExampleMap demonstrates transforming elements in a slice.


<details><summary>Code</summary>

```go
func ExampleMap() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Transform each number by squaring it
	squares := Map(numbers, func(n int) int {
		return n * n
	})

	fmt.Println(squares)
	// Output: [1 4 9 16 25]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Reduce

ExampleReduce demonstrates combining all elements into a single value.


<details><summary>Code</summary>

```go
func ExampleReduce() {
	// Create a slice of numbers
	numbers := []int{1, 2, 3, 4, 5}

	// Sum all numbers
	sum := Reduce[int, int](numbers, func(acc, curr int) int {
		return acc + curr
	})

	fmt.Println(sum)
	// Output: 15
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices Some

ExampleSome demonstrates checking if some elements satisfy a condition.


<details><summary>Code</summary>

```go
func ExampleSome() {
	// Create a slice of words
	words := []string{"hello", "world", "go", "programming"}

	// Check if some words are short (< 4 characters)
	hasShort := Some(words, func(word string) bool {
		return len(word) < 4
	})

	fmt.Printf("Has short words: %t\n", hasShort)
	// Output:
	// Has short words: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### slices ToMap

ExampleToMap demonstrates converting a slice to a map using a key function.


<details><summary>Code</summary>

```go
func ExampleToMap() {
	// Create a slice of words
	words := []string{"apple", "banana", "cherry"}

	// Convert to map with first letter as key
	wordMap := ToMap(words, func(word string) rune {
		return rune(word[0])
	})

	fmt.Printf("'a' word: %s\n", wordMap['a'])
	fmt.Printf("'b' word: %s\n", wordMap['b'])
	// Output:
	// 'a' word: apple
	// 'b' word: banana
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="streams"></a>🌊 Streams

Package streams provides interfaces and types for reading and writing streams of data.


### Functions

- [Batch](#streams-batch)
- [CSV](#streams-csv)
- [CSVTransform](#streams-csvtransform)
- [CSVTransform_tabSeparated](#streams-csvtransform_tabseparated)
- [ConsumeErrSkip](#streams-consumeerrskip)
- [DatabaseStream](#streams-databasestream)
- [Filter](#streams-filter)
- [FilterMap](#streams-filtermap)
- [Flatten](#streams-flatten)
- [Group](#streams-group)
- [JSON](#streams-json)
- [JSONEachRowTransform](#streams-jsoneachrowtransform)
- [JSONTransform](#streams-jsontransform)
- [Lines](#streams-lines)
- [Map](#streams-map)
- [MemWriter](#streams-memwriter)
- [Multicast](#streams-multicast)
- [Pipe](#streams-pipe)
- [PipeCSV](#streams-pipecsv)
- [PipeJSON](#streams-pipejson)
- [PipeJSONEachRow](#streams-pipejsoneachrow)
- [Reader](#streams-reader)
- [Reduce](#streams-reduce)
- [ReduceMap](#streams-reducemap)
- [ReduceSlice](#streams-reduceslice)
- [WriteAll](#streams-writeall)

#### streams Batch

ExampleBatch demonstrates grouping stream elements into batches.


<details><summary>Code</summary>

```go
func ExampleBatch() {
	// Create a stream from a slice of integers
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream := MemReader(data, nil)

	// Group into batches of 3
	batchStream := Batch(stream, 3)

	// Collect the results
	result, _ := Consume(batchStream)
	for i, batch := range result {
		fmt.Printf("Batch %d: %v\n", i+1, batch)
	}
	// Output:
	// Batch 1: [1 2 3]
	// Batch 2: [4 5 6]
	// Batch 3: [7 8 9]
	// Batch 4: [10]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams CSV

ExampleCSV demonstrates reading CSV data from a string.


<details><summary>Code</summary>

```go
func ExampleCSV() {
	// Create a CSV reader from a string
	csvData := "name,age,city\nAlice,25,NYC\nBob,30,LA\nCharlie,35,Chicago"
	reader := io.NopCloser(strings.NewReader(csvData))

	// Create a CSV stream directly
	csvStream, _ := CSV[[]string](
		WithCSVReader(reader),
		WithCSVSeparator(","),
	)

	// Collect the results
	result, _ := Consume(csvStream)
	for i, row := range result {
		fmt.Printf("Row %d: %v\n", i+1, row)
	}
	// Output:
	// Row 1: [name age city]
	// Row 2: [Alice 25 NYC]
	// Row 3: [Bob 30 LA]
	// Row 4: [Charlie 35 Chicago]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams CSVTransform

ExampleCSVTransform demonstrates converting a stream to CSV format.


<details><summary>Code</summary>

```go
func ExampleCSVTransform() {
	// Create a stream of employees
	employees := []Employee{
		{ID: 1, Name: "Alice Johnson", Department: "Engineering", Salary: 75000.00},
		{ID: 2, Name: "Bob Smith", Department: "Marketing", Salary: 65000.00},
		{ID: 3, Name: "Charlie Brown", Department: "Engineering", Salary: 80000.00},
	}
	stream := MemReader(employees, nil)

	// Transform to CSV with comma separator
	transform := CSVTransform(stream, CSVSeparatorComma)
	transform.WriteTo(os.Stdout)

	// Output:
	// ID,Name,Department,Salary
	// 1,Alice Johnson,Engineering,75000.00
	// 2,Bob Smith,Marketing,65000.00
	// 3,Charlie Brown,Engineering,80000.00
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams CSVTransform_tabSeparated

ExampleCSVTransform_tabSeparated demonstrates CSV with tab separator.


<details><summary>Code</summary>

```go
func ExampleCSVTransform_tabSeparated() {

	// Create a stream of products
	products := []ProductCSV{
		{SKU: "LAPTOP-001", Name: "Gaming Laptop", Price: 1299.99},
		{SKU: "MOUSE-002", Name: "Wireless Mouse", Price: 49.99},
		{SKU: "KEYBOARD-003", Name: "Mechanical Keyboard", Price: 129.99},
	}
	stream := MemReader(products, nil)

	// Transform to CSV with tab separator
	transform := CSVTransform(stream, CSVSeparatorTab)
	transform.WriteTo(os.Stdout)

	// Output:
	// SKU	Product Name	Price
	// LAPTOP-001	Gaming Laptop	$1299.99
	// MOUSE-002	Wireless Mouse	$49.99
	// KEYBOARD-003	Mechanical Keyboard	$129.99
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams ConsumeErrSkip

ExampleConsumeErrSkip demonstrates consuming a stream while skipping errors.


<details><summary>Code</summary>

```go
func ExampleConsumeErrSkip() {
	// Create a filter stream that may produce errors
	reader := strings.NewReader("1\n2\ninvalid\n4\n5")
	numbersStream := Lines(reader)

	// Create a filter that converts strings to numbers (may fail)
	filterStream := FilterMap(numbersStream, func(s string) (int, bool) {
		// Simulate conversion that might fail
		if s == "invalid" {
			return 0, false // This will be skipped
		}
		// Simple conversion for demonstration
		switch s {
		case "1":
			return 1, true
		case "2":
			return 2, true
		case "4":
			return 4, true
		case "5":
			return 5, true
		default:
			return 0, false
		}
	})

	// Consume all valid numbers, skipping errors
	numbers := ConsumeErrSkip(filterStream)

	fmt.Printf("Valid numbers: %v\n", numbers)

	// Output:
	// Valid numbers: [1 2 4 5]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams DatabaseStream

ExampleStream demonstrates how to use Stream with database rows.


<details><summary>Code</summary>

```go
func ExampleDatabaseStream() {
	// Mock data that simulates database rows
	mockData := &mockRows{
		data: [][]any{
			{1, "Alice"},
			{2, "Bob"},
			{3, "Charlie"},
		},
	}

	// Create a stream with a scan function
	stream := DB(mockData, func(rows DBRows, user *User) error {
		return rows.Scan(&user.ID, &user.Name)
	})

	// Iterate through the stream
	for stream.Next() {
		user := stream.Data()
		fmt.Printf("User ID: %d, Name: %s\n", user.ID, user.Name)
	}

	// Check for errors
	if err := stream.Err(); err != nil {
		log.Printf("Error during streaming: %v", err)
	}

	// Output:
	// User ID: 1, Name: Alice
	// User ID: 2, Name: Bob
	// User ID: 3, Name: Charlie
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Filter

ExampleFilter demonstrates filtering a stream of integers to keep only even numbers.


<details><summary>Code</summary>

```go
func ExampleFilter() {
	// Create a stream from a slice of integers
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream := MemReader(data, nil)

	// Filter to keep only even numbers
	evenStream := Filter(stream, func(n int) bool {
		return n%2 == 0
	})

	// Collect the results
	result, _ := Consume(evenStream)
	fmt.Println(result)
	// Output: [2 4 6 8 10]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams FilterMap

ExampleFilterMap demonstrates filtering and transforming in a single operation.


<details><summary>Code</summary>

```go
func ExampleFilterMap() {
	// Create a stream from a slice of integers
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream := MemReader(data, nil)

	// Filter even numbers and convert them to strings
	evenStrings := FilterMap(stream, func(n int) (string, bool) {
		if n%2 == 0 {
			return strconv.Itoa(n), true
		}
		return "", false
	})

	// Collect the results
	result, _ := Consume(evenStrings)
	fmt.Println(result)
	// Output: [2 4 6 8 10]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Flatten

ExampleFlatten demonstrates flattening a stream of slices.


<details><summary>Code</summary>

```go
func ExampleFlatten() {
	// Create a stream from a slice of slices
	data := [][]int{{1, 2}, {3, 4, 5}, {6}, {7, 8, 9}}
	stream := MemReader(data, nil)

	// Flatten the slices
	flattened := Flatten(stream)

	// Collect the results
	result, _ := Consume(flattened)
	fmt.Println(result)
	// Output: [1 2 3 4 5 6 7 8 9]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Group

ExampleGroup demonstrates grouping consecutive items with the same key.


<details><summary>Code</summary>

```go
func ExampleGroup() {
	// Create a stream from a slice of strings
	data := []string{"apple", "apricot", "banana", "blueberry", "cherry", "coconut"}
	stream := MemReader(data, nil)

	// Group by first letter
	Grouped := Group(stream, func(s string) rune {
		return rune(s[0])
	})

	// Collect the results
	result, _ := Consume(Grouped)
	for _, group := range result {
		fmt.Printf("Group: %v\n", group)
	}
	// Output:
	// Group: [apple apricot]
	// Group: [banana blueberry]
	// Group: [cherry coconut]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams JSON

ExampleJSON demonstrates reading JSON data line by line.


<details><summary>Code</summary>

```go
func ExampleJSON() {
		// Create JSON-lines data (each line is a separate JSON object)
		jsonData := `{"name":"Alice","age":25}
	{"name":"Bob","age":30}
	{"name":"Charlie","age":35}`

		reader := strings.NewReader(jsonData)

		// Create a JSON stream for a simple struct
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		jsonStream := JSON[Person](io.NopCloser(reader))

		// Collect the results
		result, _ := Consume(jsonStream)
		for _, person := range result {
			fmt.Printf("Person: %s, Age: %d\n", person.Name, person.Age)
		}
		// Output:
		// Person: Alice, Age: 25
		// Person: Bob, Age: 30
		// Person: Charlie, Age: 35
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams JSONEachRowTransform

ExampleJSONEachRowTransform demonstrates converting a stream to JSON lines format.


<details><summary>Code</summary>

```go
func ExampleJSONEachRowTransform() {
	// Define a simple structure for demonstration
	type LogEntry struct {
		Timestamp string `json:"timestamp"`
		Level     string `json:"level"`
		Message   string `json:"message"`
	}

	// Create a stream of log entries
	logs := []LogEntry{
		{Timestamp: "2025-06-28T10:00:00Z", Level: "INFO", Message: "Application started"},
		{Timestamp: "2025-06-28T10:01:00Z", Level: "WARN", Message: "High memory usage detected"},
		{Timestamp: "2025-06-28T10:02:00Z", Level: "ERROR", Message: "Database connection failed"},
	}
	stream := MemReader(logs, nil)

	// Transform to JSON lines format and write to stdout
	transform := JSONEachRowTransform(stream)
	transform.WriteTo(os.Stdout)

	// Output:
	// {"timestamp":"2025-06-28T10:00:00Z","level":"INFO","message":"Application started"}
	// {"timestamp":"2025-06-28T10:01:00Z","level":"WARN","message":"High memory usage detected"}
	// {"timestamp":"2025-06-28T10:02:00Z","level":"ERROR","message":"Database connection failed"}
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams JSONTransform

ExampleJSONTransform demonstrates converting a stream to JSON array format.


<details><summary>Code</summary>

```go
func ExampleJSONTransform() {
	// Define a simple structure for demonstration
	type Person struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Create a stream from a slice of persons
	people := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	stream := MemReader(people, nil)

	// Transform to JSON and write to stdout
	transform := JSONTransform(stream)
	transform.WriteTo(os.Stdout)

	// Output:
	// [{"id":1,"name":"Alice"},{"id":2,"name":"Bob"},{"id":3,"name":"Charlie"}]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Lines

ExampleLines demonstrates reading lines from a string.


<details><summary>Code</summary>

```go
func ExampleLines() {
	// Create a reader from a multiline string
	text := "line1\nline2\nline3\n"
	reader := strings.NewReader(text)

	// Create a lines stream
	lineStream := Lines(reader)

	// Collect the results
	result, _ := Consume(lineStream)
	fmt.Println(result)
	// Output: [line1 line2 line3]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Map

ExampleMap demonstrates transforming elements in a stream.


<details><summary>Code</summary>

```go
func ExampleMap() {
	// Create a stream from a slice of integers
	data := []int{1, 2, 3, 4, 5}
	stream := MemReader(data, nil)

	// Transform integers to their string representation
	stringStream := Map(stream, func(n int) string {
		return fmt.Sprintf("number_%d", n)
	})

	// Collect the results
	result, _ := Consume(stringStream)
	fmt.Println(result)
	// Output: [number_1 number_2 number_3 number_4 number_5]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams MemWriter

ExampleMemWriter demonstrates writing items to memory.


<details><summary>Code</summary>

```go
func ExampleMemWriter() {
	// Create a memory writer for strings
	writer := MemWriter[string]()

	// Write some items
	items := []string{"hello", "world", "from", "memory"}
	for _, item := range items {
		writer.Write(item)
	}

	// Get all items
	result := writer.Items()

	fmt.Printf("Items written: %d\n", len(result))
	fmt.Printf("Items: %v\n", result)
	// Output:
	// Items written: 4
	// Items: [hello world from memory]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Multicast

ExampleMulticast demonstrates broadcasting a stream to multiple destinations.


<details><summary>Code</summary>

```go
func ExampleMulticast() {
	// Create a stream of numbers
	reader := strings.NewReader("1\n2\n3\n4\n5")
	source := Lines(reader)

	// Create two memory writers to collect data separately
	dest1 := MemWriter[string]()
	dest2 := MemWriter[string]()

	// Multicast the stream to both destinations
	counts, err := Multicast(source, dest1, dest2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Written to dest1: %d items\n", counts[0])
	fmt.Printf("Written to dest2: %d items\n", counts[1])
	fmt.Printf("Dest1 data: %v\n", dest1.Items())
	fmt.Printf("Dest2 data: %v\n", dest2.Items())

	// Output:
	// Written to dest1: 5 items
	// Written to dest2: 5 items
	// Dest1 data: [1 2 3 4 5]
	// Dest2 data: [1 2 3 4 5]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Pipe

ExamplePipe demonstrates piping data from one stream to another.


<details><summary>Code</summary>

```go
func ExamplePipe() {
	// Create a source stream
	data := []string{"hello", "world", "from", "streams"}
	source := MemReader(data, nil)

	// Create a destination
	dest := MemWriter[string]()

	// Pipe data from source to destination
	bytesWritten, _ := Pipe(source, dest)

	fmt.Printf("Items written: %d\n", bytesWritten)
	fmt.Printf("Items: %v\n", dest.Items())
	// Output:
	// Items written: 4
	// Items: [hello world from streams]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams PipeCSV

ExamplePipeCSV demonstrates using the PipeCSV convenience function.


<details><summary>Code</summary>

```go
func ExamplePipeCSV() {
	// Create a stream of employees
	employees := []Employee{
		{ID: 101, Name: "Diana Prince", Department: "Legal", Salary: 90000.00},
		{ID: 102, Name: "Clark Kent", Department: "Journalism", Salary: 55000.00},
	}
	stream := MemReader(employees, nil)

	// Use PipeCSV to write directly to stdout with comma separator
	rowsWritten, err := PipeCSV(stream, os.Stdout, CSVSeparatorComma)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Rows written: %d\n", rowsWritten)

	// Output:
	// ID,Name,Department,Salary
	// 101,Diana Prince,Legal,90000.00
	// 102,Clark Kent,Journalism,55000.00
	// Rows written: 3
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams PipeJSON

ExamplePipeJSON demonstrates using the PipeJSON convenience function.


<details><summary>Code</summary>

```go
func ExamplePipeJSON() {
	// Define a simple structure
	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	// Create a stream of products
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99},
		{ID: 2, Name: "Mouse", Price: 29.99},
	}
	stream := MemReader(products, nil)

	// Use PipeJSON to write directly to stdout
	bytesWritten, err := PipeJSON(stream, os.Stdout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\nBytes written: %d\n", bytesWritten)

	// Output:
	// [{"id":1,"name":"Laptop","price":999.99},{"id":2,"name":"Mouse","price":29.99}]
	// Bytes written: 79
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams PipeJSONEachRow

ExamplePipeJSONEachRow demonstrates using the PipeJSONEachRow convenience function.


<details><summary>Code</summary>

```go
func ExamplePipeJSONEachRow() {
	// Define a simple metric structure
	type Metric struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	}

	// Create a stream of metrics
	metrics := []Metric{
		{Name: "cpu_usage", Value: 85.5, Unit: "percent"},
		{Name: "memory_usage", Value: 1024, Unit: "MB"},
		{Name: "disk_usage", Value: 75.2, Unit: "percent"},
	}
	stream := MemReader(metrics, nil)

	// Use PipeJSONEachRow to write to stdout
	bytesWritten, err := PipeJSONEachRow(stream, os.Stdout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Bytes written: %d\n", bytesWritten)

	// Output:
	// {"name":"cpu_usage","value":85.5,"unit":"percent"}
	// {"name":"memory_usage","value":1024,"unit":"MB"}
	// {"name":"disk_usage","value":75.2,"unit":"percent"}
	// Bytes written: 152
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Reader

ExampleReader demonstrates reading byte chunks from an io.Reader.


<details><summary>Code</summary>

```go
func ExampleReader() {
	// Create data to read
	data := "line1\nline2\nline3\n"
	reader := strings.NewReader(data)

	// Create a byte stream
	stream := Reader(reader)

	// Read all chunks
	var chunks []string
	for stream.Next() {
		chunks = append(chunks, string(stream.Data()))
	}

	fmt.Printf("Chunks: %d\n", len(chunks))
	fmt.Printf("First chunk: %q\n", chunks[0])
	// Output:
	// Chunks: 3
	// First chunk: "line1\n"
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams Reduce

ExampleReduce demonstrates using Reduce to sum numbers from a stream.


<details><summary>Code</summary>

```go
func ExampleReduce() {
	// Create a stream of numbers from strings
	reader := strings.NewReader("10\n20\n30\n40\n50")
	lines := Lines(reader)

	// Convert strings to numbers and sum them
	sum, _ := Reduce(Map(lines, func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}), func(acc, n int) int {
		return acc + n
	}, 0)

	fmt.Printf("Sum: %d\n", sum)

	// Output:
	// Sum: 150
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams ReduceMap

ExampleReduceMap demonstrates reducing a stream to a map with aggregated values.


<details><summary>Code</summary>

```go
func ExampleReduceMap() {
	// Create a stream of words
	reader := strings.NewReader("apple\nbanana\napple\ncherry\nbanana\napple")
	stream := Lines(reader)

	// Count occurrences of each word
	counts, _ := ReduceMap(stream, func(acc map[string]int, word string) map[string]int {
		acc[word]++
		return acc
	})

	fmt.Printf("apple: %d\n", counts["apple"])
	fmt.Printf("banana: %d\n", counts["banana"])
	fmt.Printf("cherry: %d\n", counts["cherry"])

	// Output:
	// apple: 3
	// banana: 2
	// cherry: 1
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams ReduceSlice

ExampleReduceSlice demonstrates collecting filtered items from a stream.


<details><summary>Code</summary>

```go
func ExampleReduceSlice() {
	// Create a stream of words
	reader := strings.NewReader("cat\ndog\nelephant\nant\nbutterfly\nbird")
	stream := Lines(reader)

	// Collect only words longer than 3 characters
	longWords, _ := ReduceSlice(stream, func(acc []string, word string) []string {
		if len(word) > 3 {
			return append(acc, word)
		}
		return acc
	})

	fmt.Printf("Long words: %v\n", longWords)

	// Output:
	// Long words: [elephant butterfly bird]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### streams WriteAll

ExampleWriteAll demonstrates writing a slice to a stream.


<details><summary>Code</summary>

```go
func ExampleWriteAll() {
	// Create data to write
	data := []string{"hello", "world", "streams"}

	// Create a memory writer
	writer := MemWriter[string]()

	// Write all data
	bytesWritten, err := WriteAll(writer, data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Bytes written: %d\n", bytesWritten)
	fmt.Printf("Items: %v\n", writer.Items())
	// Output:
	// Bytes written: 3
	// Items: [hello world streams]
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="zero"></a>🔞 Zero

Zero-value utilities and string manipulation functions.

### Functions

- [B2S](#zero-b2s)
- [S2B](#zero-s2b)
- [ZeroAllocConversions](#zero-zeroallocconversions)

#### zero B2S

ExampleB2S demonstrates converting []byte to string without memory allocation


<details><summary>Code</summary>

```go
func ExampleB2S() {
	// Convert []byte to string using zero-allocation conversion
	b := []byte("Hello, Gophers!")
	s := B2S(b)

	fmt.Printf("Original bytes: %v\n", b)
	fmt.Printf("Converted to string: %s\n", s)
	fmt.Printf("Length: %d\n", len(s))
	fmt.Printf(
		"Same underlying data: %t\n",
		uintptr(unsafe.Pointer(&b[0])) == uintptr(unsafe.Pointer(unsafe.StringData(s))),
	)

	// Output:
	// Original bytes: [72 101 108 108 111 44 32 71 111 112 104 101 114 115 33]
	// Converted to string: Hello, Gophers!
	// Length: 15
	// Same underlying data: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### zero S2B

ExampleS2B demonstrates converting a string to []byte without memory allocation


<details><summary>Code</summary>

```go
func ExampleS2B() {
	// Convert string to []byte using zero-allocation conversion
	s := "Hello, World!"
	b := S2B(s)

	fmt.Printf("Original string: %s\n", s)
	fmt.Printf("Converted to bytes: %v\n", b)
	fmt.Printf("Bytes as string: %s\n", string(b))
	fmt.Printf(
		"Same underlying data: %t\n",
		uintptr(unsafe.Pointer(unsafe.StringData(s))) == uintptr(unsafe.Pointer(&b[0])),
	)

	// Output:
	// Original string: Hello, World!
	// Converted to bytes: [72 101 108 108 111 44 32 87 111 114 108 100 33]
	// Bytes as string: Hello, World!
	// Same underlying data: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---

#### zero ZeroAllocConversions

ExampleZeroAllocConversions demonstrates the performance benefits of zero-allocation conversions


<details><summary>Code</summary>

```go
func ExampleZeroAllocConversions() {
	// Traditional conversion (allocates memory)
	original := "Performance matters!"
	traditionalBytes := []byte(original)
	traditionalString := string(traditionalBytes)

	// Zero-allocation conversion (shares memory)
	zeroAllocBytes := S2B(original)
	zeroAllocString := B2S(zeroAllocBytes)

	fmt.Printf("Original: %s\n", original)
	fmt.Printf("Traditional conversion: %s\n", traditionalString)
	fmt.Printf("Zero-alloc conversion: %s\n", zeroAllocString)
	fmt.Printf(
		"All results equal: %t\n",
		original == traditionalString && traditionalString == zeroAllocString,
	)

	// Output:
	// Original: Performance matters!
	// Traditional conversion: Performance matters!
	// Zero-alloc conversion: Performance matters!
	// All results equal: true
}
```

</details>


[⬆️ Back to Top](#table-of-contents)

---


[⬆️ Back to Top](#table-of-contents)


<br/>

