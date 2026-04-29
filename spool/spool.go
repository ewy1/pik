package spool

import (
	"fmt"
	"os"
)

var Print = fmt.Printf

var Warn = func(format string, values ...any) (any, error) {
	return fmt.Fprintf(os.Stderr, format, values...)
}
