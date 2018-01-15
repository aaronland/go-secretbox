CWD=$(shell pwd)
GOPATH := $(CWD)

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src; then mkdir src; fi
	mkdir -p src/github.com/aaronland/go-secretbox
	cp -r config src/github.com/aaronland/go-secretbox/
	cp -r salt src/github.com/aaronland/go-secretbox/
	cp *.go src/github.com/aaronland/go-secretbox/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "golang.org/x/crypto/nacl/secretbox"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/crypto/scrypt"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/crypto/ssh/terminal"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt config/*.go
	go fmt salt/*.go
	go fmt *.go

# PLEASE CLEAN ALL OF THIS UP... (20180115/thisisaaronland)

bin: 	self
	if test ! -d bin/$(OS); then mkdir -p bin/$(OS); fi
	@GOPATH=$(GOPATH) go build -o bin/$(OS)/secretbox cmd/secretbox.go
	@GOPATH=$(GOPATH) go build -o bin/$(OS)/saltshaker cmd/saltshaker.go

darwin: self
	if test ! -d bin/darwin; then mkdir -p bin/darwin; fi
	@GOPATH=$(GOPATH) GOOS=darwin GOARCH=386 go build -o bin/darwin/secretbox cmd/secretbox.go
	@GOPATH=$(GOPATH) GOOS=darwin GOARCH=386 go build -o bin/darwin/saltshaker cmd/saltshaker.go

linux: self
	if test ! -d bin/linux; then mkdir -p bin/linux; fi
	@GOPATH=$(GOPATH) GOOS=linux GOARCH=386 go build -o bin/linux/secretbox cmd/secretbox.go
	@GOPATH=$(GOPATH) GOOS=linux GOARCH=386 go build -o bin/linux/saltshaker cmd/saltshaker.go

android: self
	if test ! -d bin/android; then mkdir -p bin/android; fi
	@GOPATH=$(GOPATH) GOOS=android GOARCH=386 go build -o bin/android/secretbox cmd/secretbox.go
	@GOPATH=$(GOPATH) GOOS=android GOARCH=386 go build -o bin/android/saltshaker cmd/saltshaker.go
