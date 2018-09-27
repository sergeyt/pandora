FROM golang:latest

ENV GOPATH /go
RUN go get -u -v github.com/cosmtrek/air
WORKDIR /go/src/github.com/sergeyt/pandora

ENTRYPOINT ["/go/bin/air"]
