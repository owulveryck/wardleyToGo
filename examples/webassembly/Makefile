WASMEXEC=$(shell go env GOROOT)/misc/wasm/wasm_exec.js

all: assets/wtg.wasm assets/wasm_exec.js assets/index.html assets/script.js assets/style.css

assets/wtg.wasm: wasm/*.go
	GOOS=js GOARCH=wasm go build -o $@ $^

assets/wasm_exec.js: $(WASMEXEC)
	cp $^ $@

assets/style.css: server/assets/style.css
	cp $^ $@
assets/script.js: server/assets/script.js
	cp $^ $@
assets/index.html: server/assets/index.html
	cp $^ $@

install: all
	cp assets/index.html assets/script.js assets/style.css assets/wtg.wasm ../../docs

clean:
	rm assets/wtg.wasm
