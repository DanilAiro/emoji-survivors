FROM golang:latest AS builder

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /build

COPY shared/jwt ./shared/jwt
COPY game-service ./game-service

WORKDIR /build/game-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/bin/game-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/bin/game-service .
RUN chmod +x ./game-service
EXPOSE 8082
CMD ["./game-service"]
