FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/app /app/app
COPY .env /app/.env
EXPOSE 8080
CMD ["/app/app"]
