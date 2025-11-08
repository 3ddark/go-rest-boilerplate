# Start from the official Golang image to build the application
FROM golang:1.25.2 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

RUN go install github.com/go-delve/delve/cmd/dlv@v1.25.2


# Build the application for a specific platform
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/api

# Start a new stage from scratch for a smaller image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/main .

# Copy the .env file
COPY .env.development .env

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
