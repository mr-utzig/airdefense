app-build:
	@go build -v -o build/airdefense ./cmd/airdefense/*.go

wasm-build:
	env GOOS=js GOARCH=wasm go build -v -o web/js/airdefense.wasm ./cmd/airdefense/*.go

app-run:
	@go run ./cmd/airdefense/*.go

app-build-run:
	./build/airdefense

app-wasm-run:
	@serve -addr :42069 -dir ./web