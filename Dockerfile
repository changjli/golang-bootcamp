# Dockerfile (place in project root)

# --- Builder Stage ---
# Use an official Golang image as the builder.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy all Go module and workspace files first.
# This helps Docker cache the dependency layer.
COPY go.work go.work.sum ./
COPY core-service/go.mod core-service/go.sum ./core-service/
COPY events/go.mod ./events/
COPY utils/go.mod  ./utils/

# Copy the entire source code of your project.
COPY . .

# Build the core-service application.
# The 'go build' command will automatically resolve and download
# any necessary dependencies before compiling the code.
# -o specifies the output binary.
RUN CGO_ENABLED=0 go build -o /app/core-service-binary ./core-service

# --- Final Stage ---
# Use a minimal base image for the final container.
FROM alpine:latest

RUN apk add --no-cache tzdata

# Set the working directory
WORKDIR /app/

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/core-service-binary .

COPY ./core-service/.env .

# Expose the port the application runs on.
# Change 8080 to whatever port your core-service listens on.
EXPOSE 8080

# Command to run the executable when the container starts.
CMD ["./core-service-binary"]