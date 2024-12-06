# Use the official Golang image for building
FROM golang:1.22.3 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image to run the compiled binary
FROM alpine:3.18

WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/main .

# Expose the port that the service will run on
EXPOSE 8080

# Run the binary
CMD ["./main"]