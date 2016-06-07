VERSION = $(shell git tag | tail -n 1)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC := $(wildcard *.go) $(wildcard src/*/*.go) $(wildcard src/*/*/*.go)

all: gcredstash

gcredstash: go-get $(SRC)
	GOPATH=$(RUNTIME_GOPATH) go build

go-get:
	go get github.com/mitchellh/cli
	go get github.com/aws/aws-sdk-go
	go get github.com/ryanuber/go-glob

clean:
	rm -f gcredstash *.gz

package: clean gcredstash
	gzip -c gcredstash > gcredstash-$(VERSION)-$(GOOS)-$(GOARCH).gz

deb:
	dpkg-buildpackage -us -uc
