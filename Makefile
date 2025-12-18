
BINARY=operator
IMAGE?=ghcr.io/sevgingalibov/auto-ns-opentelemetry-instrumentation:latest

.PHONY: build docker-build docker-buildx
build:
	CGO_ENABLED=0 GOOS=linux go build -o ${BINARY} ./

docker-build:
	docker build -t ${IMAGE} .

docker-buildx:
	docker buildx build --platform linux/amd64,linux/arm64 -t ${IMAGE} --push .
