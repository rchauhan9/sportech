FROM golang:1.18-alpine as builder

WORKDIR /build
COPY ../../go.mod go.sum ./
RUN go mod download

# TODO: refactor to copy only files you need
COPY ../.. .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o sportech main.go

FROM busybox:latest
WORKDIR /app
RUN mkdir config
COPY --from=builder build/sportech .
COPY --from=builder build/config ./config

CMD ["./sportech"]
