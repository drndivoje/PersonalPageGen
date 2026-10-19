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

func TestParseMetadataRelativePathLanguage(t *testing.T) {
	content := []byte("+++\ntitle = \"Post\"\n+++")

	tests := []struct {
		name                 string
		relativePath         string
		expectedRelativePath string
		expectedLanguage     string
		expectError          bool
	}{
		{
			name:                 "no language suffix keeps path unchanged",
			relativePath:         "blog/firstPost",
			expectedRelativePath: "blog/firstPost",
			expectedLanguage:     "en",
		},
		{
			name:                 "valid language suffix updates path and language",
			relativePath:         "blog/firstPost.de",
			expectedRelativePath: "blog/firstPost",
			expectedLanguage:     "de",
		},
		{
			name:                 "valid language suffix with leading slash updates path and language",
			relativePath:         "/blog/firstPost.de",
			expectedRelativePath: "/blog/firstPost",
			expectedLanguage:     "de",
		},
		{
			name:         "invalid two letter language suffix returns error",
			relativePath: "blog/firstPost.zz",
			expectError:  true,
		},
		{
			name:                 "long extension is not treated as language code",
			relativePath:         "blog/firstPost.markdown",
			expectedRelativePath: "blog/firstPost.markdown",
			expectedLanguage:     "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := parseMetadata(content, tt.relativePath)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if metadata.RelativePath != tt.expectedRelativePath {
				t.Fatalf("expected relative path %q, got %q", tt.expectedRelativePath, metadata.RelativePath)
			}

			if metadata.Language != tt.expectedLanguage {
				t.Fatalf("expected language %q, got %q", tt.expectedLanguage, metadata.Language)
			}
		})
	}
}
