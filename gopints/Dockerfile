# Build on the runner's native platform; cross-compile to TARGETOS/TARGETARCH.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG TARGETVARIANT
ARG VERSION=dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOARM=$([ "${TARGETVARIANT}" = "v7" ] && echo "7" || echo "") \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -X main.version=${VERSION}" \
    -o /server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /server .
VOLUME ["/app/data"]
EXPOSE 8080 9876/udp
CMD ["./server"]
