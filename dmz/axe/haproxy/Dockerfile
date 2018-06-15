FROM centos:centos7

# HTTP Port
EXPOSE 80/tcp
# HTTPS Port
EXPOSE 443/tcp

VOLUME /etc/letsencrypt

ENV CERTS nothingsbland.com, tryhard.tv
ENV EMAIL brianbland42@gmail.com

ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH $PATH:/go/bin
ENV PATH=$PATH:/usr/local/go/bin

ENV LUA_VERSION 5-3-4
ENV LUA_VERSION_SHORT 53

ENV HAPROXY_MAJOR 1.8             
ENV HAPROXY_VERSION 1.8.8

# Remove Base Repo (random issue)
# RUN mv /etc/yum.repos.d/CentOS-Base.repo /tmp

RUN cat /etc/yum.repos.d/CentOS-Base.repo
RUN cat /etc/resolv.conf

#RUN ping https://www.google.com

# System Updates
RUN yum update -y && \
    yum groups mark install "Development Tools" && \
    yum groups mark convert "Development Tools" && \
    yum groupinstall -y "Development Tools" && \
    yum install -y \
        gcc gcc-c++ make \
        readline-devel \
        openssl-devel \
        ncurses-devel \
        glibc-devel  \
        kernel-devel \
        pcre-static \
        pcre-devel \
        curl \
        which \
        dnsmasq \
        ca-certificates \
        python-ndg-httpsclient \
        wget \
        nano \
        haproxy \
        cronie \
        certbot \
        python-setuptools && \ 
    yum clean all

# EPEL
RUN cd /tmp && wget https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
RUN yum install -y /tmp/epel-release-latest-*.noarch.rpm

# Lua Install
RUN cd /usr/src \
    && curl -R -O https://github.com/lua/lua/archive/v${LUA_VERSION}.tar.gz \
    && tar zxf v${LUA_VERSION}.tar.gz \
    && rm v${LUA_VERSION}.tar.gz \
    && cd v${LUA_VERSION} \
    && make linux \
&& make INSTALL_TOP=/opt/v${LUA_VERSION_SHORT} install

# HAProxy Install
RUN cd / \
    && wget http://www.haproxy.org/download/${HAPROXY_MAJOR}/src/haproxy-${HAPROXY_VERSION}.tar.gz -O haproxy.tar.gz \
    && mkdir -p /usr/src/haproxy \
    && tar -xzf haproxy.tar.gz -C /usr/src/haproxy --strip-components=1 \
    && rm haproxy.tar.gz \
    && make -C /usr/src/haproxy \
        TARGET=linux2628 \
        USE_PCRE=1 \
        PCREDIR= \
		USE_OPENSSL=1 \
		USE_ZLIB=1 \
		USE_LUA=yes \
        LUA_LIB=/opt/lua53/lib/ \
        LUA_INC=/opt/lua53/include/ \
        LDFLAGS=-ldl \
		all \
		install-bin \
    && mkdir -p /usr/local/etc/haproxy \
    && cp -R /usr/src/haproxy/examples/errorfiles /usr/local/etc/haproxy/errors \
    && rm -rf /usr/src/haproxy

# Install Go
WORKDIR /tmp
RUN wget https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz

COPY haproxy.https.cfg /usr/local/etc/haproxy/haproxy.https.cfg
COPY haproxy.http.cfg /usr/local/etc/haproxy/haproxy.http.cfg

# Get into shawshank-haproxy directory
RUN mkdir -p /go/src/shawshank-haproxy
COPY shawshank-haproxy/entrypoint.go /go/src/shawshank-haproxy/entrypoint.go
WORKDIR /go/src/shawshank-haproxy
RUN go build

ENTRYPOINT ["./shawshank-haproxy"]
# ENTRYPOINT ["tail", "-f", "/dev/null"]

# Build
# docker build -t haproxy-img:1.1.0 -f Dockerfile .

# Run
# docker run -tdi -p 80:80 -p 443:443 -v /vault/apps/letsencrypt:/etc/letsencrypt --name haproxy haproxy-img:1.0.5
# docker run -tdi -p 80:80 --name haproxy haproxy-img:1.0.0


# Remove untagged images
# docker rmi $(docker images -a | grep "^<none>" | awk '{print $3}')
# Need Troubleshooting Steps
