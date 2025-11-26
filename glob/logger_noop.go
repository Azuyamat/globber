//go:build !logger

package glob

//nolint:unused,deadcode
func log(format string, a ...interface{}) {
	// No-op
}
