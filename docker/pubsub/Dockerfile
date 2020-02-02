FROM golang:1.13-alpine

ENV GOPATH /go
ENV DOCKERIZE_VERSION v0.6.1

RUN apk add --no-cache openssl curl

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN apk add --no-cache git \
    && go get -u -v github.com/gocontrib/pubsub/cmd/pubsubd \
    && apk del git

WORKDIR /go/src/github.com/gocontrib/pubsub/cmd/pubsubd

ENTRYPOINT ["dockerize", "--wait", "tcp://nats:4222", "/go/bin/pubsubd"]

HEALTHCHECK --interval=10s --timeout=10s CMD /usr/bin/curl --fail http://localhost:4302/api/pubsub/health || exit 1
