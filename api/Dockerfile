FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o vault-api .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/vault-api ./vault-api
COPY --from=builder /app/init.sql ./init.sql
EXPOSE 8080
CMD ["./vault-api"]
