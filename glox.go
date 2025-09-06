package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mbanszel/glox/lox"
)

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()

	// for _, a_token := range tokens {
	// 	fmt.Println(a_token)
	// }
	parser := lox.NewParser(tokens)
	expression := parser.Parse()

	if lox.HadError {
		return
	}

	fmt.Println(lox.NewAstPrinter().Print(expression))

	interpreter := lox.NewIterpreter()
	result := interpreter.Interpret(expression)
	fmt.Printf("%v\n", result)
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file", filename)
		panic("exiting")
	}

	run(string(bytes))
	if lox.HadError {
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
		lox.HadError = false
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
