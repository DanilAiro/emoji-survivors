FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY shared/jwt ./shared/jwt
COPY auth-service ./auth-service

WORKDIR /build/auth-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /build/bin/auth-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /build/bin/auth-service ./auth-service

EXPOSE 8081

ENTRYPOINT ["./auth-service"]