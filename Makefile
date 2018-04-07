export GOPATH := $(PWD)
export GOBIN := $(GOPATH)/bin
export PACKAGES := $(shell env GOPATH=$(GOPATH) go list ./src/...)
export REVISION := $(shell git describe --exact-match --tags $(git log -n1 --pretty='%h') || git rev-parse --verify --short HEAD || echo ${REVISION})

all: clean predependencies dependencies build

clean:
	rm -vf bin/*

build: build-macos build-linux build-windows

build-macos: build-macos-amd64 build-macos-i386

build-linux: build-linux-amd64 build-linux-i386

build-windows: build-windows-amd64 build-windows-i386

build-macos-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=${REVISION}" -o bin/svcproxy-darwin-amd64 svcproxy/cmd

build-macos-i386:
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -ldflags "-X main.Version=${REVISION}" -o bin/svcproxy-darwin-i386 svcproxy/cmd

build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=${REVISION}" -o bin/svcproxy-linux-amd64 svcproxy/cmd

build-linux-i386:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-X main.Version=${REVISION}" -o bin/svcproxy-linux-i386 svcproxy/cmd

build-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=${REVISION}" -o bin/svcproxy-windows-amd64.exe svcproxy/cmd

build-windows-i386:
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags "-X main.Version=${REVISION}" -o bin/svcproxy-windows-i386.exe svcproxy/cmd

dependencies:
	cd src && trash

predependencies:
	go get -u github.com/rancher/trash

sign:
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/svcproxy-darwin-amd64.sig 				bin/svcproxy-darwin-amd64
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/svcproxy-darwin-i386.sig 				bin/svcproxy-darwin-i386
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/svcproxy-linux-amd64.sig 				bin/svcproxy-linux-amd64
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/svcproxy-linux-i386.sig 					bin/svcproxy-linux-i386
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/svcproxy-windows-amd64.exe.sig 	bin/svcproxy-windows-amd64.exe
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/svcproxy-windows-i386.exe.sig 		bin/svcproxy-windows-i386.exe

test:
	go test ./src/...

verify:
	gpg --verify bin/svcproxy-darwin-amd64.sig 				bin/svcproxy-darwin-amd64
	gpg --verify bin/svcproxy-darwin-i386.sig 				bin/svcproxy-darwin-i386
	gpg --verify bin/svcproxy-linux-amd64.sig 				bin/svcproxy-linux-amd64
	gpg --verify bin/svcproxy-linux-i386.sig 					bin/svcproxy-linux-i386
	gpg --verify bin/svcproxy-windows-amd64.exe.sig 	bin/svcproxy-windows-amd64.exe
	gpg --verify bin/svcproxy-windows-i386.exe.sig 		bin/svcproxy-windows-i386.exe
