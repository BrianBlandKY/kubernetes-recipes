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

# ENTRYPOINT ["tail", "-f", "/dev/null"]
CMD ["haproxy", "-f", "/etc/haproxy/haproxy.cfg", "-db", "-W"]

# Build
# docker build -t nothingsbland-img -f Dockerfile .

# Run
# docker run -tdi -p 8080:8080 --name nothingsbland nothingsbland-img
# docker ruxn -tdi -p 8080:8080 -v /Users/bland/Development/Go/src/NothingsBland.com:/apps/nothingsbland --name nothingsbland nothingsbland-img


# Install and Configure Let's Encrypt
# Make this configurable for other installations: 
# - nothignsbland.com
# - quizzicaldad.com
# - tryhard.me
# Use base HAProxy CMD
