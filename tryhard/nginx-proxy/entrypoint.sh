#!/bin/sh

# Substitute Template
sed "s/PROXY_SERVICE_NAME/$PROXY_SERVICE_NAME/g" default.conf.template >> default.conf

# Start nginx 
nginx -g 'daemon off;'