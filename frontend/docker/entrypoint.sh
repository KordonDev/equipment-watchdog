#!/bin/bash

echo $URL
sed -i "s/.*BASE_URL.*/window.BASE_URL=\"$(echo $URL | sed -En 's/\//\\\//pg')\";/g" /var/www/index.html

nginx -g "daemon off;"
