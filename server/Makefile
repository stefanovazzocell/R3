.PHONY: build clean test

@:
	go fmt
	go clean
	go run .

build:
	go build -ldflags "-w -s"

clean:
	go fmt
	go clean

test:
	go fmt
	go clean
	go test
	go clean