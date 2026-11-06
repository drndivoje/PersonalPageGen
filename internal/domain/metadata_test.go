package domain

import (
	"testing"
)
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
