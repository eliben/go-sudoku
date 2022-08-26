.PHONY: build serve

build:
	GOOS=js GOARCH=wasm go build -o assets/gosudoku.wasm main.go

serve:
	go run ./server
