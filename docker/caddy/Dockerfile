# Builder based on https://github.com/abiosoft/caddy-docker/blob/master/builder/Dockerfile
FROM golang:1.13-alpine as builder

ARG version="1.0.4"
ARG plugins=cors,realip,expires,cache,jwt,git
ARG enable_telemetry="true"

RUN apk add --no-cache git gcc musl-dev

# process wrapper
RUN go get -v github.com/abiosoft/parent

# build caddy
COPY ./docker/caddy/build.sh /usr/bin/caddy_build.sh
RUN VERSION=${version} PLUGINS=${plugins} /bin/sh /usr/bin/caddy_build.sh

# Dist image
FROM alpine:latest

ARG version="1.0.3"
LABEL caddy_version="$version"

ENV HTTP_PORT 4200
# Let's Encrypt Agreement
ENV ACME_AGREE="false"
# Telemetry Stats
ENV ENABLE_TELEMETRY="$enable_telemetry"

# install deps
RUN apk update \
    && apk add --no-cache ca-certificates git mailcap openssh-client tzdata \
    && addgroup -S -g 1001 app \
    && adduser -S -D -h /app -u 1001 -G app app

# install caddy
COPY --from=builder /install/caddy /usr/bin/caddy

# install process wrapper
COPY --from=builder /go/bin/parent /bin/parent

# validate install
RUN /usr/bin/caddy -version
RUN /usr/bin/caddy -plugins

COPY config/caddy/Caddyfile /etc/Caddyfile

EXPOSE 80 443 4200
VOLUME ["/www"]
WORKDIR /www

ENTRYPOINT ["/bin/parent", "caddy"]
CMD ["--conf", "/etc/Caddyfile", "--log", "stdout", "--agree=$ACME_AGREE", "--email=stodyshev@gmail.com"]
