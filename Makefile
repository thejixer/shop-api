build:
	@go build -o bin/shopapi

run: build
	@./bin/shopapi

test:
	@go test -v ./..