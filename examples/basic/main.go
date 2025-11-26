package main

import (
	"fmt"
	"io/fs"
	"os"
	"runtime/pprof"
	"time"

	"github.com/azuyamat/globber/glob"
)

func main() {
	folder := os.Args[1]
	pattern := os.Args[2]
	fmt.Printf("Scanning folder: %s with pattern: %s\n", folder, pattern)
	fsMatcher := glob.FSMatcher(pattern)
	i := 0
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	start := time.Now()
	fsMatcher.WalkDirFS(folder, func(path string, entry fs.DirEntry) error {
		i++
		fmt.Printf("(%d) Matched: %s\n", i, path)
		return nil
	})
	elapsed := time.Since(start)
	fmt.Printf("Walk completed in %d milliseconds\n", elapsed.Milliseconds())
	fmt.Println("Files found:", i)
}
