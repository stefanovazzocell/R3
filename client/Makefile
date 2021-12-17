.PHONY: build clean release

@:
	go fmt
	go clean
	go run .

build:
	go build -ldflags "-w -s"

clean:
	go fmt
	go clean

release:
	@echo "Building for Linux"
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o rocket_linux_amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o rocket_linux_arm64
	@echo "Building for Windows"
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o rocket_windows_amd64.exe
	GOOS=windows GOARCH=arm go build -ldflags "-w -s" -o rocket_windows_arm.exe
	@echo "Building for Mac OS"
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o rocket_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-w -s" -o rocket_darwin_arm64