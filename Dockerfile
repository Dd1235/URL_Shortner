# Use official Golang image as base
FROM golang:1.21-alpine

# Set environment variables
ENV GO111MODULE=on

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the Go app
RUN go build -o url-shortner main.go

# Expose the app port (adjust if your app runs on a different port)
EXPOSE 8080

# Command to run the executable
CMD ["./url-shortner"]
