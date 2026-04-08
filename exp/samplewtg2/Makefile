.PHONY: run clean

# Find wasm_exec.js (location changed in Go 1.24)
GOROOT := $(shell go env GOROOT)
WASMEXEC := $(shell test -f "$(GOROOT)/misc/wasm/wasm_exec.js" && echo "$(GOROOT)/misc/wasm/wasm_exec.js" || echo "$(GOROOT)/lib/wasm/wasm_exec.js")

main.wasm: wasm/main.go
	GOOS=js GOARCH=wasm go build -o $@ ./wasm/

wasm_exec.js:
	cp "$(WASMEXEC)" .

run: main.wasm wasm_exec.js
	@echo "Serving on http://localhost:8080"
	@(sleep 1 && (xdg-open http://localhost:8080 2>/dev/null || open http://localhost:8080 2>/dev/null || true)) &
	@go run main.go

clean:
	rm -f main.wasm wasm_exec.js
