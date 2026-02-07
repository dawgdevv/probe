.PHONY: build-ui build-cli build-all dev-ui clean

# Build the React frontend
build-ui:
	cd web && npm run build
	rm -rf internal/web/static
	mkdir -p internal/web/static
	cp -r web/dist/* internal/web/static/

# Build the Go binary (with embedded UI)
build-cli: build-ui
	go build -o apitestercli

# Build everything
build-all: build-cli

# Run frontend in dev mode (proxies API to Go server)
dev-ui:
	cd web && npm run dev

# Run Go server in dev mode
dev-server:
	go run main.go serve

# Run tests via CLI
test:
	go run main.go run tests.yaml

# Clean build artifacts
clean:
	rm -f apitestercli
	rm -rf web/dist
	rm -rf internal/web/static
