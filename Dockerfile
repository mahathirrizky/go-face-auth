# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder-go

WORKDIR /app

# Install git and other build deps if needed by Go
RUN apk add --no-cache git

# Download dependencies first (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o go-face-auth main.go

# Stage 2: Final runtime image (Python + Go binary)
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies for OpenCV and DeepFace
RUN apt-get update && apt-get install -y \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy python dependencies file
COPY requirements.txt .

# Setup Python virtual environment so main.go can find it at .venv/bin/python3
RUN python -m venv .venv
RUN .venv/bin/pip install --no-cache-dir -r requirements.txt

# Copy the Python scripts
COPY face_recognition_server.py .

# Copy the Go binary and other necessary folders from builder
COPY --from=builder-go /app/go-face-auth .
COPY config ./config

# Expose Go API port
EXPOSE 8080

# Environment variables (Coolify will likely pass these in at runtime)
ENV PORT=8080
ENV GIN_MODE=release

# Start the Go application
CMD ["./go-face-auth"]
