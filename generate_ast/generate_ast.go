package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func defineAst(outputDir string, baseName string, types []string) {
	path := fmt.Sprintf("%s/%s.go", outputDir, baseName)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open output file: %s\n", path))
	}
	defer file.Close()

	fmt.Fprintln(file, "package parser")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "import (")
	fmt.Fprintln(file, "  \"github.com/mbanszel/glox/scanner\"")
	fmt.Fprintln(file, ")")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "type "+baseName+" interface {")
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)

	for _, tokenTokenType := range types {
		className := strings.Split(tokenTokenType, ":")[0]
		fields := strings.Split(tokenTokenType, ":")[1]
		className = strings.TrimSpace(className)
		fields = strings.TrimSpace(fields)
		defineTokenType(file, baseName, className, fields)
	}
}

func defineTokenType(file io.Writer, _ string, className string, fieldList string) {
	fields := strings.Split(fieldList, ",")
	fmt.Fprintln(file, "type "+className+" struct {")

	for _, field := range fields {
		field = strings.TrimSpace(field)
		fmt.Fprintln(file, "  "+field)
	}
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)

	fmt.Fprintf(file, "func New%s(%s) %s {\n", className, fieldList, className)
	fmt.Fprintf(file, "  return %s{\n", className)
	for _, field := range fields {
		field = strings.TrimSpace(field)
		name := strings.Split(field, " ")[0]
		fmt.Fprintf(file, "    %s:%s,\n", name, name)
	}
	fmt.Fprintln(file, "  }")
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]
	fmt.Printf("output directory: %s\n", outputDir)

	defineAst(outputDir, "Expr", []string{
		"Binary	  : left Expr, operator scanner.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator scanner.TokenType, right Expr",
	})
}
