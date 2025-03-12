clean:
	rm -Rfv bin
	mkdir bin

build: clean
	go build -o bin/mailcow-prometheus-exporter cmd/main.go

build-all: clean
	GOOS="linux"   GOARCH="amd64"       go build -o bin/mailcow-prometheus-exporter__linux-amd64 cmd/main.go
	GOOS="linux"   GOARCH="arm" GOARM=6 go build -o bin/mailcow-prometheus-exporter__linux-armv6 cmd/main.go
	GOOS="linux"   GOARCH="arm" GOARM=7 go build -o bin/mailcow-prometheus-exporter__linux-armv7 cmd/main.go
	GOOS="linux"   GOARCH="arm64"       go build -o bin/mailcow-prometheus-exporter__linux-arm64 cmd/main.go
	GOOS="darwin"  GOARCH="amd64"       go build -o bin/mailcow-prometheus-exporter__macos-amd64 cmd/main.go
	GOOS="darwin"  GOARCH="arm64"       go build -o bin/mailcow-prometheus-exporter__macos-arm64 cmd/main.go
	GOOS="windows" GOARCH="amd64"       go build -o bin/mailcow-prometheus-exporter__win-amd64   cmd/main.go
	GOOS="windows" GOARCH="arm64"       go build -o bin/mailcow-prometheus-exporter__win-arm64   cmd/main.go

docker:
	docker build . -t themailcow/prometheus-exporter
