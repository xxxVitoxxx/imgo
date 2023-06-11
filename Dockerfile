FROM golang:1.19 AS builder

LABEL stage=latest

WORKDIR /app

ADD . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o imgo main.go

FROM alpine:latest

COPY --from=builder /app/imgo /bin/imgo
COPY config.toml .

EXPOSE 80

CMD ["/bin/imgo"]
