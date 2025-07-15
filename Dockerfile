FROM golang:1.24.3-alpine

WORKDIR /app

# Create a new user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Switch to the new user
USER appuser

#  Copy go.mod and go.sum first to leverage Docker layer caching for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o app

CMD ["./app"]