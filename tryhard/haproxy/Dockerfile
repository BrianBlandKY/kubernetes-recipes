FROM centos:centos7

# HTTP Port
EXPOSE 80/tcp
# HTTPS Port
EXPOSE 443/tcp

RUN yum update -y && \
    yum install -y epel-release && \
    yum install -y gcc pcre-static pcre-devel && \
    yum install -y wget nano haproxy && \ 
    yum group install "Development Tools" -y && \
    yum clean all

RUN mkdir -p /usr/local/src/haproxy

WORKDIR /usr/local/src/haproxy

RUN wget http://www.haproxy.org/download/1.8/src/haproxy-1.8.4.tar.gz -O haproxy.tar.gz
RUN tar xzvf haproxy.tar.gz -C .

WORKDIR /usr/local/src/haproxy/haproxy-1.8.4
RUN make TARGET=linux2628
RUN make install

# Copy over config
COPY haproxy.cfg /etc/haproxy

# Let's Encrypt Stuffs

# ... 

ENTRYPOINT ["haproxy", "-f", "/etc/haproxy/haproxy.cfg", "-db", "-W"]

# Testing
# ENTRYPOINT ["tail", "-f", "/dev/null"]

# Build
# docker build -t haproxy-img:1.0.0 -f Dockerfile .

# Run
# docker run -tdi -p 80:80 -p 443:443 --name haproxy haproxy-img
