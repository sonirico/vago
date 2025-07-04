package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParseTsFromMongoObjectID tests the MongoParseTsFromObjectID function.
func TestParseTsFromMongoObjectID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTime  time.Time
		expectedError bool
	}{
		{
			name: "Valid ObjectId",
			id:   "64ebcbf2e47b5d95fb000000", // Replace with a known valid ObjectId
			expectedTime: time.Unix(1693174770, 0).
				UTC(),
			// Expected time from the ObjectId's timestamp
			expectedError: false,
		},
		{
			name:          "Invalid ObjectId format",
			id:            "invalid_object_id",
			expectedTime:  time.Time{},
			expectedError: true,
		},
		{
			name:          "Empty ObjectId",
			id:            "",
			expectedTime:  time.Time{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTime, err := MongoParseTsFromObjectID(tt.id)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tt.expectedTime, actualTime, "Expected time does not match actual time")
			}
		})
	}
}
