#!/bin/sh

CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o server main.go

docker build . -t xeno:1.0
