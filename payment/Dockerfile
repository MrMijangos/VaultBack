FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o vault-payment .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/vault-payment ./vault-payment
EXPOSE 8005
CMD ["./vault-payment"]
