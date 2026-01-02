# Stage 1: Build the Go app
FROM golang:1.24.6-bookworm AS builder

# Set working directory inside container
WORKDIR /usr/src/app

# Install git if needed for go modules
RUN apt-get update && apt-get install -y git

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire project
COPY . .

# Build the Go app (point to your main.go)
RUN go build -v -o /run-app ./cmd/myapp

# Stage 2: Minimal final image
FROM debian:bookworm

WORKDIR /app

# Copy the compiled Go binary
COPY --from=builder /run-app .

# Copy static files and templates
COPY ./api/static ./api/static
COPY ./api/templates ./api/templates

# Expose port
EXPOSE 8080

# Start the app
CMD ["/app/run-app"]
