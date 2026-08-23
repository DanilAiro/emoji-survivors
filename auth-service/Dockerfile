FROM golang:latest AS builder

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /build

COPY shared/jwt ./shared/jwt
COPY auth-service ./auth-service

WORKDIR /build/auth-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/bin/auth-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/bin/auth-service .
RUN chmod +x ./auth-service
EXPOSE 8081
CMD ["./auth-service"]
