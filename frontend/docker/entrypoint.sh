#!/bin/sh

echo $URL
sed -i "s/.*BASE_URL.*/window.BASE_URL=\"$(echo $URL | sed -En 's/\//\\\//pg')\";/g" /var/www/index.html

caddy run --config /etc/caddy/Caddyfile --adapter caddyfile