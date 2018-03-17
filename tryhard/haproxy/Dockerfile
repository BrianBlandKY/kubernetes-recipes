FROM centos:centos7

# HTTP Port
EXPOSE 80/tcp
# HTTPS Port
EXPOSE 443/tcp

ENV CERTS www.nothingsbland.com
ENV EMAIL brianbland42@gmail.com

# Software Versions
ENV SUPERVISOR_VERSION 3.3.4

ENV LUA_VERSION 5.3.4
ENV LUA_VERSION_SHORT 53

# ENV HAPROXY_MAJOR 1.8
# ENV HAPROXY_VERSION 1.8.4

ENV HAPROXY_MAJOR 1.7
ENV HAPROXY_VERSION 1.7.2

# System Updates
RUN yum update -y && \
    yum groups mark install "Development Tools" && \
    yum groups mark convert "Development Tools" && \
    yum groupinstall -y "Development Tools" && \
    yum install -y epel-release && \
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

# Supervisor Install
RUN wget https://github.com/Supervisor/supervisor/archive/${SUPERVISOR_VERSION}.tar.gz \
    && tar -xvf ${SUPERVISOR_VERSION}.tar.gz \
    && cd supervisor-${SUPERVISOR_VERSION} && python setup.py install

# Lua Install
RUN cd /usr/src \
    && curl -R -O http://www.lua.org/ftp/lua-${LUA_VERSION}.tar.gz \
    && tar zxf lua-${LUA_VERSION}.tar.gz \
    && rm lua-${LUA_VERSION}.tar.gz \
    && cd lua-${LUA_VERSION} \
    && make linux \
    && make INSTALL_TOP=/opt/lua${LUA_VERSION_SHORT} install

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

# See https://github.com/janeczku/haproxy-acme-validation-plugin
COPY acme/acme-http01-webroot.lua /usr/local/etc/haproxy
COPY acme/cert-renewal-haproxy.sh /

# Cron job for auto renewals
COPY crontab.txt /var/crontab.txt
RUN crontab /var/crontab.txt && chmod 600 /etc/crontab

COPY supervisord.conf /etc/supervisord.conf
COPY certs.sh /
COPY bootstrap.sh /
COPY docker-entrypoint.sh /

RUN chmod 700 /bootstrap.sh && \
    chmod 700 /usr/local/etc/haproxy/acme-http01-webroot.lua && \
    chmod 700 /cert-renewal-haproxy.sh && \
    chmod 700 /certs.sh && \
    chmod 700 /docker-entrypoint.sh

# Used for ACME renewal
RUN mkdir /jail && mkdir /etc/letsencrypt

# Copy over config
COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

# ENTRYPOINT ["/bootstrap.sh"]
# ENTRYPOINT ["haproxy", "-f", "/etc/haproxy/haproxy.cfg", "-db", "-W"]

# Testing
ENTRYPOINT ["tail", "-f", "/dev/null"]

# Build
# docker build -t haproxy-img:1.0.0 -f Dockerfile .

# Run
# docker run -tdi -p 80:80 -p 443:443 --name haproxy haproxy-img
# docker run -tdi -p 80:80 -p 443:443 \
# -e CERTS=my.domain,my.other.domain \
# -e EMAIL=my.email@my.domain haproxy:1.0.0
