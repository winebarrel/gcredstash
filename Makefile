SHELL:=/bin/bash
VERSION:=v0.3.5
GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)
SRC:=$(wildcard *.go) $(wildcard src/*/*.go) $(wildcard src/*/*/*.go)
TEST_SRC:=$(wildcard src/gcredstash/*_test.go)
CMD_TEST_SRC:=$(wildcard src/gcredstash/command/*_test.go)

UBUNTU_IMAGE=docker-go-pkg-build-ubuntu-trusty
UBUNTU_CONTAINER_NAME=docker-go-pkg-build-ubuntu-trusty-$(shell date +%s)
CENTOS_IMAGE=docker-go-pkg-build-centos6
CENTOS_CONTAINER_NAME=docker-go-pkg-build-centos6-$(shell date +%s)

all: gcredstash

gcredstash: $(SRC)
	CGO_ENABLED=0 go build -a -ldflags "-w -s" -tags netgo -installsuffix netgo -o gcredstash
ifeq ($(GOOS),linux)
	[[ "`ldd gcredstash`" =~ "not a dynamic executable" ]] || exit 1
endif

test: $(TEST_SRC) $(CMD_TEST_SRC)
	go test -v $(TEST_SRC)
	go test -v $(CMD_TEST_SRC)

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

mock:
	go get github.com/golang/mock/mockgen
	go get github.com/aws/aws-sdk-go
	mockgen -package mockaws -destination src/mockaws/dynamodbmock.go github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI
	mockgen -package mockaws -destination src/mockaws/kmsmock.go github.com/aws/aws-sdk-go/service/kms/kmsiface KMSAPI

tag:
ifdef FORCE
	git tag $(VERSION) -f
else
	git tag $(VERSION)
endif
