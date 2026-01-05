#!/bin/sh

# Escape characters that would break the sed replacement (/, &)
ESCAPED_URL=$(printf '%s' "$URL" | sed -e 's/[\/&]/\\&/g')

# Replace the BASE_URL JS snippet in every .html file in /var/www (non-recursive)
for f in /var/www/*.html; do
  if [ -f "$f" ]; then
    echo "Updating BASE_URL in $f to $URL"
    sed -i "s/.*BASE_URL.*/window.BASE_URL=\"$ESCAPED_URL\";/" "$f"
  fi
done

caddy run --config /etc/caddy/Caddyfile --adapter caddyfile