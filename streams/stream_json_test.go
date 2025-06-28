package streams

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJSONEachRowStream(t *testing.T) {
	buf := io.NopCloser(strings.NewReader(strings.TrimSpace(`
{"p":1.1,"v":123.345,"t":"1682298820000"}
{"p":1.2,"v":123.346,"t":"1682298821000"}
{"p":1.3,"v":123.347,"t":"1682298890000"}
{"p":1,  "v":123.349,"t":"1682298891000"}
`)))

	type trade struct {
		Price  float64 `json:"p"`
		Volume float64 `json:"v"`
		Time   string  `json:"t"`
	}

	s := JSON[trade](buf)

	expected := []trade{
		{Price: 1.1, Volume: 123.345, Time: "1682298820000"},
		{Price: 1.2, Volume: 123.346, Time: "1682298821000"},
		{Price: 1.3, Volume: 123.347, Time: "1682298890000"},
		{Price: 1, Volume: 123.349, Time: "1682298891000"},
	}
	i := 0

	for s.Next() {
		if err := s.Err(); err != nil {
			t.Fatalf("unexpected err %v", err)
		}

		actual := s.Data()

		assert.Equal(t, actual, expected[i])
		i++
	}
}

// ExampleJSON demonstrates reading JSON data line by line.
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
