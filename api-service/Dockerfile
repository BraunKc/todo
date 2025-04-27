FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/main ./cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/main /main

RUN mkdir -p /api-service/config
COPY --from=builder /app/config/config.yml /api-service/config/config.yml

EXPOSE 8080

CMD ["/main"]