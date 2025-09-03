package main

import (
  "fmt"
)

var hadError = false

func error(line int, message string) {
  report(line, "", message)
}

func report(line int, where string, message string) {
  fmt.Println("[line", line, "] Error", where, ":", message)
  hadError = true
}
