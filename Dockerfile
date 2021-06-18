FROM ubuntu:xenial-20210429

RUN apt-get update && apt-get -y install git wget unzip ruby

RUN apt-get -y install build-essential libcurl4-openssl-dev libxml2 libxml2-dev libxslt1-dev ruby-dev  libgmp-dev zlib1g-dev

#Installs at /usr/local/bin/wpscan
RUN gem install wpscan

RUN wget https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz \
        && tar -xvf go1.15.2.linux-amd64.tar.gz \
        && mv go /usr/local \
        && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.5/protoc-3.15.5-linux-x86_64.zip \
        && unzip protoc-3.15.5-linux-x86_64.zip -d $HOME/.local

ADD . /appsrc/

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH
ENV GO111MODULE=on
ENV PATH=$PATH:/root/go/bin:/root/.local/bin

RUN cd /appsrc/
    && go build .


CMD ["/appsrc/wpscandocker"]
