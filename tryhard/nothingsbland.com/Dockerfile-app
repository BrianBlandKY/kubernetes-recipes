FROM centos:centos7

# Web Port
EXPOSE 8080/tcp 

VOLUME ["/apps/nothingsbland"]

ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH $PATH:/go/bin
ENV PATH=$PATH:/usr/local/go/bin
ENV container=docker

# General Updates and Installs
RUN yum install -y epel-release && \
    yum update -y && \
    yum install -y wget git nano && \
    yum clean all

# Go 1.9.2
WORKDIR /tmp
RUN wget https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz

# Make volume log directory
RUN mkdir -p "/apps/nothingsbland/logs"

# Download NothingsBland.com Release v1.0.0
RUN mkdir -p /go/src/NothingsBland.com
RUN wget https://github.com/BrianBlandKY/NothingsBland.com/archive/v1.0.0.tar.gz && \
    tar -C /go/src/NothingsBland.com -xzf v1.0.0.tar.gz
RUN mv /go/src/NothingsBland.com/NothingsBland.com-1.0.0/* /go/src/NothingsBland.com/
RUN rm -fr /go/src/NothingsBland.com/NothingsBland.com-1.0.0

WORKDIR /go/src/NothingsBland.com

# Install Go Dep
RUN go get -u github.com/golang/dep/cmd/dep

# Install Application Dependencies
RUN dep ensure

# Build and Run minifier
WORKDIR /go/src/NothingsBland.com/minifier
RUN go build
RUN ./minifier

# Build Web
WORKDIR /go/src/NothingsBland.com/web
RUN go build

ENTRYPOINT [ "./web", "--config", "app.prod.yaml" ]

#ENTRYPOINT ["tail", "-f", "/dev/null"]

# Build
# docker build -t nothingsbland-img -f Dockerfile .

# Run
# docker run -tdi -p 8080:8080 --name nothingsbland nothingsbland-img
# docker run -tdi -p 8080:8080 -v /Users/bland/Development/Go/src/NothingsBland.com:/apps/nothingsbland --name nothingsbland nothingsbland-img

# Exec
# docker exec -ti nothingsbland sh -c "go run main.go"
# docker exec -ti nothingsbland sh -c "go-wrapper run"

# Terminal 
# docker exec -ti nothingsbland /bin/bash

# Stop Containers
# docker stop $(docker ps -aq)

# Remove Containers
# docker rm $(docker ps -aq)

# Remove Images
# docker rmi $(docker images -q)

# Helpful Commands
# docker stop $(docker ps -aq) && docker rm $(docker ps -aq) && docker rmi $(docker images -q)
# docker exec -t -i 50f331760ba7 /bin/bash
# docker start -a -i `docker ps -q -l`