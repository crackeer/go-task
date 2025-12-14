# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Build the application
ENV GOPROXY https://goproxy.cn,direct
RUN GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Copy any other necessary files (if needed)
# COPY .env .

# Expose the default port
EXPOSE 80

# Set environment variables
ENV PORT=80

# Run the application
CMD ["./main"]