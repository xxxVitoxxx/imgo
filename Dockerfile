FROM golang:1.19 as builder

LABEL stage=latest

WORKDIR /app

ADD . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o imgo main.go

EXPOSE 8090

CMD ["/app/imgo"]
