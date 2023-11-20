build:
	@go build -o cmd/bin/shopapi cmd/shop-api/main.go

run: build
	@./cmd/bin/shopapi

test:
	@go test -v ./..