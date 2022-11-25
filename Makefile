CMD=cityheaven
SRC=$(shell find . -name "*.go")

.PHONY: lint test build clean

lint:
	golangci-lint run

test:
	go test -v ./...

build:
	go build -o ${CMD} cmd/${CMD}/main.go

clean:
	rm -f ${CMD}
