FROM alpine:3.18.4

ADD ./dist/equipment-watchdog /usr/local/bin/equipment-watchdog
EXPOSE 8080
ENV GIN_MODE=release

ENTRYPOINT ["equipment-watchdog"]
