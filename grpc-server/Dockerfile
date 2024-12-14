# Use a minimal Go base image
FROM golang:1.23.2 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY ./grpc-server/go.mod ./grpc-server/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code (project root files)
COPY ./grpc-server /app

# Build the application from the main directory
RUN go build -o grpc-server ./main

# Use a minimal runtime base image for the final stage
FROM debian:bullseye-slim

# Set the working directory in the final runtime image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/grpc-server ./main

# Expose the gRPC server port
EXPOSE 50051

# Run the binary (no .go files here, only the binary)
CMD ["./grpc-server"]
