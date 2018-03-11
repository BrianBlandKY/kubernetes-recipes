FROM nginx:1.13.8

EXPOSE 8080/tcp

ENV PROXY_SERVICE_NAME localhost

WORKDIR /etc/nginx/conf.d

# Remove default configuration
RUN rm default.conf

# Copy configuration template
COPY default.conf.template .

# Using ENV variables we can dynamically set the nginx config at run-time
# One image can be used for multiple applications
CMD sh -c "sed 's/PROXY_SERVICE_NAME/$PROXY_SERVICE_NAME/g' default.conf.template >> default.conf && nginx -g 'daemon off;'"

# Build
# docker build -t nginx-img -f Dockerfile .

# Run
# docker run -tdi -p 8080:8080 -e "PROXY_SERVICE_NAME=192.168.1.8" --name nginx nginx-img

# Terminal 
# docker exec -ti nginx /bin/bash
