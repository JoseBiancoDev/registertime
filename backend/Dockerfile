FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod file
COPY go.mod ./

# Copy the rest of the code
COPY . .

# Tidy and download dependencies to generate go.sum
RUN go mod tidy && go mod download

# Build the application
RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
