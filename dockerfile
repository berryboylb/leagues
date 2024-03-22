# Use the official Go image as a base
FROM golang:1.22-alpine AS builder

# Install air for live reloading
RUN go install github.com/cosmtrek/air@latest

# init air
RUN air init

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

EXPOSE 8000

# Set the entry point to air with the configuration file
# ENTRYPOINT ["air", "-c", ".air.toml"]
