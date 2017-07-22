FROM centos:centos6
MAINTAINER Genki Sugawara <sgwr_dts@yahoo.co.jp>

USER root
WORKDIR /

RUN yum -y install yum-plugin-fastestmirror
RUN echo "include_only=.jp" >> /etc/yum/pluginconf.d/fastestmirror.conf
RUN yum install -y rpmdevtools make git gcc
RUN rpmdev-setuptree

ARG GO_VERSION=1.8.3
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/.go
ENV PATH $GOROOT/bin:$PATH
RUN wget -O- -q https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz | tar zxf - && \
    mv go /usr/local/
RUN mkdir $GOPATH
