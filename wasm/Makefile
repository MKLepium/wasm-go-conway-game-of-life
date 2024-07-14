basedir = ./..



all:
	make build_wasm

build_wasm:
	echo "Building wasm file"
	GOOS=js GOARCH=wasm go build -o main.wasm
