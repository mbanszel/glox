generate_ast: generate_ast/generate_ast.go
	cd generate_ast && go build .

lox/err.go: generate_ast.go
	cd generate_ast && ./generate_ast ../lox

glox: glox.go lox/*.go 
	go build .

