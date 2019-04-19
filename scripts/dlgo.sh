#!/usr/bin/env bash

GO_VERSION=1.12.4
FILE=go${GO_VERSION}.linux-amd64.tar.gz

curl -O https://storage.googleapis.com/golang/${FILE}
tar -xvf ${FILE}
rm -rf ${HOME}/tools/go
mv go ${HOME}/tools
rm -f ${FILE}
