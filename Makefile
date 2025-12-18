.PHONY: build build-wasm test clean convert-examples

# Build native binary
build:
	go build -o bin/jww-dxf ./cmd/jww-dxf

# Build WebAssembly
build-wasm:
	GOOS=js GOARCH=wasm go build -o dist/jww-dxf.wasm ./wasm/

# Copy wasm_exec.js from Go installation
copy-wasm-exec:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/

# Build WASM and copy support files
dist: build-wasm copy-wasm-exec

# Run tests
test:
	go test -v ./...

stat:
	go run ./cmd/jww-stats/ examples/jww

# Convert all JWW files in examples/jww to DXF and save to examples/converted
convert-examples: build
	@mkdir -p examples/converted
	@for f in examples/jww/*.jww; do \
		if [ -f "$$f" ]; then \
			echo "Converting $$f..."; \
			./bin/jww-dxf "$$f" -o "examples/converted/$$(basename "$$f" .jww).dxf"; \
		fi \
	done
	@echo "Done. Converted files are in examples/converted/"

# Clean build artifacts
clean:
	rm -rf bin/ dist/
