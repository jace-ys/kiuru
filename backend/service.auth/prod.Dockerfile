FROM golang:1.13 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o service

FROM alpine:3.10
WORKDIR /src
COPY --from=builder /src/service /bin/service
COPY --from=builder /src/serve /bin/serve
CMD ["serve"]
