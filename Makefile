
all:
	make build
	make move
	make run_server

run_server:
	echo "Starting Webserver" 
	go run main.go

build:
	echo "Building wasm file"
	cd wasm && make build_wasm

move:
	mv wasm/main.wasm ./frontend/main.wasm

clean:
	rm -f ./frontend/main.wasm