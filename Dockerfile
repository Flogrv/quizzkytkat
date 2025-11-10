# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o quiz-server .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/quiz-server .

# Create necessary directories
RUN mkdir -p /app/data /app/.ssh

# Generate SSH host key
RUN apk add --no-cache openssh && \
    ssh-keygen -t ed25519 -f /app/.ssh/id_ed25519 -N "" && \
    apk del openssh

# Copy questions file if exists (will be overridden by volume)
COPY questions.json* ./questions.json* 2>/dev/null || echo '{"questions":[]}' > questions.json

# Expose SSH port
EXPOSE 2222

# Run the application
CMD ["./quiz-server"]
