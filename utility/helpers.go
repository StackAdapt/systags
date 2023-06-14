package utility

import (
	"fmt"
	"io"
)

// Fprintf performs a Fprintf with error handling
func Fprintf(w io.Writer, format string, a ...any) {

	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		panic(err)
	}
}

// Fprintln performs a Fprintln with error handling
func Fprintln(w io.Writer, a ...any) {

	_, err := fmt.Fprintln(w, a...)
	if err != nil {
		panic(err)
	}
}
