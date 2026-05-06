package cmd

import (
	"fmt"
	"os"
)

var verbose bool

func logf(format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, "> "+format+"\n", args...)
	}
}
