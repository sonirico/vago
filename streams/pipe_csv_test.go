package streams

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of csvMarshaler for testing
type mockCsvMarshaler struct {
	ID    int
	Name  string
	Email string
	Err   error
}

func (m mockCsvMarshaler) MarshalCSV() ([]string, []string, error) {
	if m.Err != nil {
		return nil, nil, m.Err
	}
	header := []string{"ID", "Name", "Email"}
	record := []string{strconv.FormatInt(int64(m.ID), 10), m.Name, m.Email}
	return header, record, nil
}

func TestPipeCSVTransform_WriteTo(t *testing.T) {
	tests := []struct {
		name          string
		separator     rune
		stream        ReadStream[mockCsvMarshaler]
		expectedCSV   string
		expectedError error
	}{
		{
			name:      "valid stream with multiple records comma separated",
			separator: CSVSeparatorComma,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedCSV:   "ID,Name,Email\n1,Alice,alice@example.com\n2,Bob,bob@example.com\n",
			expectedError: nil,
		},
		{
			name:      "valid stream with multiple records tab separated",
			separator: CSVSeparatorTab,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedCSV:   "ID\tName\tEmail\n1\tAlice\talice@example.com\n2\tBob\tbob@example.com\n",
			expectedError: nil,
		},
		{
			name:          "empty stream",
			separator:     CSVSeparatorComma,
			stream:        MemReader([]mockCsvMarshaler{}, nil),
			expectedCSV:   "",
			expectedError: nil,
		},
		{
			name:      "stream with single record",
			separator: CSVSeparatorComma,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
			}, nil),
			expectedCSV:   "ID,Name,Email\n1,Alice,alice@example.com\n",
			expectedError: nil,
		},
		{
			name:          "stream with error",
			separator:     CSVSeparatorComma,
			stream:        MemReader([]mockCsvMarshaler{}, errors.New("stream error")),
			expectedCSV:   "",
			expectedError: errors.New("stream error"),
		},
		{
			name:          "stream with error io.EOF",
			separator:     CSVSeparatorComma,
			stream:        MemReader([]mockCsvMarshaler{}, io.EOF),
			expectedCSV:   "",
			expectedError: nil,
		},
		{
			name:      "marshaler with error",
			separator: CSVSeparatorComma,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{
					ID:    2,
					Name:  "Error",
					Email: "error@example.com",
					Err:   errors.New("marshal error"),
				},
			}, nil),
			expectedCSV:   "",
			expectedError: errors.New("csv marshaling error: marshal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			transform := CSVTransform(tt.stream, tt.separator)

			written, err := transform.WriteTo(buffer)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(len(strings.Split(tt.expectedCSV, "\n")))-1, written)
				assert.Equal(t, tt.expectedCSV, buffer.String())
			}
		})
	}
}

func TestPipeCSVTransform_WriteToFile(t *testing.T) {
	stream := MemReader([]mockCsvMarshaler{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}, nil)

	_, err := CSVTransform(stream, CSVSeparatorTab).WriteTo(os.Stdout)
	assert.NoError(t, err)
}

// Employee struct for CSV examples
type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

// MarshalCSV implements the csvMarshaler interface for Employee
func (e Employee) MarshalCSV() ([]string, []string, error) {
	header := []string{"ID", "Name", "Department", "Salary"}
	record := []string{
		strconv.Itoa(e.ID),
		e.Name,
		e.Department,
		fmt.Sprintf("%.2f", e.Salary),
	}
	return header, record, nil
}

// ExampleCSVTransform demonstrates converting a stream to CSV format.
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

// ProductCSV wraps Product to implement csvMarshaler interface
type ProductCSV struct {
	SKU   string
	Name  string
	Price float64
}

func (p ProductCSV) MarshalCSV() ([]string, []string, error) {
	header := []string{"SKU", "Product Name", "Price"}
	record := []string{p.SKU, p.Name, fmt.Sprintf("$%.2f", p.Price)}
	return header, record, nil
}

// ExampleCSVTransform_tabSeparated demonstrates CSV with tab separator.
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

// ExamplePipeCSV demonstrates using the PipeCSV convenience function.
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
