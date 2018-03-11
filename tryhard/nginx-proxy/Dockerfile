FROM nginx:1.13.8

EXPOSE 8080/tcp

ENV SERVER_ADDRESS=localhost
ENV PROXY_PORT=8080

WORKDIR /etc/nginx/conf.d

# Remove default configuration
RUN rm default.conf

# Copy configuration template
COPY default.conf.template .

# Using ENV variables we can dynamically set the nginx config at run-time
# One image can be used for multiple applications
CMD sh -c "sed 's/SERVER_ADDRESS/$SERVER_ADDRESS/g; s/PROXY_PORT/$PROXY_PORT/g' default.conf.template >> default.conf && nginx -g 'daemon off;'"

# Build
# docker build -t nginx-img -f Dockerfile .

# Run
# docker run -tdi -p 80:80 -e "SERVER_ADDRESS=192.168.1.8" -e "PROXY_PORT=8080" --name nginx nginx-img

# Terminal 
# docker exec -ti nginx /bin/bash
