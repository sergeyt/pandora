FROM golang:1.16-alpine AS build

ENV GOPATH /go
ENV DOCKERIZE_VERSION v0.6.1

RUN apk add --no-cache git

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN go get -u -v github.com/gocontrib/pubsub/cmd/pubsubd

FROM alpine

RUN apk add --no-cache openssl curl

COPY --from=build /usr/local/bin/dockerize /usr/local/bin
COPY --from=build /go/bin/pubsubd /usr/local/bin

ENTRYPOINT ["dockerize", "--wait", "tcp://nats:4222", "/usr/local/bin/pubsubd"]

HEALTHCHECK --interval=60s --timeout=60s CMD /usr/bin/curl --fail http://localhost:4302/api/pubsub/health || exit 1
