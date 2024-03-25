# Use the official Go image as a base
FROM golang:1.22-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Use a minimal base image for the final build
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root

# Copy the built Go executable from the builder stage
COPY --from=builder /app/app .

# Expose the port the app runs on
EXPOSE 8000

# Command to run the executable
CMD ["./app"]
