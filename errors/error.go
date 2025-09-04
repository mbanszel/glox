package errors

import (
	"fmt"
)

var HadError = false

func Emit(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("[line %d ] Error %s: %s\n", line, where, message)
	HadError = true
}
