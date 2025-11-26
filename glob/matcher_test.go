package glob

import (
	"regexp"
	"strings"
	"testing"
)

func TestMatches(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{
			name:    "exact match",
			pattern: "file.txt",
			path:    "file.txt",
			want:    true,
		},
		{
			name:    "exact match no match",
			pattern: "file.txt",
			path:    "other.txt",
			want:    false,
		},
		{
			name:    "single wildcard",
			pattern: "*.txt",
			path:    "file.txt",
			want:    true,
		},
		{
			name:    "single wildcard no match",
			pattern: "*.txt",
			path:    "file.md",
			want:    false,
		},
		{
			name:    "double wildcard",
			pattern: "**/file.txt",
			path:    "path/to/file.txt",
			want:    true,
		},
		{
			name:    "double wildcard nested",
			pattern: "**/*.go",
			path:    "src/main.go",
			want:    true,
		},
		{
			name:    "path with wildcard",
			pattern: "src/*.go",
			path:    "src/main.go",
			want:    true,
		},
		{
			name:    "path with wildcard no match",
			pattern: "src/*.go",
			path:    "lib/main.go",
			want:    false,
		},
		{
			name:    "double wildcard deep nesting",
			pattern: "**/*.go",
			path:    "a/b/c/d/e/file.go",
			want:    true,
		},
		{
			name:    "double wildcard at end",
			pattern: "src/**",
			path:    "src/a/b/c/file.txt",
			want:    true,
		},
		{
			name:    "double wildcard middle",
			pattern: "src/**/test.go",
			path:    "src/pkg/utils/test.go",
			want:    true,
		},
		{
			name:    "double wildcard middle no match",
			pattern: "src/**/test.go",
			path:    "lib/pkg/test.go",
			want:    false,
		},
		{
			name:    "multiple wildcards",
			pattern: "*.test.*",
			path:    "component.test.js",
			want:    true,
		},
		{
			name:    "empty path",
			pattern: "*.txt",
			path:    "",
			want:    false,
		},
		{
			name:    "root level double wildcard",
			pattern: "**",
			path:    "any/path/file.txt",
			want:    true,
		},
		{
			name:    "complex path with double wildcard",
			pattern: "src/**/*.test.js",
			path:    "src/components/Button/Button.test.js",
			want:    true,
		},
		{
			name:    "complex path no match",
			pattern: "src/**/*.test.js",
			path:    "src/components/Button/Button.js",
			want:    false,
		},
		{
			name:    "hidden file with wildcard",
			pattern: ".*",
			path:    ".gitignore",
			want:    true,
		},
		{
			name:    "no wildcard deep path match",
			pattern: "src/lib/util.go",
			path:    "src/lib/util.go",
			want:    true,
		},
		{
			name:    "no wildcard deep path no match",
			pattern: "src/lib/util.go",
			path:    "src/lib/helper.go",
			want:    false,
		},
		{
			name:    "wildcard in middle of path",
			pattern: "src/*/main.go",
			path:    "src/cmd/main.go",
			want:    true,
		},
		{
			name:    "wildcard in middle no match subdirs",
			pattern: "src/*/main.go",
			path:    "src/a/b/main.go",
			want:    false,
		},
		{
			name:    "negate exact match",
			pattern: "!file.txt",
			path:    "file.txt",
			want:    false,
		},
		{
			name:    "negate exact match - different file",
			pattern: "!file.txt",
			path:    "other.txt",
			want:    true,
		},
		{
			name:    "negate wildcard",
			pattern: "!*.txt",
			path:    "file.txt",
			want:    false,
		},
		{
			name:    "negate wildcard - different extension",
			pattern: "!*.txt",
			path:    "file.md",
			want:    true,
		},
		{
			name:    "negate double wildcard",
			pattern: "!**/*.go",
			path:    "src/main.go",
			want:    false,
		},
		{
			name:    "negate double wildcard - different extension",
			pattern: "!**/*.go",
			path:    "src/main.js",
			want:    true,
		},
		{
			name:    "negate path pattern",
			pattern: "!src/*.go",
			path:    "src/main.go",
			want:    false,
		},
		{
			name:    "negate path pattern - different directory",
			pattern: "!src/*.go",
			path:    "lib/main.go",
			want:    true,
		},
		{
			name:    "negate with dot",
			pattern: "!.gitignore",
			path:    ".gitignore",
			want:    false,
		},
		{
			name:    "negate with dot - different file",
			pattern: "!.gitignore",
			path:    ".env",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := Matcher(tt.pattern)
			got, err := matcher.Matches(tt.path)
			if err != nil {
				t.Fatalf("Matches() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v, pattern: %s, path: %s", got, tt.want, tt.pattern, tt.path)
			}
		})
	}
}

// globToRegex converts a glob pattern to a regex pattern
func globToRegex(pattern string) *regexp.Regexp {
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.ReplaceAll(pattern, `\*\*`, "DOUBLESTAR")
	pattern = strings.ReplaceAll(pattern, `\*`, "[^/]*")
	pattern = strings.ReplaceAll(pattern, "DOUBLESTAR", ".*")
	return regexp.MustCompile("^" + pattern + "$")
}

// Benchmark test cases
var benchmarkCases = []struct {
	name    string
	pattern string
	path    string
}{
	{
		name:    "SimpleWildcard",
		pattern: "*.txt",
		path:    "file.txt",
	},
	{
		name:    "DoubleWildcard",
		pattern: "**/*.go",
		path:    "src/pkg/utils/main.go",
	},
	{
		name:    "ComplexPattern",
		pattern: "src/**/*.test.js",
		path:    "src/components/Button/Button.test.js",
	},
	{
		name:    "DeepNesting",
		pattern: "**/*.go",
		path:    "a/b/c/d/e/f/g/h/i/j/file.go",
	},
	{
		name:    "MiddleWildcard",
		pattern: "src/**/test.go",
		path:    "src/pkg/utils/subpkg/test.go",
	},
	{
		name:    "ExactMatch",
		pattern: "src/lib/util.go",
		path:    "src/lib/util.go",
	},
	{
		name:    "MultipleWildcards",
		pattern: "*.test.*",
		path:    "component.test.js",
	},
	{
		name:    "LongPath",
		pattern: "**/*.config.js",
		path:    "project/src/app/modules/admin/components/settings/webpack.config.js",
	},
}

func BenchmarkGlobMatcher(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			matcher := Matcher(bc.pattern)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = matcher.Matches(bc.path)
			}
		})
	}
}

func BenchmarkRegexMatcher(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			re := globToRegex(bc.pattern)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = re.MatchString(bc.path)
			}
		})
	}
}

func BenchmarkGlobCompilation(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Matcher(bc.pattern)
			}
		})
	}
}

func BenchmarkRegexCompilation(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = globToRegex(bc.pattern)
			}
		})
	}
}

func BenchmarkGlobMatcherWithCompilation(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matcher := Matcher(bc.pattern)
				_, _ = matcher.Matches(bc.path)
			}
		})
	}
}

func BenchmarkRegexMatcherWithCompilation(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				re := globToRegex(bc.pattern)
				_ = re.MatchString(bc.path)
			}
		})
	}
}
