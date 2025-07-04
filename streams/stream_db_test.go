package streams

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRows is a mock implementation of DBRows for testing
type mockRows struct {
	data      [][]any
	index     int
	closed    bool
	scanError error
	nextError bool
}

func (m *mockRows) Next() bool {
	if m.nextError && m.index == 1 {
		return false
	}
	if m.index >= len(m.data) {
		return false
	}
	// Don't increment index here, do it in Scan()
	return true
}

func (m *mockRows) Scan(dest ...any) error {
	if m.scanError != nil {
		m.index++ // Still increment on error to move to next row
		return m.scanError
	}
	if m.index >= len(m.data) {
		return errors.New("no more rows")
	}

	row := m.data[m.index]
	for i, d := range dest {
		if i < len(row) {
			switch v := d.(type) {
			case *int:
				*v = row[i].(int)
			case *string:
				*v = row[i].(string)
			}
		}
	}
	m.index++ // Increment after successful scan
	return nil
}

func (m *mockRows) Close() error {
	m.closed = true
	return nil
}

// Test data structures
type User struct {
	ID   int
	Name string
}

type Product struct {
	ID    int
	Title string
}

func TestStream_Next(t *testing.T) {
	tests := []struct {
		name        string
		rows        *mockRows
		scanFn      func(DBRows, *User) error
		expectNext  []bool
		expectData  []User
		expectError bool
	}{
		{
			name: "successful streaming with multiple rows",
			rows: &mockRows{
				data: [][]any{
					{1, "Alice"},
					{2, "Bob"},
					{3, "Charlie"},
				},
			},
			scanFn: func(rows DBRows, user *User) error {
				return rows.Scan(&user.ID, &user.Name)
			},
			expectNext: []bool{true, true, true, false},
			expectData: []User{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 3, Name: "Charlie"},
			},
			expectError: false,
		},
		{
			name: "empty result set",
			rows: &mockRows{
				data: [][]any{},
			},
			scanFn: func(rows DBRows, user *User) error {
				return rows.Scan(&user.ID, &user.Name)
			},
			expectNext:  []bool{false},
			expectData:  nil, // Changed from []User{} to nil
			expectError: false,
		},
		{
			name: "single row",
			rows: &mockRows{
				data: [][]any{
					{1, "Alice"},
				},
			},
			scanFn: func(rows DBRows, user *User) error {
				return rows.Scan(&user.ID, &user.Name)
			},
			expectNext: []bool{true, false},
			expectData: []User{
				{ID: 1, Name: "Alice"},
			},
			expectError: false,
		},
		{
			name: "scan error",
			rows: &mockRows{
				data: [][]any{
					{1, "Alice"},
				},
				scanError: errors.New("scan failed"),
			},
			scanFn: func(rows DBRows, user *User) error {
				return rows.Scan(&user.ID, &user.Name)
			},
			expectNext: []bool{true, false},
			expectData: []User{
				{ID: 0, Name: ""}, // Zero values due to scan error
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := DB(tt.rows, tt.scanFn)

			var actualData []User
			var actualNext []bool

			// Test Next() behavior
			for i := 0; i < len(tt.expectNext); i++ {
				next := stream.Next()
				actualNext = append(actualNext, next)

				if next {
					actualData = append(actualData, stream.Data())
				}
			}

			assert.Equal(t, tt.expectNext, actualNext, "Next() results should match")

			// Handle nil vs empty slice comparison
			if tt.expectData == nil {
				assert.Nil(t, actualData, "Data() results should be nil")
			} else {
				assert.Equal(t, tt.expectData, actualData, "Data() results should match")
			}

			// Check error state
			if tt.expectError {
				assert.Error(t, stream.Err(), "Expected an error")
			} else {
				assert.NoError(t, stream.Err(), "Expected no error")
			}

			// Verify rows were closed
			assert.True(t, tt.rows.closed, "Rows should be closed after streaming")
		})
	}
}

func TestStream_Data_BeforeNext(t *testing.T) {
	rows := &mockRows{
		data: [][]any{
			{1, "Alice"},
		},
	}

	stream := DB(rows, func(rows DBRows, user *User) error {
		return rows.Scan(&user.ID, &user.Name)
	})

	// Calling Data() before Next() should return zero value
	data := stream.Data()
	assert.Equal(t, User{}, data, "Data() before Next() should return zero value")
}

func TestStream_Err_NoError(t *testing.T) {
	rows := &mockRows{
		data: [][]any{
			{1, "Alice"},
		},
	}

	stream := DB(rows, func(rows DBRows, user *User) error {
		return rows.Scan(&user.ID, &user.Name)
	})

	assert.NoError(t, stream.Err(), "Err() should return nil initially")

	stream.Next()
	assert.NoError(t, stream.Err(), "Err() should return nil after successful scan")
}

func TestStream_Close(t *testing.T) {
	rows := &mockRows{
		data: [][]any{
			{1, "Alice"},
		},
	}

	stream := DB(rows, func(rows DBRows, user *User) error {
		return rows.Scan(&user.ID, &user.Name)
	})

	err := stream.Close()
	assert.NoError(t, err, "Close() should not return an error")
}

func TestStream_DifferentTypes(t *testing.T) {
	tests := []struct {
		name     string
		rows     *mockRows
		scanFn   func(DBRows, *Product) error
		expected []Product
	}{
		{
			name: "Product streaming",
			rows: &mockRows{
				data: [][]any{
					{1, "Laptop"},
					{2, "Mouse"},
				},
			},
			scanFn: func(rows DBRows, product *Product) error {
				return rows.Scan(&product.ID, &product.Title)
			},
			expected: []Product{
				{ID: 1, Title: "Laptop"},
				{ID: 2, Title: "Mouse"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := DB(tt.rows, tt.scanFn)

			var results []Product
			for stream.Next() {
				results = append(results, stream.Data())
			}

			assert.Equal(t, tt.expected, results)
			assert.NoError(t, stream.Err())
		})
	}
}

func TestStream_ReadStreamInterface(t *testing.T) {
	rows := &mockRows{
		data: [][]any{
			{1, "Alice"},
		},
	}

	stream := DB(rows, func(rows DBRows, user *User) error {
		return rows.Scan(&user.ID, &user.Name)
	})

	// Verify it implements ReadStream interface
	var readStream ReadStream[User] = stream
	require.NotNil(t, readStream, "Stream should implement ReadStream interface")

	// Test interface methods
	assert.True(t, readStream.Next())
	assert.Equal(t, User{ID: 1, Name: "Alice"}, readStream.Data())
	assert.NoError(t, readStream.Err())
	assert.NoError(t, readStream.Close())
}

func TestDB(t *testing.T) {
	rows := &mockRows{
		data: [][]any{
			{1, "Alice"},
		},
	}

	scanFn := func(rows DBRows, user *User) error {
		return rows.Scan(&user.ID, &user.Name)
	}

	stream := DB(rows, scanFn)

	require.NotNil(t, stream, "DB should return a non-nil stream")
	assert.NotNil(t, stream.rows, "Stream should have rows")
	assert.NotNil(t, stream.scanFn, "Stream should have scan function")
	assert.Nil(t, stream.current, "Stream should have nil current initially")
	assert.NoError(t, stream.err, "Stream should have no error initially")
}

// ExampleStream demonstrates how to use Stream with database rows.
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
