FROM golang:1.15-alpine

ENV GOPATH /go
ENV DOCKERIZE_VERSION v0.6.1

RUN apk add --no-cache openssl curl

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN wget https://github.com/cosmtrek/air/raw/master/bin/linux/air && chmod +x air && mv air /usr/local/bin

COPY docker/app/main.sh /main.sh

WORKDIR /pandora

ENTRYPOINT ["/main.sh"]

HEALTHCHECK --interval=10s --timeout=10s CMD /usr/bin/curl --fail http://localhost:4201/api/health || exit 1
