#!/usr/bin/env bash

set -a
source .env
go test
set +a
