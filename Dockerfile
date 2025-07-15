FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

# Create a new user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Switch to the new user
USER appuser

CMD ["./app"]