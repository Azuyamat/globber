package glob

import (
	"fmt"
	"strings"
	"sync"
)

type node struct {
	Type   tokenType
	Value  string
	Next   *node
	Negate bool

	Children []*node
}

func (n *node) String() string {
	return fmt.Sprintf("Node{Type: %d, Value: %s}\n", n.Type, n.Value)
}

func (n *node) Tree() string {
	sb := ""
	sb += fmt.Sprintf("- %s\n", n.String())
	for _, child := range n.Children {
		sb += "  " + child.Tree()
	}
	if n.Next != nil {
		sb += n.Next.Tree()
	}
	return sb
}

type matchResult struct {
	node *node
	pos  int
}

var stackPool = sync.Pool{
	New: func() interface{} {
		s := make([]matchResult, 0, 32)
		return &s
	},
}

func (n *node) MatchString(path string, pos int) (bool, error) {
	if n.Negate {
		match, err := n.matchWithoutNegate(path, pos)
		return !match, err
	}
	return n.matchWithoutNegate(path, pos)
}

func (n *node) matchWithoutNegate(path string, pos int) (bool, error) {
	stackPtr := stackPool.Get().(*[]matchResult)
	stack := (*stackPtr)[:0]
	stack = append(stack, matchResult{node: n, pos: pos})
	defer func() {
		*stackPtr = stack
		stackPool.Put(stackPtr)
	}()

	for len(stack) > 0 {
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		current := state.node
		i := state.pos

	nodeLoop:
		for current != nil {
			if i > len(path) {
				break nodeLoop
			}

			switch current.Type {
			case tokenLiteral:
				if !matchLiteral(path, i, current.Value) {
					break nodeLoop
				}
				i += len(current.Value)

			case tokenDot:
				if i >= len(path) || path[i] != '.' {
					break nodeLoop
				}
				i++

			case tokenSlash:
				if i >= len(path) || path[i] != '/' {
					break nodeLoop
				}
				i++

			case tokenStar:
				if current.Next == nil {
					hasSlashInRemainder := false
					for j := i; j < len(path); j++ {
						if path[j] == '/' {
							hasSlashInRemainder = true
							break
						}
					}
					if !hasSlashInRemainder {
						return true, nil
					}
					break nodeLoop
				}

				hasSlash, endOfSegment := findSlash(path, i)
				if !hasSlash {
					endOfSegment = len(path)
				}

				for j := endOfSegment; j >= i; j-- {
					stack = append(stack, matchResult{node: current.Next, pos: j})
				}
				break nodeLoop

			case tokenDoubleStar:
				if current.Next == nil {
					return true, nil
				}

				if match, _ := current.Next.MatchString(path, i); match {
					return true, nil
				}

				for j := len(path); j > i; j-- {
					stack = append(stack, matchResult{node: current.Next, pos: j})
				}
				break nodeLoop

			default:
				break nodeLoop
			}

			current = current.Next
		}

		if current == nil && i == len(path) {
			return true, nil
		}
	}

	return false, nil
}

//go:inline
func matchLiteral(path string, pos int, literal string) bool {
	if pos+len(literal) > len(path) {
		return false
	}
	if len(literal) == 1 {
		return path[pos] == literal[0]
	}
	return path[pos:pos+len(literal)] == literal
}

func findSlash(path string, pos int) (bool, int) {
	if pos >= len(path) {
		return false, len(path)
	}
	idx := strings.IndexByte(path[pos:], '/')
	if idx == -1 {
		return false, len(path)
	}
	return true, pos + idx
}
