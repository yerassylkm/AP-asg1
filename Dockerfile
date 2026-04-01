FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o payment-app cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/payment-app .

EXPOSE 8081

CMD ["./payment-app"]