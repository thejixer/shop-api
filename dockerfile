# syntax=docker/dockerfile:1

FROM golang:1.21.5-bookworm as base

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the code into the container
COPY . .

# Download the dependencies
RUN go mod download

# Build the Go app
RUN go build -o bin cmd/shop-api/main.go 

# Command to run when starting the container
CMD ["./bin"]