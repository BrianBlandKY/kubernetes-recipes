FROM nginx:1.13.8

EXPOSE 8080/tcp

ENV PROXY_SERVICE_NAME localhost

WORKDIR /etc/nginx/conf.d

# Remove default configuration
RUN rm default.conf

# Copy configuration template
COPY default.conf.template .

# Copy entrypoint script
COPY entrypoint.sh .

RUN chmod 700 entrypoint.sh

# Using ENV variables we can dynamically set the nginx config at run-time
# One image can be used for multiple applications
ENTRYPOINT ["./entrypoint.sh"]

# Testing
# ENTRYPOINT ["tail", "-f", "/dev/null"]

# Build
# docker build -t nothingsbland-img:1.0.0 -f Dockerfile .

# Run
# docker run -tdi -p 8080:8080 -e "PROXY_SERVICE_NAME=192.168.1.8" --name nginx nginx-img

# Terminal 
# docker exec -ti nginx /bin/bash
