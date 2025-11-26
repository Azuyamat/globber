# Globber

A high-performance Go library for glob pattern matching and file system traversal.

## Overview

Globber provides an optimized implementation for matching file paths against glob patterns (like `**/*.go`, `*.txt`, `src/**/test.js`). Built from the ground up for performance, it uses a custom lexer/parser pipeline and stack-based matching algorithm that significantly outperforms regex-based approaches.

## Features

### Pattern Matching

- **Single wildcard (`*`)** - Matches any characters except `/` within a single path segment
- **Double wildcard (`**`)** - Matches zero or more complete path segments
- **Literal matching** - Exact string matches for file names and paths
- **Dot files** - Patterns like `.*` for hidden files
- **Path separators** - Explicit `/` matching for directory boundaries
- **Negation patterns** - `!` prefix to invert match logic
- **Complex combinations** - Support for patterns like `src/**/*.test.js`

### Performance Optimizations

- **Fast-path detection** - Optimized code paths for exact matches and simple suffix patterns
- **Stack-based matching** - Iterative algorithm eliminates recursion overhead
- **Object pooling** - `sync.Pool` usage reduces memory allocations
- **Zero external dependencies** - Built entirely on Go standard library

### File System Support

- Works with Go's `fs.FS` interface (in-memory, embedded, custom file systems)
- Real file system traversal via `filepath.WalkDir`
- Cross-platform path handling with forward slash normalization

## Installation

```bash
go get github.com/azuyamat/globber
```

## Usage

### Basic Pattern Matching

```go
package main

import (
    "fmt"
    "github.com/azuyamat/globber/glob"
)

func main() {
    matcher := glob.Matcher("**/*.go")

    if matcher.Match("src/main.go") {
        fmt.Println("Match found!")
    }
}
```

### File System Walking

```go
package main

import (
    "fmt"
    "io/fs"
    "github.com/azuyamat/globber/glob"
)

func main() {
    fsMatcher := glob.FSMatcher("**/*.go")

    err := fsMatcher.WalkDirFS("/path/to/scan", func(path string, entry fs.DirEntry) error {
        fmt.Printf("Matched: %s\n", path)
        return nil
    })

    if err != nil {
        panic(err)
    }
}
```

### Using fs.FS Interface

```go
package main

import (
    "fmt"
    "io/fs"
    "os"
    "github.com/azuyamat/globber/glob"
)

func main() {
    fsys := os.DirFS("/path/to/scan")
    fsMatcher := glob.FSMatcher("src/**/*.test.js")

    err := fsMatcher.Walk(fsys, func(path string, entry fs.DirEntry) error {
        fmt.Printf("Test file: %s\n", path)
        return nil
    })

    if err != nil {
        panic(err)
    }
}
```

## Pattern Syntax

| Pattern | Description | Example | Matches |
|---------|-------------|---------|---------|
| `*` | Single segment wildcard | `*.go` | `main.go`, `test.go` |
| `**` | Multi-segment wildcard | `**/*.go` | `src/main.go`, `a/b/c/test.go` |
| `/` | Path separator | `src/*.go` | `src/main.go` (not `src/a/b.go`) |
| `.` | Literal dot | `.*` | `.gitignore`, `.env` |
| `!` | Negation | `!*.test.js` | Inverts the match |

### Example Patterns

- `*.txt` - All text files in root
- `**/*.go` - All Go files in any directory
- `src/**/*.test.js` - All test files under src/
- `.*` - All hidden files
- `*.test.*` - Files with `.test.` in the name

## Architecture

### Compilation Pipeline

1. **Scanner** ([scanner.go](glob/scanner.go)) - Character-by-character tokenization
2. **Lexer** ([lexer.go](glob/lexer.go)) - Token stream generation
3. **Parser** ([parser.go](glob/parser.go)) - Abstract Syntax Tree construction
4. **Compiler** ([compiler.go](glob/compiler.go)) - Orchestrates the pipeline

### Matching Engine

- **AST Nodes** ([node.go](glob/node.go)) - Linked list structure with stack-based matching
- **Matcher** ([matcher.go](glob/matcher.go)) - Public API with fast-path optimizations
- **Stack Pool** - Reusable stacks via `sync.Pool` for reduced allocations

### File System Layer

The library provides both serial and parallel walking implementations:

- **Serial Walking** ([fs.go:116](glob/fs.go#L116)) - Single-threaded directory traversal
- **Parallel Walking v1** ([fs.go:29](glob/fs.go#L29)) - Worker pool with semaphore limiting (CPU cores × 4)
- **Parallel Walking v2** ([fs.go:151](glob/fs.go#L151)) - Alternative parallel implementation (CPU cores × 2)
- **FS Abstraction** - Works with any `fs.FS` implementation

Currently, the `Walk()` method uses serial walking. The parallel implementations are available in the codebase for future optimization.

## Performance

Globber is designed for high-throughput scenarios:

- Significantly faster than regex-based matching
- Minimal memory allocations through object pooling
- Inline function hints for critical hot paths

Run benchmarks:

```bash
cd glob
go test -bench=. -benchmem
```

## Development

### Requirements

- Go 1.23.3 or later

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./glob/...
```

### Running Examples

```bash
cd examples/basic
go run main.go /path/to/scan "**/*.go"
```

Enable CPU profiling:

```bash
cd examples/basic
go run main.go /path/to/scan "**/*.go"
go tool pprof cpu.prof
```

### Debug Logging

Build with the `logger` tag to enable debug output:

```bash
go build -tags logger ./...
```

## Project Structure

```
globber/
├── glob/                   # Core library package
│   ├── compiler.go        # Pattern compilation orchestration
│   ├── lexer.go           # Tokenization
│   ├── parser.go          # AST generation
│   ├── scanner.go         # Character scanning
│   ├── token.go           # Token definitions
│   ├── node.go            # AST node and matching logic
│   ├── matcher.go         # Main matcher interface
│   ├── fs.go              # File system walking
│   ├── logger.go          # Debug logging (build tag)
│   └── *_test.go          # Tests and benchmarks
├── examples/
│   ├── basic/             # CLI example with profiling
│   └── fs.go              # Simple FS implementation
└── go.mod                 # Module definition
```

## License

See the repository for license information.

## Contributing

Contributions are welcome! Please ensure:

- All tests pass: `go test ./...`
- Code is formatted: `go fmt ./...`
- Benchmarks show no regression: `go test -bench=.`

## Use Cases

- Build tools and task runners
- File search utilities
- `.gitignore` style pattern matching
- Asset bundlers and compilers
- Code generation tools
- Development tool watchers
