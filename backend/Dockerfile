FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o net-sentry-api ./cmd/api

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/net-sentry-api /app/net-sentry-api

EXPOSE 8080

ENTRYPOINT ["/app/net-sentry-api"]