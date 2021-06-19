FROM quay.io/aptible/ruby:2.7-ubuntu-16.04
#Source: https://github.com/aptible/docker-ruby

RUN apt-get update && apt-get -y install git wget unzip curl jq

RUN apt-get -y install build-essential libcurl4-openssl-dev libxml2 libxml2-dev libxslt1-dev libgmp-dev zlib1g-dev

#Installs at /usr/local/bin/wpscan
RUN gem install wpscan

RUN wget https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz \
        && tar -xvf go1.15.2.linux-amd64.tar.gz \
        && mv go /usr/local

ADD . /appsrc/

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH
ENV GO111MODULE=on
ENV PATH=$PATH:/root/go/bin:/root/.local/bin

RUN cd /appsrc/ \
    && chmod 777 -R . \
    && go build .


CMD ["/appsrc/wpscandocker"]