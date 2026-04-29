package spool

import (
	"fmt"
	"os"
)

var (
	Stderr = os.Stderr
	Stdout = os.Stdout
)

var Print = func(format string, values ...any) (any, error) {
	return fmt.Fprintf(Stdout, format, values...)
}

var Warn = func(format string, values ...any) (any, error) {
	return fmt.Fprintf(Stderr, format, values...)
}
