# Start with a minimal base image containing Go runtime
FROM golang:1.22.0 AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod ./

# Download and install Go dependencies
RUN go mod download

# Copy the entire project source code into the container
COPY . .

# Ensure config.yml is copied to the correct location
COPY config.yml /app/config.yml

# Build the Go application
RUN go build -o main cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

CMD ["./main"]
