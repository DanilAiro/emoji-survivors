FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY shared/jwt ./shared/jwt
COPY game-service ./game-service

WORKDIR /build/game-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /build/bin/game-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /build/bin/game-service ./game-service

EXPOSE 8082

ENTRYPOINT ["./game-service"]