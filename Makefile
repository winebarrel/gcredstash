SHELL:=/bin/bash
VERSION=v0.3.0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
RUNTIME_GOPATH=$(GOPATH):$(shell pwd)
SRC=$(wildcard *.go) $(wildcard src/*/*.go) $(wildcard src/*/*/*.go)
TEST_SRC=$(wildcard src/gcredstash/*_test.go)
CMD_TEST_SRC=$(wildcard src/gcredstash/command/*_test.go)

UBUNTU_IMAGE=docker-go-pkg-build-ubuntu-trusty
UBUNTU_CONTAINER_NAME=docker-go-pkg-build-ubuntu-trusty-$(shell date +%s)
CENTOS_IMAGE=docker-go-pkg-build-centos6
CENTOS_CONTAINER_NAME=docker-go-pkg-build-centos6-$(shell date +%s)

all: gcredstash

gcredstash: go-get $(SRC)
	GOPATH=$(RUNTIME_GOPATH) go build -a -tags netgo -installsuffix netgo -o gcredstash
ifeq ($(GOOS),linux)
	[[ "`ldd gcredstash`" =~ "not a dynamic executable" ]] || exit 1
endif

test: go-get $(TEST_SRC) $(CMD_TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test $(TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test $(CMD_TEST_SRC)

go-get:
	go get github.com/mitchellh/cli
	go get github.com/aws/aws-sdk-go
	go get github.com/ryanuber/go-glob
	go get github.com/golang/mock/gomock
	go get github.com/mattn/go-shellwords

clean:
	rm -f gcredstash{,.exe} *.gz *.zip
	rm -f pkg/*
	rm -f debian/gcredstash.debhelper.log
	rm -f debian/gcredstash.substvars
	rm -f debian/files
	rm -rf debian/gcredstash/

package: clean gcredstash
ifeq ($(GOOS),windows)
	zip gcredstash-$(VERSION)-$(GOOS)-$(GOARCH).zip gcredstash.exe
else
	gzip -c gcredstash > gcredstash-$(VERSION)-$(GOOS)-$(GOARCH).gz
endif

package\:linux:
	docker run --name $(UBUNTU_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) make -C /tmp/src package:linux:docker
	docker rm $(UBUNTU_CONTAINER_NAME)

package\:linux\:docker: package
	mv gcredstash-*.gz pkg/

deb:
	docker run --name $(UBUNTU_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) make -C /tmp/src deb:docker
	docker rm $(UBUNTU_CONTAINER_NAME)

deb\:docker: clean
	dpkg-buildpackage -us -uc
	mv ../gcredstash_* pkg/

docker\:build\:ubuntu-trusty:
	docker build -f docker/Dockerfile.ubuntu-trusty -t $(UBUNTU_IMAGE) .

rpm:
	docker run --name $(CENTOS_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(CENTOS_IMAGE) make -C /tmp/src rpm:docker
	docker rm $(CENTOS_CONTAINER_NAME)

rpm\:docker: clean
	cd ../ && tar zcf gcredstash.tar.gz src
	mv ../gcredstash.tar.gz /root/rpmbuild/SOURCES/
	cp gcredstash.spec /root/rpmbuild/SPECS/
	rpmbuild -ba /root/rpmbuild/SPECS/gcredstash.spec
	mv /root/rpmbuild/RPMS/x86_64/gcredstash-*.rpm pkg/
	mv /root/rpmbuild/SRPMS/gcredstash-*.src.rpm pkg/

docker\:build\:centos6:
	docker build -f docker/Dockerfile.centos6 -t $(CENTOS_IMAGE) .

mock: go-get
	go get github.com/golang/mock/mockgen
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface/interface.go -destination src/mockaws/dynamodbmock.go -package mockaws
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/kms/kmsiface/interface.go -destination src/mockaws/kmsmock.go -package mockaws

tag:
ifdef FORCE
	git tag $(VERSION) -f
else
	git tag $(VERSION)
endif
