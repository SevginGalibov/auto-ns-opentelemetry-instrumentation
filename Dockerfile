FROM golang:1.20-alpine AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN apk add --no-cache git
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -a -o /out/operator ./

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/operator /usr/local/bin/operator
RUN mkdir -p /etc/instrumentation
WORKDIR /
ENTRYPOINT ["/usr/local/bin/operator"]
