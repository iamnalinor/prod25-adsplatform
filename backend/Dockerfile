# Dockerfile was generated from
# https://github.com/lodthe/dockerfiles/blob/main/go/Dockerfile

FROM golang:1.23-alpine3.21 AS builder

# Setup base software for building an app.
RUN apk update && apk add ca-certificates git gcc g++ libc-dev binutils

WORKDIR /opt

# Download dependencies.
COPY go.mod go.sum ./
RUN go mod download && go install github.com/swaggo/swag/cmd/swag@latest && go mod verify

# Copy application source.
COPY . .

# Generate the docs.
RUN swag init

# Build the application.
RUN go build -o bin/application main.go

# Prepare executor image.
FROM alpine:3.21 AS runner

RUN apk update && apk add ca-certificates libc6-compat openssh bash && rm -rf /var/cache/apk/*

WORKDIR /opt

COPY --from=builder /opt/bin/application ./

# Run the application.
CMD ["./application"]
HEALTHCHECK --interval=5s --timeout=5s CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1
