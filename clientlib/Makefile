.PHONY: build clean install commit

@:
	go run .

build:
	go build -ldflags "-w -s"

clean:
	go fmt
	go clean

commit:
	go fmt
	go clean
	git add .
	git commit

test:
	go fmt
	go clean
	go test
