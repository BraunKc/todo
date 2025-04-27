FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/main ./cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/main /main

EXPOSE 8282

CMD ["/main"]