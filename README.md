
# gozo

[![Go Report Card](https://goreportcard.com/badge/github.com/sonirico/gozo)](https://goreportcard.com/report/github.com/sonirico/gozo)
[![Go Reference](https://pkg.go.dev/badge/github.com/sonirico/gozo.svg)](https://pkg.go.dev/github.com/sonirico/gozo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![gozo Art](gozo.png)](https://github.com/sonirico/gozo/gozo.png)

The ultimate toolkit for Go developers. A comprehensive collection of functions, data structures, and utilities designed to enhance productivity and code quality.

## Modules

## <a name="table-of-contents"></a>Table of Contents

- [🪄 Fp](#fp)
  - [Err](#err)
  - [None](#none)
  - [Ok](#ok)
  - [OkZero](#okzero)
  - [OptionFromPtr](#optionfromptr)
  - [OptionFromTuple](#optionfromtuple)
  - [OptionFromZero](#optionfromzero)
  - [Some](#some)
- [🗝️ Maps](#maps)
  - [Equals](#equals)
  - [Filter](#filter)
  - [FilterInPlace](#filterinplace)
  - [FilterMap](#filtermap)
  - [FilterMapTuple](#filtermaptuple)
  - [Fold](#fold)
  - [Map](#map)
  - [Reduce](#reduce)
  - [Slice](#slice)
- [⛓️ Slices](#slices)
  - [All](#all)
  - [Any](#any)
  - [Append](#append)
  - [AppendVector](#appendvector)
  - [Contains](#contains)
  - [Cut](#cut)
  - [Delete](#delete)
  - [DeleteOrder](#deleteorder)
  - [Equals](#equals)
  - [Extract](#extract)
  - [ExtractIdx](#extractidx)
  - [Filter](#filter)
  - [FilterInPlace](#filterinplace)
  - [FilterInPlaceCopy](#filterinplacecopy)
  - [FilterMap](#filtermap)
  - [FilterMapTuple](#filtermaptuple)
  - [Find](#find)
  - [FindIdx](#findidx)
  - [Fold](#fold)
  - [FoldSame](#foldsame)
  - [Includes](#includes)
  - [IndexOf](#indexof)
  - [Insert](#insert)
  - [InsertVector](#insertvector)
  - [Map](#map)
  - [MapInPlace](#mapinplace)
  - [Peek](#peek)
  - [Pop](#pop)
  - [PopFront](#popfront)
  - [PushFront](#pushfront)
  - [Reduce](#reduce)
  - [ReduceSame](#reducesame)
  - [Shift](#shift)
  - [Some](#some)
  - [ToMap](#tomap)
  - [ToMapIdx](#tomapidx)
  - [Unshift](#unshift)
- [🌊 Streams](#streams)
  - [CSVTransform](#csvtransform)
  - [Collect](#collect)
  - [Iter](#iter)
  - [Iter2](#iter2)
  - [JSONEachRowTransform](#jsoneachrowtransform)
  - [JSONTransform](#jsontransform)
  - [Multiplex](#multiplex)
  - [NewBatchStream](#newbatchstream)
  - [NewCompactStream](#newcompactstream)
  - [NewCompactStreamFactory](#newcompactstreamfactory)
  - [NewFilterStream](#newfilterstream)
  - [NewFlattenerStream](#newflattenerstream)
  - [NewLineReaderStream](#newlinereaderstream)
  - [NewMapperStream](#newmapperstream)
  - [NewMemory](#newmemory)
  - [NewMemoryWriteStream](#newmemorywritestream)
  - [NewReaderStream](#newreaderstream)
  - [NewWriterStream](#newwriterstream)
  - [Pipe](#pipe)
  - [SeqKeys](#seqkeys)
  - [SeqValues](#seqvalues)
  - [WriteAll](#writeall)
  - [WriteSeq](#writeseq)
  - [WriteSeqKeys](#writeseqkeys)
  - [WriteSeqValues](#writeseqvalues)

## <a name="fp"></a>🪄 Fp

Package fp provides functional programming primitives for Go.
It implements monadic types like Option and Result for more expressive error handling.


### Functions

- [Err](#err)
- [None](#none)
- [Ok](#ok)
- [OkZero](#okzero)
- [OptionFromPtr](#optionfromptr)
- [OptionFromTuple](#optionfromtuple)
- [OptionFromZero](#optionfromzero)
- [Some](#some)

#### <a name="err"></a>Err

Err creates a new Result in the error state with the given error.
This is a constructor function for creating a Result that represents failure.


<details><summary>Code</summary>

```go
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}
```

</details>


---

#### <a name="none"></a>None

None creates a new Option in the None state.
This is a constructor function for creating an Option that does not contain a value.


<details><summary>Code</summary>

```go
func None[T any]() Option[T] {
	return Option[T]{}
}
```

</details>


---

#### <a name="ok"></a>Ok

Ok creates a new Result in the Ok state with the given value.
This is a constructor function for creating a Result that represents success.


<details><summary>Code</summary>

```go
func Ok[T any](v T) Result[T] {
	return Result[T]{value: v, err: nil}
}
```

</details>


---

#### <a name="okzero"></a>OkZero

OkZero creates a new Result in the Ok state with the zero value.
This is a constructor function for creating a Result that represents success
but doesn't carry a meaningful value.


<details><summary>Code</summary>

```go
func OkZero[T any]() Result[T] {
	return Result[T]{err: nil}
}
```

</details>


---

#### <a name="optionfromptr"></a>OptionFromPtr

OptionFromPtr creates an Option from a pointer.
If the pointer is nil, returns None, otherwise returns Some with the dereferenced value.
This is useful for converting nullable pointers to the Option type.


<details><summary>Code</summary>

```go
func OptionFromPtr[T any](x *T) Option[T] {
	if x == nil {
		return None[T]()
	}
	return Some(*x)
}
```

</details>


---

#### <a name="optionfromtuple"></a>OptionFromTuple

OptionFromTuple creates an Option from a tuple-like return (value, ok).
If ok is true, returns Some(x), otherwise returns None.
This is useful for converting Go's common (value, ok) pattern to an Option.


<details><summary>Code</summary>

```go
func OptionFromTuple[T any](x T, ok bool) Option[T] {
	if ok {
		return Some(x)
	}
	return None[T]()
}
```

</details>


---

#### <a name="optionfromzero"></a>OptionFromZero

OptionFromZero creates an Option from a value, treating zero values as None.
If the value equals the zero value for its type, returns None, otherwise returns Some(x).
This is useful when zero values are treated as invalid or unset.


<details><summary>Code</summary>

```go
func OptionFromZero[T comparable](x T) Option[T] {
	var zero T
	if x == zero {
		return None[T]()
	}
	return Some(x)
}
```

</details>


---

#### <a name="some"></a>Some

Some creates a new Option in the Some state with the given value.
This is a constructor function for creating an Option that contains a value.


<details><summary>Code</summary>

```go
func Some[T any](t T) Option[T] {
	return Option[T]{value: t, isSome: true}
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

- [Equals](#equals)
- [Filter](#filter)
- [FilterInPlace](#filterinplace)
- [FilterMap](#filtermap)
- [FilterMapTuple](#filtermaptuple)
- [Fold](#fold)
- [Map](#map)
- [Reduce](#reduce)
- [Slice](#slice)

#### <a name="equals"></a>Equals

Equals compares two maps and returns whether they are equal in values.
Two maps are considered equal if:
- They have the same length
- They contain the same keys
- For each key, the values in both maps satisfy the equality function

Maps are compared using the provided equality function for values.
This allows for deep equality checks on complex value types.


<details><summary>Code</summary>

```go
func Equals[K comparable, V any](m1, m2 map[K]V, eq func(V, V) bool) bool {
	if len(m1) != len(m2) {
		return false
	}

	if m1 == nil && m2 != nil {
		return false
	}

	if m1 != nil && m2 == nil {
		return false
	}

	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok {
			return false
		}

		if !eq(v1, v2) {
			return false
		}
	}

	return true
}
```

</details>


---

#### <a name="filter"></a>Filter

Filter creates a new map containing only the key-value pairs that satisfy the predicate.
The predicate function takes a key and value and returns a boolean indicating
whether to include the entry in the result.

Unlike FilterInPlace, this function creates a new map and does not modify the input map.


<details><summary>Code</summary>

```go
func Filter[K comparable, V any](
	m map[K]V,
	p func(K, V) bool,
) map[K]V {
	if m == nil {
		return nil
	}

	res := make(map[K]V, len(m))

	for k, v := range m {
		if p(k, v) {
			res[k] = v
		}
	}

	return res
}
```

</details>


---

#### <a name="filterinplace"></a>FilterInPlace

FilterInPlace modifies the given map by removing entries that do not satisfy the predicate.
The predicate function takes a key and value and returns a boolean indicating
whether to keep the entry in the map.

This function directly modifies the input map for better performance when
creating a new map is not necessary.
It returns the modified map for convenience in chaining operations.


<details><summary>Code</summary>

```go
func FilterInPlace[K comparable, V any](
	m map[K]V,
	p func(K, V) bool,
) map[K]V {
	if m == nil {
		return nil
	}

	for k, v := range m {
		if !p(k, v) {
			delete(m, k)
		}
	}

	return m
}
```

</details>


---

#### <a name="filtermap"></a>FilterMap

FilterMap both filters and maps a map into a new map, potentially with different key and value types.
The predicate function should return an fp.Option monad containing a tuple of the new key and value:
- fp.Some to include the entry in the result (with transformed key and value)
- fp.None to exclude the entry from the result

This provides a powerful way to simultaneously transform and filter map entries
while leveraging the Option monad for expressing presence/absence.


<details><summary>Code</summary>

```go
func FilterMap[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	p func(K1, V1) fp.Option[tuples.Tuple2[K2, V2]],
) map[K2]V2 {
	if m == nil {
		return nil
	}

	res := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		tpl := p(k1, v1)
		if tpl.IsSome() {
			v := tpl.UnwrapUnsafe()
			res[v.V1] = v.V2
		}
	}

	return res
}
```

</details>


---

#### <a name="filtermaptuple"></a>FilterMapTuple

FilterMapTuple both filters and maps the given map into a new map, potentially with different key and value types.
The predicate function returns three values:
- The new key (K2)
- The new value (V2)
- A boolean indicating whether to include this entry in the result

This function is an alternative to FilterMap that uses Go's native boolean return
instead of the Option monad for expressing presence/absence.


<details><summary>Code</summary>

```go
func FilterMapTuple[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	p func(K1, V1) (K2, V2, bool),
) map[K2]V2 {
	if m == nil {
		return nil
	}

	res := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		if k2, v2, ok := p(k1, v1); ok {
			res[k2] = v2
		}
	}

	return res
}
```

</details>


---

#### <a name="fold"></a>Fold

Fold compacts a map into a single value by iteratively applying a reduction function
with an explicit initial value.
The reduction function takes the accumulator, a key, and a value, and returns
the updated accumulator.

Unlike Reduce, Fold takes an explicit initial value for the accumulator.
This is useful when the zero value of the result type is not appropriate
as the starting value.


<details><summary>Code</summary>

```go
func Fold[K comparable, V any, R any](
	m map[K]V,
	p func(R, K, V) R,
	initial R,
) R {
	if m == nil {
		return initial
	}

	r := initial

	for k, v := range m {
		r = p(r, k, v)
	}

	return r
}
```

</details>


---

#### <a name="map"></a>Map

Map transforms a map into another map, with potentially different key and value types.
The transformation is applied to each key-value pair by the provided function,
which returns the new key and value for the resulting map.

This function preserves nil semantics: if the input map is nil, the output will also be nil.
Otherwise, a new map is created with the transformed key-value pairs.


<details><summary>Code</summary>

```go
func Map[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	p func(K1, V1) (K2, V2),
) map[K2]V2 {
	if m == nil {
		return nil
	}

	res := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		k2, v2 := p(k1, v1)
		res[k2] = v2
	}

	return res
}
```

</details>


---

#### <a name="reduce"></a>Reduce

Reduce compacts a map into a single value by iteratively applying a reduction function.
The reduction function takes the accumulator, a key, and a value, and returns
the updated accumulator.

The initial value for the accumulator is the zero value of type R.
If you need a different initial value, use Fold instead.


<details><summary>Code</summary>

```go
func Reduce[K comparable, V any, R any](
	m map[K]V,
	p func(R, K, V) R,
) R {
	var r R

	if m == nil {
		return r
	}

	for k, v := range m {
		r = p(r, k, v)
	}

	return r
}
```

</details>


---

#### <a name="slice"></a>Slice

Slice converts a map into a slice by applying a transformation function to each key-value pair.
The transformation function takes a key and value and returns an element
for the resulting slice.

The order of elements in the resulting slice is not guaranteed, as map iteration
in Go is not deterministic.


<details><summary>Code</summary>

```go
func Slice[K comparable, V, R any](
	m map[K]V,
	p func(K, V) R,
) slices.Slice[R] {
	res := make([]R, len(m))
	i := 0

	for k, v := range m {
		res[i] = p(k, v)
		i++
	}

	return res
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

- [All](#all)
- [Any](#any)
- [Append](#append)
- [AppendVector](#appendvector)
- [Contains](#contains)
- [Cut](#cut)
- [Delete](#delete)
- [DeleteOrder](#deleteorder)
- [Equals](#equals)
- [Extract](#extract)
- [ExtractIdx](#extractidx)
- [Filter](#filter)
- [FilterInPlace](#filterinplace)
- [FilterInPlaceCopy](#filterinplacecopy)
- [FilterMap](#filtermap)
- [FilterMapTuple](#filtermaptuple)
- [Find](#find)
- [FindIdx](#findidx)
- [Fold](#fold)
- [FoldSame](#foldsame)
- [Includes](#includes)
- [IndexOf](#indexof)
- [Insert](#insert)
- [InsertVector](#insertvector)
- [Map](#map)
- [MapInPlace](#mapinplace)
- [Peek](#peek)
- [Pop](#pop)
- [PopFront](#popfront)
- [PushFront](#pushfront)
- [Reduce](#reduce)
- [ReduceSame](#reducesame)
- [Shift](#shift)
- [Some](#some)
- [ToMap](#tomap)
- [ToMapIdx](#tomapidx)
- [Unshift](#unshift)

#### <a name="all"></a>All

All checks if all elements in the slice satisfy the predicate.
Returns true if all elements match the predicate, false otherwise.


<details><summary>Code</summary>

```go
func All[T any](arr []T, predicate func(t T) bool) bool {
	for _, x := range arr {
		if !predicate(x) {
			return false
		}
	}
	return true
}
```

</details>


---

#### <a name="any"></a>Any

Any checks if at least one element in the slice satisfies the predicate.
Returns true if any element matches the predicate, false otherwise.
Alias for Contains.


<details><summary>Code</summary>

```go
func Any[T any](arr []T, predicate func(t T) bool) bool {
	return Contains(arr, predicate)
}
```

</details>


---

#### <a name="append"></a>Append

Append adds an element to the end of the slice and returns the result.
Unlike the builtin append, this function is explicitly named for clarity.


<details><summary>Code</summary>

```go
func Append[T any](arr []T, item T) []T {
	return append(arr, item)
}
```

</details>


---

#### <a name="appendvector"></a>AppendVector

AppendVector adds all elements from another slice to the end of this slice.
Returns the resulting concatenated slice.


<details><summary>Code</summary>

```go
func AppendVector[T any](arr, items []T) []T {
	return append(arr, items...)
}
```

</details>


---

#### <a name="contains"></a>Contains

Contains checks if the slice contains an element that satisfies the predicate.
Returns true if any element matches the predicate, false otherwise.


<details><summary>Code</summary>

```go
func Contains[T any](arr []T, predicate func(t T) bool) bool {
	return IndexOf(arr, predicate) >= 0
}
```

</details>


---

#### <a name="cut"></a>Cut

Cut removes a sector from slice given lower and upper bounds. Bounds are
represented as indices of the slice. E.g:
Cut([1, 2, 3, 4], 1, 2) -> [1, 4]
Cut([4], 0, 0) -> []
Cut will returned the original slice without the cut subslice.


<details><summary>Code</summary>

```go
func Cut[T any](arr []T, from, to int) []T {
	if len(arr) < 1 {
		return arr
	}

	if from < 0 {
		from = 0
	}

	if from >= len(arr) {
		from = len(arr) - 1
	}

	if to < 0 {
		to = 0
	}

	if to >= len(arr) {
		to = len(arr) - 1
	}

	if len(arr) == 1 {
		return arr[:0]
	}

	if from > to {

		return append(arr[:from], arr[from+to+1:]...)
	}

	return append(arr[:from], arr[to+1:]...)
}
```

</details>


---

#### <a name="delete"></a>Delete

Delete removes the element at the specified index without preserving order.
This provides better performance than DeleteOrder but changes the order of elements.
If the index is out of bounds, returns the original slice unchanged.


<details><summary>Code</summary>

```go
func Delete[T any](arr []T, idx int) []T {
	le := len(arr) - 1
	if le < 0 || idx > le || idx < 0 {
		return arr
	}
	var t T
	arr[idx] = arr[le]
	arr[le] = t
	arr = arr[:le]
	return arr
}
```

</details>


---

#### <a name="deleteorder"></a>DeleteOrder

DeleteOrder removes the element at the specified index while preserving order.
This is slower than Delete but maintains the relative order of the remaining elements.
If the index is out of bounds, returns the original slice unchanged.


<details><summary>Code</summary>

```go
func DeleteOrder[T any](arr []T, idx int) []T {
	le := len(arr) - 1
	if le < 0 || idx > le || idx < 0 {
		return arr
	}
	var t T

	if le > 0 {
		copy(arr[idx:], arr[idx+1:])
	}

	arr[le] = t
	arr = arr[:le]
	return arr
}
```

</details>


---

#### <a name="equals"></a>Equals

Equals compares two slices and returns whether they contain equal elements.
Two slices are considered equal if they have the same length and corresponding
elements satisfy the equality function.


<details><summary>Code</summary>

```go
func Equals[T any](one, other []T, predicate func(x, y T) bool) (res bool) {
	if len(one) != len(other) {
		return
	}

	res = true

	for idx, otherItem := range other {
		res = predicate(one[idx], otherItem)
		if !res {
			return
		}
	}
	return
}
```

</details>


---

#### <a name="extract"></a>Extract

Extract gets and deletes the first element that matches the predicate.
Returns the modified slice, the extracted element, and a success flag.
If no element matches, returns the original slice, zero value, and false.


<details><summary>Code</summary>

```go
func Extract[T any](arr []T, predicate func(t T) bool) ([]T, T, bool) {
	res, idx := FindIdx(arr, predicate)
	if idx < 0 {
		return arr, res, false
	}

	arr = Delete(arr, idx)
	return arr, res, true
}
```

</details>


---

#### <a name="extractidx"></a>ExtractIdx

ExtractIdx gets and deletes the element at the given position.
Returns the modified slice, the extracted element, and a success flag.
If the index is out of bounds, returns the original slice, zero value, and false.


<details><summary>Code</summary>

```go
func ExtractIdx[T any](arr []T, idx int) (res []T, item T, ok bool) {
	if idx >= len(arr) || idx < 0 {
		return
	}

	ok = true
	item = arr[idx]
	res = Delete(arr, idx)

	return
}
```

</details>


---

#### <a name="filter"></a>Filter

Filter creates a new slice containing only the elements that satisfy the predicate.


<details><summary>Code</summary>

```go
func Filter[T any](arr []T, predicate func(t T) bool) []T {
	res := make([]T, 0, len(arr))

	for _, x := range arr {
		if predicate(x) {
			res = append(res, x)
		}
	}

	return res
}
```

</details>


---

#### <a name="filterinplace"></a>FilterInPlace

FilterInPlace modifies the slice in place to contain only elements that
satisfy the predicate. This is more efficient than Filter when creating
a new slice is not necessary.


<details><summary>Code</summary>

```go
func FilterInPlace[T any](arr []T, predicate func(t T) bool) []T {
	n := 0
	for i, x := range arr {
		if predicate(x) {
			if n != i {
				arr[n] = x
			}
			n++
		}
	}

	arr = arr[:n]

	return arr
}
```

</details>


---

#### <a name="filterinplacecopy"></a>FilterInPlaceCopy

FilterInPlaceCopy filters the slice in place and returns a copy of the result.
This combines the efficiency of FilterInPlace with the safety of creating a new slice.


<details><summary>Code</summary>

```go
func FilterInPlaceCopy[T any](arr []T, predicate func(t T) bool) []T {
	n := 0
	for i, x := range arr {
		if predicate(x) {
			if n != i {
				arr[n] = x
			}
			n++
		}
	}

	arr = arr[:n]

	res := make([]T, n)

	copy(res, arr[:n])

	return res
}
```

</details>


---

#### <a name="filtermap"></a>FilterMap

FilterMap creates a new slice by applying a transformation function that
returns an Option. Elements with Some options are included in the result,
while None options are excluded.


<details><summary>Code</summary>

```go
func FilterMap[T, U any](arr []T, predicate func(t T) fp.Option[U]) []U {
	pre := func(t T) (U, bool) {
		return predicate(t).Unwrap()
	}

	return FilterMapTuple[T, U](arr, pre)
}
```

</details>


---

#### <a name="filtermaptuple"></a>FilterMapTuple

FilterMapTuple creates a new slice by applying a transformation function that
also filters elements. The function should return the transformed value and
a boolean indicating whether to include the element.


<details><summary>Code</summary>

```go
func FilterMapTuple[T, U any](arr []T, predicate func(t T) (U, bool)) []U {
	res := make([]U, 0, len(arr))

	for _, x := range arr {
		if mapped, ok := predicate(x); ok {
			res = append(res, mapped)
		}
	}

	return res
}
```

</details>


---

#### <a name="find"></a>Find

Find returns the first element that satisfies the predicate.
Returns the element and true if found, otherwise the zero value and false.


<details><summary>Code</summary>

```go
func Find[T any](arr []T, predicate func(t T) bool) (res T, ok bool) {
	var idx int
	res, idx = FindIdx(arr, predicate)
	ok = idx > -1
	return
}
```

</details>


---

#### <a name="findidx"></a>FindIdx

FindIdx returns the first element that satisfies the predicate and its index.
Returns the element and its index if found, otherwise the zero value and -1.


<details><summary>Code</summary>

```go
func FindIdx[T any](arr []T, predicate func(t T) bool) (res T, idx int) {
	idx = IndexOf(arr, predicate)
	if idx < 0 {
		return
	}

	res = arr[idx]
	return
}
```

</details>


---

#### <a name="fold"></a>Fold

Fold compacts the slice into a single value by iteratively applying
the reduction function, starting with the provided initial value.
The accumulator type can be different from the element type.


<details><summary>Code</summary>

```go
func Fold[T, U any](arr []T, p func(U, T) U, initial U) U {
	if len(arr) < 1 {
		return initial
	}

	initial = p(initial, arr[0])

	if len(arr) < 2 {
		return initial
	}

	i := 1

	for i < len(arr) {
		initial = p(initial, arr[i])

		i++
	}

	return initial
}
```

</details>


---

#### <a name="foldsame"></a>FoldSame

FoldSame is a convenience wrapper around Fold for when the accumulator
and element types are the same.


<details><summary>Code</summary>

```go
func FoldSame[T any](arr []T, p func(T, T) T, initial T) T {
	return Fold[T, T](arr, p, initial)
}
```

</details>


---

#### <a name="includes"></a>Includes

Includes checks if the slice contains a specific element using the equality operator.
Returns true if the element is found, false otherwise.


<details><summary>Code</summary>

```go
func Includes[T comparable](arr []T, target T) bool {
	return Contains(arr, func(t T) bool {
		return t == target
	})
}
```

</details>


---

#### <a name="indexof"></a>IndexOf

IndexOf returns the index of the first element that satisfies the predicate.
Returns the index where the element was found, or -1 if not found.


<details><summary>Code</summary>

```go
func IndexOf[T any](arr []T, predicate func(t T) bool) (pos int) {
	pos = -1
	for i, x := range arr {
		if predicate(x) {
			pos = i
			return
		}
	}
	return
}
```

</details>


---

#### <a name="insert"></a>Insert

Insert places an element at the specified index in the slice.
Elements at or after the index are shifted to the right.
Returns the resulting slice with the new element inserted.
If the index is out of bounds, returns the original slice unchanged.


<details><summary>Code</summary>

```go
func Insert[T any](arr []T, item T, idx int) []T {
	if arr == nil {
		return []T{item}
	}

	if idx < 0 || idx > len(arr) {
		return arr
	}

	return append(arr[:idx], append([]T{item}, arr[idx:]...)...)
}
```

</details>


---

#### <a name="insertvector"></a>InsertVector

InsertVector places a slice of elements at the specified index in the slice.
Elements at or after the index are shifted to the right.
Returns the resulting slice with the new elements inserted.
If the index is out of bounds, returns the original slice unchanged.


<details><summary>Code</summary>

```go
func InsertVector[T any](arr, items []T, idx int) (res []T) {
	if arr == nil {
		res = items[:]
		return
	}

	if items == nil || len(items) == 0 {
		res = arr
		return
	}

	if idx < 0 || idx > len(arr) {
		return arr
	}

	return append(arr[:idx], append(items, arr[idx:]...)...)
}
```

</details>


---

#### <a name="map"></a>Map

Map creates a new slice by applying the transformation function to each element.
The transformation can change the type of the elements.


<details><summary>Code</summary>

```go
func Map[T, U any](arr []T, predicate func(t T) U) []U {
	res := make([]U, 0, len(arr))

	for _, x := range arr {
		res = append(res, predicate(x))
	}

	return res
}
```

</details>


---

#### <a name="mapinplace"></a>MapInPlace

MapInPlace transforms each element in the slice using the provided function.
Modifies the slice in place and returns it.


<details><summary>Code</summary>

```go
func MapInPlace[T any](arr []T, predicate func(t T) T) []T {
	for i, x := range arr {
		arr[i] = predicate(x)
	}

	return arr
}
```

</details>


---

#### <a name="peek"></a>Peek

Peek returns the item at the specified index without modifying the slice.
Returns the element and true if the index is valid, otherwise the zero value and false.


<details><summary>Code</summary>

```go
func Peek[T any](arr []T, idx int) (item T, ok bool) {
	if len(arr) < 1 || idx >= len(arr) {
		return
	}

	item = arr[idx]
	ok = true

	return
}
```

</details>


---

#### <a name="pop"></a>Pop

Pop deletes and returns the last item from the slice.
Returns the modified slice, the popped element, and a success flag.
If the slice is empty, returns the original slice, zero value, and false.


<details><summary>Code</summary>

```go
func Pop[T any](arr []T) (res []T, item T, ok bool) {
	if len(arr) < 1 {
		return
	}

	var t T
	le := len(arr) - 1
	res = arr[:le]
	item = arr[le]
	ok = true

	arr[le] = t

	return
}
```

</details>


---

#### <a name="popfront"></a>PopFront

PopFront removes and returns the first element of the slice.
Returns the modified slice (without the first element), the removed element, and a success flag.
If the slice is empty, returns the original slice, zero value, and false.


<details><summary>Code</summary>

```go
func PopFront[T any](arr []T) (res []T, item T, ok bool) {
	if len(arr) < 1 {
		res = arr
		return
	}

	item, res = arr[0], arr[1:]
	return
}
```

</details>


---

#### <a name="pushfront"></a>PushFront

PushFront inserts an element at the beginning of the slice.
Returns the resulting slice with the new element at the front.


<details><summary>Code</summary>

```go
func PushFront[T any](arr []T, item T) []T {
	return append([]T{item}, arr...)
}
```

</details>


---

#### <a name="reduce"></a>Reduce

Reduce compacts the slice into a single value by iteratively applying
the reduction function to each element. Starts with the zero value.


<details><summary>Code</summary>

```go
func Reduce[T, U any](arr []T, p func(T, T) T) (res T) {
	return Fold(arr, p, res)
}
```

</details>


---

#### <a name="reducesame"></a>ReduceSame

ReduceSame is a convenience wrapper around Reduce for when the accumulator
and element types are the same.


<details><summary>Code</summary>

```go
func ReduceSame[T any](arr []T, p func(T, T) T) T {
	return Reduce[T, T](arr, p)
}
```

</details>


---

#### <a name="shift"></a>Shift

Shift removes and returns the first element of the slice.
Alias for PopFront, following JavaScript array method naming conventions.


<details><summary>Code</summary>

```go
func Shift[T any](arr []T) ([]T, T, bool) {
	return PopFront(arr)
}
```

</details>


---

#### <a name="some"></a>Some

Some checks if at least one element in the slice satisfies the predicate.
Returns true if any element matches the predicate, false otherwise.
Alias for Contains.


<details><summary>Code</summary>

```go
func Some[T any](arr []T, predicate func(t T) bool) bool {
	return Contains(arr, predicate)
}
```

</details>


---

#### <a name="tomap"></a>ToMap

ToMap creates a map from a slice, using the provided function to determine the key
for each element. The element itself becomes the value in the map.


<details><summary>Code</summary>

```go
func ToMap[V any, K comparable](arr []V, predicate func(x V) K) map[K]V {
	res := make(map[K]V, len(arr))

	for _, x := range arr {
		res[predicate(x)] = x
	}

	return res
}
```

</details>


---

#### <a name="tomapidx"></a>ToMapIdx

ToMapIdx creates a map from a slice, preserving each element's original index.
Uses the provided function to determine the key for each element.
The value in the map is a WrappedIdx containing both the element and its original index.


<details><summary>Code</summary>

```go
func ToMapIdx[V any, K comparable](arr []V, predicate func(x V) K) map[K]WrappedIdx[V] {
	res := make(map[K]WrappedIdx[V], len(arr))

	for i, x := range arr {
		res[predicate(x)] = WrappedIdx[V]{value: x, idx: i}
	}

	return res
}
```

</details>


---

#### <a name="unshift"></a>Unshift

Unshift inserts an element at the beginning of the slice.
Alias for PushFront, following JavaScript array method naming conventions.


<details><summary>Code</summary>

```go
func Unshift[T any](arr []T, item T) []T {
	return PushFront(arr, item)
}
```

</details>


---


[⬆️ Back to Top](#table-of-contents)


<br/>

## <a name="streams"></a>🌊 Streams



### Functions

- [CSVTransform](#csvtransform)
- [Collect](#collect)
- [Iter](#iter)
- [Iter2](#iter2)
- [JSONEachRowTransform](#jsoneachrowtransform)
- [JSONTransform](#jsontransform)
- [Multiplex](#multiplex)
- [NewBatchStream](#newbatchstream)
- [NewCompactStream](#newcompactstream)
- [NewCompactStreamFactory](#newcompactstreamfactory)
- [NewFilterStream](#newfilterstream)
- [NewFlattenerStream](#newflattenerstream)
- [NewLineReaderStream](#newlinereaderstream)
- [NewMapperStream](#newmapperstream)
- [NewMemory](#newmemory)
- [NewMemoryWriteStream](#newmemorywritestream)
- [NewReaderStream](#newreaderstream)
- [NewWriterStream](#newwriterstream)
- [Pipe](#pipe)
- [SeqKeys](#seqkeys)
- [SeqValues](#seqvalues)
- [WriteAll](#writeall)
- [WriteSeq](#writeseq)
- [WriteSeqKeys](#writeseqkeys)
- [WriteSeqValues](#writeseqvalues)

#### <a name="csvtransform"></a>CSVTransform

CSVTransform creates a new PipeCSVTransform for a given stream.


<details><summary>Code</summary>

```go
func CSVTransform[T csvMarshaler](stream ReadStream[T], separator rune) Transform[T] {
	return &TransformCSV[T]{
		stream:    stream,
		separator: separator,
	}
}
```

</details>


---

#### <a name="collect"></a>Collect

Collect collects all elements from an iter.Seq into a slice.


<details><summary>Code</summary>

```go
func Collect[T any](stream iter.Seq[T]) []T {
	return slices.Collect(stream)
}
```

</details>


---

#### <a name="iter"></a>Iter

Iter converts a ReadStream into an iter.Seq.


<details><summary>Code</summary>

```go
func Iter[T any](stream ReadStream[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer stream.Close()

		for stream.Next() {
			if !yield(stream.Data()) {
				break
			}
		}

		return
	}
}
```

</details>


---

#### <a name="iter2"></a>Iter2

Iter2 converts a ReadStream into an iter.Seq2, yielding both data and error.


<details><summary>Code</summary>

```go
func Iter2[T any](stream ReadStream[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		defer stream.Close()

		for stream.Next() {
			if !yield(stream.Data(), nil) {
				break
			}
		}
		if err := stream.Err(); err != nil {
			var x T
			yield(x, err)
		}
		if err := stream.Close(); err != nil {
			var x T
			yield(x, err)
		}
	}
}
```

</details>


---

#### <a name="jsoneachrowtransform"></a>JSONEachRowTransform

JSONEachRowTransform creates a Transform that converts a ReadStream to JSON-lines format.
Each element in the stream becomes a separate JSON object on its own line,
which is useful for streaming JSON processing and log formats.

This format is also known as NDJSON (Newline Delimited JSON) and is commonly
used for streaming APIs and log processing systems.


<details><summary>Code</summary>

```go
func JSONEachRowTransform[T any](stream ReadStream[T]) Transform[T] {
	return &TransformJSONEachRow[T]{
		stream: stream,
	}
}
```

</details>


---

#### <a name="jsontransform"></a>JSONTransform

JSONTransform creates a Transform that converts a ReadStream to JSON format.
The resulting Transform can be used with WriteTo to output the stream data
as a JSON array to any io.Writer.

This is useful for converting structured data to JSON for APIs, file output,
or network transmission.


<details><summary>Code</summary>

```go
func JSONTransform[T any](r ReadStream[T]) Transform[T] {
	return &TransformJSON[T]{
		stream: r,
	}
}
```

</details>


---

#### <a name="multiplex"></a>Multiplex

Multiplex copies all items from a ReadStream to multiple WriteStreams
Returns a slice with bytes written to each destination and any error


<details><summary>Code</summary>

```go
func Multiplex[T any](src ReadStream[T], destinations ...WriteStream[T]) ([]int64, error) {
	if len(destinations) == 0 {
		return []int64{}, nil
	}

	bytesWritten := make([]int64, len(destinations))

	for src.Next() {
		if err := src.Err(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return bytesWritten, fmt.Errorf("read error: %w", err)
		}

		data := src.Data()
		for i, dst := range destinations {
			n, err := dst.Write(data)
			if err != nil {
				return bytesWritten, fmt.Errorf("write error to destination %d: %w", i, err)
			}
			bytesWritten[i] += n
		}
	}

	for i, dst := range destinations {
		if err := dst.Flush(); err != nil {
			return bytesWritten, fmt.Errorf("flush error for destination %d: %w", i, err)
		}
	}

	return bytesWritten, nil
}
```

</details>


---

#### <a name="newbatchstream"></a>NewBatchStream

NewBatchStream creates a new batch-oriented stream that reads items from
`inner` in chunks of `batchSize` and returns them as []T.


<details><summary>Code</summary>

```go
func NewBatchStream[T any](inner ReadStream[T], batchSize int) ReadStream[[]T] {
	return &BatchStream[T]{
		inner:     inner,
		batchSize: batchSize,
	}
}
```

</details>


---

#### <a name="newcompactstream"></a>NewCompactStream

NewCompactStream creates a new stream that groups consecutive items with the same key.
The keyFunc is used to extract the grouping key from each item (like Python's itemgetter).


<details><summary>Code</summary>

```go
func NewCompactStream[T any, K comparable](inner ReadStream[T], keyFunc func(T) K) ReadStream[[]T] {
	return &CompactStream[T, K]{
		inner:   inner,
		keyFunc: keyFunc,
	}
}
```

</details>


---

#### <a name="newcompactstreamfactory"></a>NewCompactStreamFactory

NewCompactStreamFactory creates a factory for CompactStream instances.


<details><summary>Code</summary>

```go
func NewCompactStreamFactory[T any, K comparable](
	innerFactory ReadStreamFactory[T],
	keyFunc func(T) K,
) ReadStreamFactory[[]T] {
	return func(rc io.ReadCloser) ReadStream[[]T] {
		return NewCompactStream(innerFactory(rc), keyFunc)
	}
}
```

</details>


---

#### <a name="newfilterstream"></a>NewFilterStream

NewFilterStream creates a new ReadStream that filters elements from the inner stream
using the provided predicate function. Only elements that satisfy the predicate
(return true) will be included in the resulting stream.

This is useful for creating data processing pipelines where you need to exclude
certain elements based on custom criteria.


<details><summary>Code</summary>

```go
func NewFilterStream[T any](inner ReadStream[T], predicate func(T) bool) ReadStream[T] {
	return &FilterStream[T]{
		inner:     inner,
		predicate: predicate,
	}
}
```

</details>


---

#### <a name="newflattenerstream"></a>NewFlattenerStream

NewFlattenerStream creates a new ReadStream that flattens slices from the inner stream.
It takes a ReadStream[[]T] and converts it to ReadStream[T] by emitting each
element from the inner slices individually.

This is useful when you have a stream of slices and want to process each
individual element, such as flattening batched data or expanding grouped results.


<details><summary>Code</summary>

```go
func NewFlattenerStream[T any](inner ReadStream[[]T]) ReadStream[T] {
	return &FlattenerStream[T]{
		inner: inner,
	}
}
```

</details>


---

#### <a name="newlinereaderstream"></a>NewLineReaderStream

NewLineReaderStream creates a new ReadStream that reads lines as strings from an io.Reader


<details><summary>Code</summary>

```go
func NewLineReaderStream(reader io.Reader) ReadStream[string] {
	return &LineReaderStream{
		original: reader,
		reader:   bufio.NewReader(reader),
	}
}
```

</details>


---

#### <a name="newmapperstream"></a>NewMapperStream

NewMapperStream creates a new ReadStream that transforms elements from the inner stream
using the provided mapper function. Each element of type T is converted to type V
using the mapper function.

This is useful for creating data processing pipelines where you need to transform
data from one type to another, such as converting strings to uppercase or
extracting specific fields from structs.


<details><summary>Code</summary>

```go
func NewMapperStream[T, V any](inner ReadStream[T], mapper func(T) V) ReadStream[V] {
	return &MapperStream[T, V]{
		inner:  inner,
		mapper: mapper,
	}
}
```

</details>


---

#### <a name="newmemory"></a>NewMemory

NewMemory creates a new ReadStream that reads from a slice in memory.
This is useful for testing, converting slices to streams, or creating
simple data sources for streaming pipelines.

The error parameter allows you to simulate error conditions during streaming.
If err is not nil, the stream will return this error when Err() is called.


<details><summary>Code</summary>

```go
func NewMemory[T any](items []T, err error) *MemoryStream[T] {
	return &MemoryStream[T]{
		items:  items,
		cursor: -1,
		error:  err,
	}
}
```

</details>


---

#### <a name="newmemorywritestream"></a>NewMemoryWriteStream

NewMemoryWriteStream creates a new memory-based WriteStream


<details><summary>Code</summary>

```go
func NewMemoryWriteStream[T any]() *MemoryWriteStream[T] {
	return &MemoryWriteStream[T]{
		items: make([]T, 0),
	}
}
```

</details>


---

#### <a name="newreaderstream"></a>NewReaderStream

NewReaderStream creates a new ReadStream that reads from an io.Reader


<details><summary>Code</summary>

```go
func NewReaderStream(reader io.Reader) ReadStream[[]byte] {
	return &ReaderStream{
		original: reader,
		reader:   bufio.NewReader(reader),
	}
}
```

</details>


---

#### <a name="newwriterstream"></a>NewWriterStream

NewWriterStream creates a new WriteStream that writes to an io.Writer


<details><summary>Code</summary>

```go
func NewWriterStream(writer io.Writer) *WriterStream {
	return &WriterStream{
		writer: writer,
	}
}
```

</details>


---

#### <a name="pipe"></a>Pipe

Pipe copies all items from a ReadStream to a WriteStream
Returns the total number of bytes written and any error


<details><summary>Code</summary>

```go
func Pipe[T any](src ReadStream[T], dst WriteStream[T]) (int64, error) {
	var totalBytes int64

	for src.Next() {
		if err := src.Err(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return totalBytes, fmt.Errorf("read error: %w", err)
		}

		n, err := dst.Write(src.Data())
		if err != nil {
			return totalBytes, fmt.Errorf("write error: %w", err)
		}
		totalBytes += n
	}

	if err := dst.Flush(); err != nil {
		return totalBytes, fmt.Errorf("flush error: %w", err)
	}

	return totalBytes, nil
}
```

</details>


---

#### <a name="seqkeys"></a>SeqKeys

SeqKeys collects all keys from an iter.Seq2 into another Seq.


<details><summary>Code</summary>

```go
func SeqKeys[K, V any](iter iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range iter {
			if !yield(k) {
				break
			}
		}
	}
}
```

</details>


---

#### <a name="seqvalues"></a>SeqValues

SeqValues collects all values from an iter.Seq2 into another Seq.


<details><summary>Code</summary>

```go
func SeqValues[K, V any](iter iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range iter {
			if !yield(v) {
				break
			}
		}
	}
}
```

</details>


---

#### <a name="writeall"></a>WriteAll

WriteAll writes all items from a slice to a WriteStream
Returns the total number of bytes written and any error


<details><summary>Code</summary>

```go
func WriteAll[T any](stream WriteStream[T], items []T) (int64, error) {
	return WriteSeq(stream, slices.Values(items))
}
```

</details>


---

#### <a name="writeseq"></a>WriteSeq

WriteSeq writes all items from an iter.Seq to a WriteStream
Returns the total number of bytes written and any error


<details><summary>Code</summary>

```go
func WriteSeq[T any](stream WriteStream[T], items iter.Seq[T]) (int64, error) {
	bytesWritten := int64(0)

	for v := range items {
		n, err := stream.Write(v)
		if err != nil {
			return 0, fmt.Errorf("write error: %w", err)
		}
		if n == 0 {
			continue
		}
		bytesWritten += n
	}
	if err := stream.Flush(); err != nil {
		return 0, fmt.Errorf("flush error: %w", err)
	}
	return bytesWritten, nil
}
```

</details>


---

#### <a name="writeseqkeys"></a>WriteSeqKeys

WriteSeqKeys writes all keys from an iter.Seq2 to a WriteStream
Returns the total number of bytes written and any error


<details><summary>Code</summary>

```go
func WriteSeqKeys[K, V any](stream WriteStream[K], items iter.Seq2[K, V]) (int64, error) {
	return WriteSeq(stream, SeqKeys(items))
}
```

</details>


---

#### <a name="writeseqvalues"></a>WriteSeqValues

WriteSeqValues writes all values from an iter.Seq2 to a WriteStream
Returns the total number of bytes written and any error


<details><summary>Code</summary>

```go
func WriteSeqValues[K, V any](stream WriteStream[V], items iter.Seq2[K, V]) (int64, error) {
	return WriteSeq(stream, SeqValues(items))
}
```

</details>


---


[⬆️ Back to Top](#table-of-contents)


<br/>

