package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParsePoolingParams(t *testing.T) {
	type testCase struct {
		name                    string
		inputURI                string
		expectedCleanURI        string
		expectedMaxOpenConns    int
		expectedMinConns        int
		expectedMaxConnLifetime time.Duration
		expectedMaxConnIdleTime time.Duration
		expectError             bool
	}

	testCases := []testCase{
		{
			name:             "no pool parameters",
			inputURI:         "postgres://user:pass@localhost/dbname",
			expectedCleanURI: "postgres://user:pass@localhost/dbname",
			// all numeric/durations default to 0
			expectedMaxOpenConns:    0,
			expectedMinConns:        0,
			expectedMaxConnLifetime: 0,
			expectedMaxConnIdleTime: 0,
			expectError:             false,
		},
		{
			name: "all valid parameters",
			inputURI: "postgres://user:pass@localhost/dbname?" +
				"pool_max_conns=10&pool_min_conns=2&pool_max_conn_lifetime=5m&pool_max_conn_idle_time=1m",
			// After cleaning, no query parameters remain.
			expectedCleanURI:        "postgres://user:pass@localhost/dbname",
			expectedMaxOpenConns:    10,
			expectedMinConns:        2,
			expectedMaxConnLifetime: 5 * time.Minute,
			expectedMaxConnIdleTime: 1 * time.Minute,
			expectError:             false,
		},
		{
			name: "unsupported parameters are removed",
			inputURI: "postgres://user:pass@localhost/dbname?" +
				"pool_max_conns=15&pool_health_check_period=30s&pool_max_conn_lifetime_jitter=10s",
			expectedCleanURI:        "postgres://user:pass@localhost/dbname",
			expectedMaxOpenConns:    15,
			expectedMinConns:        0,
			expectedMaxConnLifetime: 0,
			expectedMaxConnIdleTime: 0,
			expectError:             false,
		},
		{
			name:     "invalid pool_max_conns value",
			inputURI: "postgres://user:pass@localhost/dbname?pool_max_conns=abc",
			// Error expected, so the clean URI and values are irrelevant.
			expectError: true,
		},
		{
			name:        "invalid pool_max_conn_lifetime value",
			inputURI:    "postgres://user:pass@localhost/dbname?pool_max_conn_lifetime=notaduration",
			expectError: true,
		},
		{
			name:        "invalid pool_min_conns value",
			inputURI:    "postgres://user:pass@localhost/dbname?pool_min_conns=xyz",
			expectError: true,
		},
		{
			name:        "invalid pool_max_conn_idle_time value",
			inputURI:    "postgres://user:pass@localhost/dbname?pool_max_conn_idle_time=5minutes", // should be like "5m"
			expectError: true,
		},
		{
			name: "mix of valid and extraneous parameters",
			inputURI: "postgres://user:pass@localhost/dbname?" +
				"pool_max_conns=20&pool_min_conns=4&dummy_param=foo&pool_max_conn_idle_time=30s",
			// Expected: dummy_param remains in the URI if not explicitly removed; adjust your expectations as needed.
			expectedCleanURI:        "postgres://user:pass@localhost/dbname?dummy_param=foo",
			expectedMaxOpenConns:    20,
			expectedMinConns:        4,
			expectedMaxConnLifetime: 0,
			expectedMaxConnIdleTime: 30 * time.Second,
			expectError:             false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params, err := parsePoolingParams(tc.inputURI)
			if tc.expectError {
				assert.Error(t, err, "expected error but got nil")
				return
			}
			assert.NoError(t, err, "unexpected error")

			assert.Equal(t, tc.expectedCleanURI, params.CleanURI, "clean URI mismatch")
			assert.Equal(t, tc.expectedMaxOpenConns, params.MaxOpenConns, "MaxOpenConns mismatch")
			assert.Equal(t, tc.expectedMinConns, params.MinConns, "MinConns mismatch")
			assert.Equal(
				t,
				tc.expectedMaxConnLifetime,
				params.MaxConnLifetime,
				"MaxConnLifetime mismatch",
			)
			assert.Equal(
				t,
				tc.expectedMaxConnIdleTime,
				params.MaxConnIdleTime,
				"MaxConnIdleTime mismatch",
			)
		})
	}
}
