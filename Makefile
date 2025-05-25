all: build test

build:
	@echo "Building..."
	@go build -o bin/bot .

test:
	@echo "Testing..."
	@go test -v ./...

run:
	@echo "Running..."
	@go run .

clean:
	@echo "Cleaning..."
	@rm bin/bot

.PHONY: all build test run clean
