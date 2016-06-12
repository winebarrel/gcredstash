VERSION = $(shell git tag | tail -n 1)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC := $(wildcard *.go) $(wildcard src/*/*.go) $(wildcard src/*/*/*.go)
TEST_SRC := $(wildcard src/gcredstash/*_test.go)
CMD_TEST_SRC := $(wildcard src/gcredstash/command/*_test.go)

all: gcredstash

gcredstash: go-get $(SRC)
	GOPATH=$(RUNTIME_GOPATH) go build

test: go-get $(TEST_SRC) $(CMD_TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test $(TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test $(CMD_TEST_SRC)

go-get:
	go get github.com/mitchellh/cli
	go get github.com/aws/aws-sdk-go
	go get github.com/ryanuber/go-glob
	go get github.com/golang/mock/gomock

clean:
	rm -f gcredstash *.gz

package: clean gcredstash
	gzip -c gcredstash > gcredstash-$(VERSION)-$(GOOS)-$(GOARCH).gz

deb:
	dpkg-buildpackage -us -uc

mock:
	go get github.com/golang/mock/mockgen
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface/interface.go -destination src/mockaws/dynamodbmock.go -package mockaws
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/kms/kmsiface/interface.go -destination src/mockaws/kmsmock.go -package mockaws
