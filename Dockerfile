FROM golang:1.20.5-bullseye AS Build

RUN mkdir /build
ADD . /build
WORKDIR /build
# RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

RUN CGO_ENABLED=1 go build -o equipment-watchdog ./main.go #-mod=vendor ./main.go

FROM scratch

VOLUME /etc/equipment-watchdog
#COPY --from=Build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=Build /build/equipment-watchdog /app/
WORKDIR /app

ENTRYPOINT ["./equipment-watchdog"]
