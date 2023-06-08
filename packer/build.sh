#!/bin/bash

[ ! -d "../build" ] && mkdir ../build
[ -f "../build/imgo" ] && rm ../build/imgo
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/imgo ../main.go 
