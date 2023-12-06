package error

import (
	"fmt"
	"os"
)

var HadError = false

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stdin, "[line %d] Error %s: %s\n", line, where, message)
	HadError = true
}

func Error(line int, message string) {
	Report(line, "", message)
}
