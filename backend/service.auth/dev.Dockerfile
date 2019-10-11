FROM golang:1.13 AS builder
WORKDIR /service.auth
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o service

FROM alpine:3.10
WORKDIR /service
COPY --from=builder /service.auth/serve-dev .
COPY --from=builder /service.auth/service .
CMD ["./serve-dev"]
