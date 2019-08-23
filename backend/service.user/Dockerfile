FROM golang:1.12 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main

FROM alpine:3.10
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .
EXPOSE 80
CMD ["./main"]
