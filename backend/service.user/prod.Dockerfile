FROM golang:1.13 AS builder
WORKDIR /service.user
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o service

FROM alpine:3.10
WORKDIR /service
COPY --from=builder /service.user/service /bin/service
COPY --from=builder /service.user/serve /bin/serve
CMD ["serve"]
