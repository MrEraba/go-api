.PHONY: test test-coverage

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out
