.PHONY: build clean install commit

@:
	go run .

build:
	go build

clean:
	go fmt
	go clean

install:
	go get github.com/julienschmidt/httprouter
	go get github.com/garyburd/redigo/redis

commit:
	go fmt
	go clean
	git add .
	git commit
