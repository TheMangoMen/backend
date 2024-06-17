# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22.4 as builder
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source to the container
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app


# Run stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary file from the build stage
COPY --from=builder /app/main .

# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["./main"]

