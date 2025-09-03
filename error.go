package main

import (
	"fmt"
)

var hadError = false

func emit_error(line int, message string) {
	report_error(line, "", message)
}

func report_error(line int, where string, message string) {
	fmt.Println("[line", line, "] Error", where, ":", message)
	hadError = true
}
