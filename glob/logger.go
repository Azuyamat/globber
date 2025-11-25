//go:build logger

package glob

import "fmt"

func log(format string, a ...interface{}) {
	fmt.Printf("[GLOB] "+format+"\n", a...)
}
