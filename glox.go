package main

import (
	"bufio"
	"fmt"
	"os"
)

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, a_token := range tokens {
		fmt.Println(a_token)
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file", filename)
		panic("exiting")
	}

	run(string(bytes))
	if hadError {
		panic("exiting")
	}

}

func runPrompt() {
	fmt.Println("This is glox.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		run(line)
		hadError = false
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		panic("exiting")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
