PHONY: run-app

build-app:
	go build -o bin/ ./cmd/app

run-app: build-app
	./bin/app

debug-app: build-app
	DEBUG=TRUE ./bin/app

