package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/azuyamat/gear/command"
	"github.com/azuyamat/globber/glob"
)

var rootCommand = command.NewRootCommand("globber", "A file globbing utility")
var countCommand = command.NewExecutableCommand("count", "Count files matching the glob pattern").
	Args(
		command.NewStringArg("path", "Path to scan"),
		command.NewStringArg("pattern", "Glob pattern to match"),
	).
	Flags(
		command.NewBoolFlag("verbose", "v", "Enable verbose output", false),
	).
	Handler(func(ctx *command.Context, args command.ValidatedArgs) error {
		path := args.String("path")
		pattern := args.String("pattern")
		verbose := args.FlagBool("verbose")
		if verbose {
			fmt.Println("Verbose mode enabled")
		}
		fmt.Printf("Counting files in folder: %s with pattern: %s\n", path, pattern)

		fsMatcher := glob.FSMatcher(pattern)
		count := 0
		start := time.Now()
		err := fsMatcher.WalkDirFS(path, func(path string, entry fs.DirEntry) error {
			count++
			if verbose {
				fmt.Printf("(%d) Matched: %s\n", count, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
		elapsed := time.Since(start)
		fmt.Printf("Counted %d files in %d milliseconds\n", count, elapsed.Milliseconds())
		return nil
	})

func init() {
	rootCommand.AddChild(countCommand)
}

func main() {
	err := rootCommand.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
