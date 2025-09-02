# --- Build Stage ---
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Compile the application into a static binary
# CGO_ENABLED=0 is important for minimal images like Alpine
RUN CGO_ENABLED=0 go build -o /oasnake main.go

# --- Final Stage ---
FROM alpine:latest

# Copy the compiled binary from the build stage
COPY --from=builder /oasnake /usr/local/bin/oasnake

# Set the binary as the container's entrypoint.
# The container will behave like the oasnake executable.
ENTRYPOINT ["oasnake"]

# Default command (e.g., display help)
CMD ["--help"]
