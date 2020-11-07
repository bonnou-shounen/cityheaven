#!/bin/sh

cd `dirname $0`

go build -o dist/cityheaven cmd/cityheaven/main.go
