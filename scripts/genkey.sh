#!/usr/bin/env bash

head -c 12 /dev/urandom | shasum | cut -d' ' -f1
