package glob

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestMatcherWithFS(t *testing.T) {
	testFS := fstest.MapFS{
		"file.txt": &fstest.MapFile{
			Data: []byte("content"),
		},
		"dir/nested.go": &fstest.MapFile{
			Data: []byte("package main"),
		},
		"dir/another.txt": &fstest.MapFile{
			Data: []byte("text"),
		},
		"other/test.go": &fstest.MapFile{
			Data: []byte("package test"),
		},
	}

	tests := []struct {
		pattern string
		want    []string
	}{
		{
			pattern: "*.txt",
			want:    []string{"file.txt"},
		},
		{
			pattern: "**/*.go",
			want:    []string{"dir/nested.go", "other/test.go"},
		},
		{
			pattern: "dir/*",
			want:    []string{"dir/another.txt", "dir/nested.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			m := FSMatcher(tt.pattern)

			var got []string
			err := m.Walk(testFS, func(path string, entry fs.DirEntry) error {
				if !entry.IsDir() {
					got = append(got, path)
				}
				return nil
			})

			if err != nil {
				t.Fatalf("Walk() error = %v", err)
			}

			if len(got) != len(tt.want) {
				t.Errorf("got %d matches, want %d: %v", len(got), len(tt.want), got)
				return
			}

			wantMap := make(map[string]bool)
			for _, w := range tt.want {
				wantMap[w] = true
			}

			for _, g := range got {
				if !wantMap[g] {
					t.Errorf("unexpected match: %s", g)
				}
			}
		})
	}
}
