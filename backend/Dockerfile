# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build main server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Build migration tool binary
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

# Build seed binary
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed/main.go

# Final stage
FROM alpine:3.21
WORKDIR /app

# Copy binaries and migrations
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/internal/database/migrations ./migrations
# For seeding example data
COPY --from=builder /app/seed .

RUN apk update && apk add --no-cache go

# Create entrypoint script
RUN echo $'#!/bin/sh\n\
./migrate -action up\n\
exec ./server' > entrypoint.sh && \
    chmod +x entrypoint.sh

EXPOSE 8080
CMD ["./entrypoint.sh"]