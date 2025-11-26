package glob

import (
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"
)

func createLargeFS(dirs, filesPerDir int) fstest.MapFS {
	testFS := make(fstest.MapFS)

	for d := 0; d < dirs; d++ {
		dirName := fmt.Sprintf("dir%d", d)
		for f := 0; f < filesPerDir; f++ {
			path := fmt.Sprintf("%s/file%d.txt", dirName, f)
			testFS[path] = &fstest.MapFile{Data: []byte("content")}
		}

		for f := 0; f < filesPerDir/2; f++ {
			path := fmt.Sprintf("%s/code%d.go", dirName, f)
			testFS[path] = &fstest.MapFile{Data: []byte("package main")}
		}

		subDir := fmt.Sprintf("%s/sub", dirName)
		for f := 0; f < filesPerDir/4; f++ {
			path := fmt.Sprintf("%s/nested%d.txt", subDir, f)
			testFS[path] = &fstest.MapFile{Data: []byte("nested")}
		}
	}

	return testFS
}

func BenchmarkWalk(b *testing.B) {
	sizes := []struct {
		name        string
		dirs        int
		filesPerDir int
	}{
		{"Small", 10, 10},
		{"Medium", 50, 20},
		{"Large", 100, 50},
		{"XLarge", 200, 100},
	}

	for _, size := range sizes {
		testFS := createLargeFS(size.dirs, size.filesPerDir)

		b.Run(size.name, func(b *testing.B) {
			m := FSMatcher("**/*.go")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				count := 0
				_ = m.Walk(testFS, func(path string, entry fs.DirEntry) error {
					count++
					return nil
				})
			}
		})
	}
}
