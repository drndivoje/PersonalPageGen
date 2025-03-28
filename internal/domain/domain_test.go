package domain

import (
	"testing"
)

func TestParseHeaderAttribute(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedKey   string
		expectedValue []string
		expectError   bool
	}{
		{
			name:          "Valid single value",
			input:         "author = John Doe",
			expectedKey:   "author",
			expectedValue: []string{"John Doe"},
			expectError:   false,
		},
		{
			name:          "Valid list of values",
			input:         "tags = [Go, Programming, Tutorial]",
			expectedKey:   "tags",
			expectedValue: []string{"Go", "Programming", "Tutorial"},
			expectError:   false,
		},
		{
			name:        "Invalid format",
			input:       "invalid_line",
			expectError: true,
		},
		{
			name:        "Missing value",
			input:       "key = ",
			expectError: true,
		},
		{
			name:          "Valid single value with spaces",
			input:         "title = \"My Blog Post\"",
			expectedKey:   "title",
			expectedValue: []string{"\"My Blog Post\""},
			expectError:   false,
		},
		{
			name:        "Invalid list format",
			input:       "tags = [Go, Programming, Tutorial",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHeaderAttribute(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.name != tt.expectedKey {
				t.Errorf("expected key %s, got %s", tt.expectedKey, result.name)
			}

			if len(result.values) != len(tt.expectedValue) {
				t.Errorf("expected values %v, got %v", tt.expectedValue, result.values)
				return
			}

			for i, v := range result.values {
				if v != tt.expectedValue[i] {
					t.Errorf("expected value %s at index %d, got %s", tt.expectedValue[i], i, v)
				}
			}
		})
	}
}
