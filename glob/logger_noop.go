//go:build !logger

package glob

func log(format string, a ...interface{}) {
	// No-op
}
