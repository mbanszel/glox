generate_ast: generate_ast/generate_ast.go
	cd generate_ast && go build .

glox: glox.go scanner/*.go errors/*.go
	go build .

