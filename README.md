
# gozo

[![Go Report Card](https://goreportcard.com/badge/github.com/sonirico/gozo)](https://goreportcard.com/report/github.com/sonirico/gozo)
[![Go Reference](https://pkg.go.dev/badge/github.com/sonirico/gozo.svg)](https://pkg.go.dev/github.com/sonirico/gozo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/sonirico/gozo.svg)](https://github.com/sonirico/gozo/releases)

[![gozo Art](gozo.png)](https://github.com/sonirico/gozo/gozo.png)

The ultimate toolkit for Go developers. A comprehensive collection of functions, data structures, and utilities designed to enhance productivity and code quality.

## Modules

## <a name="table-of-contents"></a>Table of Contents

- [🪄 Fp](#fp)
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
- [🗝️ Maps](#maps)
  - [Equals](#maps-equals)
  - [Filter](#maps-filter)
  - [FilterMap](#maps-filtermap)
  - [FilterMapTuple](#maps-filtermaptuple)
  - [Fold](#maps-fold)
  - [Map](#maps-map)
  - [Reduce](#maps-reduce)
  - [Slice](#maps-slice)
- [⛓️ Slices](#slices)
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
- [🌊 Streams](#streams)
  - [Batch](#streams-batch)
  - [CSV](#streams-csv)
  - [Compact](#streams-compact)
  - [ConsumeErrSkip](#streams-consumeerrskip)
  - [Filter](#streams-filter)
  - [FilterMap](#streams-filtermap)
  - [Flatten](#streams-flatten)
  - [JSON](#streams-json)
  - [Lines](#streams-lines)
  - [Map](#streams-map)
  - [MemWriter](#streams-memwriter)
  - [Multicast](#streams-multicast)
  - [Pipe](#streams-pipe)
  - [Reader](#streams-reader)
  - [Reduce](#streams-reduce)
  - [WriteAll](#streams-writeall)

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


---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="streams"></a>🌊 Streams

Powerful data streaming and processing utilities with fluent API for functional programming patterns.

### Functions

- [Batch](#streams-batch)
- [CSV](#streams-csv)
- [Compact](#streams-compact)
- [ConsumeErrSkip](#streams-consumeerrskip)
- [Filter](#streams-filter)
- [FilterMap](#streams-filtermap)
- [Flatten](#streams-flatten)
- [JSON](#streams-json)
- [Lines](#streams-lines)
- [Map](#streams-map)
- [MemWriter](#streams-memwriter)
- [Multicast](#streams-multicast)
- [Pipe](#streams-pipe)
- [Reader](#streams-reader)
- [Reduce](#streams-reduce)
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


---

#### streams Compact

ExampleCompact demonstrates grouping consecutive items with the same key.


<details><summary>Code</summary>

```go
func ExampleCompact() {
	// Create a stream from a slice of strings
	data := []string{"apple", "apricot", "banana", "blueberry", "cherry", "coconut"}
	stream := MemReader(data, nil)

	// Group by first letter
	compacted := Compact(stream, func(s string) rune {
		return rune(s[0])
	})

	// Collect the results
	result, _ := Consume(compacted)
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


---

#### streams Reduce

ExampleReduce demonstrates reducing a stream to a map with aggregated values.


<details><summary>Code</summary>

```go
func ExampleReduce() {
	// Create a stream of words
	reader := strings.NewReader("apple\nbanana\napple\ncherry\nbanana\napple")
	stream := Lines(reader)

	// Count occurrences of each word
	counts, _ := Reduce(stream, func(acc map[string]int, word string) map[string]int {
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


---


[⬆️ Back to Top](#table-of-contents)


<br/>

