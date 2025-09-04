all: glox.go scanner/*.go errors/*.go
	CGO_ENABLED=0 GOFLAGS="-ldflags=-s -w" go build -trimpath -o glox .
