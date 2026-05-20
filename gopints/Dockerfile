FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o kegerator-server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /app/kegerator-server .
VOLUME ["/app/data"]
EXPOSE 8080
EXPOSE 9876/udp
CMD ["./kegerator-server"]
