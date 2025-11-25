package main

import (
	"fmt"

	"github.com/azuyamat/globber/glob"
)

func main() {
	matcher := glob.Matcher("*.txt")
	ast := matcher.AST()
	fmt.Println(ast.Tree())
	matches, err := matcher.Matches("file.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Matches:", matches)
}
