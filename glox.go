package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/mbanszel/glox/lox"
)

var interpreter = lox.NewInterpreter()

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()

	// for _, a_token := range tokens {
	// 	fmt.Println(a_token)
	// }
	parser := lox.NewParser(tokens)
	statements := parser.Parse()
	// fmt.Println(lox.NewAstPrinter().Print(expression))
	if lox.HadError {
		return
	}
	interpreter.Interpret(statements)
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file", filename)
		panic("exiting")
	}

	run(string(bytes))
	if lox.HadError {
		os.Exit(65)
	}
	if lox.HadRuntimeError {
		os.Exit(70)
	}

}

func runPrompt() {
	prompt := os.Getenv("GLOX_PROMPT")
	if prompt == "" {
		prompt = ">>> "
	}

	rl, err := readline.New(prompt)
	if err != nil {
		panic("Readline failed.")
	}
	defer rl.Close()

	fmt.Println("Type 'exit' to quit.")

	for {
		line, err := rl.Readline()

		if err != nil {
			break
		}

		if line == "exit" {
			break
		}

		run(line)

	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		// TODO: diceide on the exit code
		os.Exit(1)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
