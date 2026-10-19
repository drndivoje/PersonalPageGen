package input

import "testing"

func TestParsePagePath(t *testing.T) {
	tests := []struct {
		name       string
		rootFolder string
		pageFolder string
		expected   string
	}{
		{
			name:       "nested markdown path trims extension",
			rootFolder: "/tmp/site",
			pageFolder: "/tmp/site/blog/posts/first-post.md",
			expected:   "blog/posts/first-post",
		},
		{
			name:       "file without markdown extension stays unchanged",
			rootFolder: "/tmp/site",
			pageFolder: "/tmp/site/blog/posts/first-post.txt",
			expected:   "blog/posts/first-post.txt",
		},
		{
			name:       "page equal to root returns dot",
			rootFolder: "/tmp/site",
			pageFolder: "/tmp/site",
			expected:   ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePagePath(tt.rootFolder, tt.pageFolder)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
