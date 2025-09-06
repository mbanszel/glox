package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
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
	result, err := interpreter.Interpret(expression)
	if err != nil {
		rtError := err.(lox.RuntimeError)
		token := rtError.GetToken()
		lox.Report(
			token.Line,
			token.Lexeme,
			fmt.Sprintf(
				"Runtime error: %s",
				rtError.GetMessage(),
			),
		)
		lox.HadError = false
	}
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
		panic("exiting")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
