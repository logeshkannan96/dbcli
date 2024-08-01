BINARY_NAME=dbcli
VERSION=1.0.0

build:
	go build -o ${BINARY_NAME} -v ./cmd

install:
	go install -v ./cmd

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test -v ./...

run:
	go run ./cmd

.PHONY: build install clean test run