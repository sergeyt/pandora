#!/bin/sh -e

cd /pyadmin && gunicorn admin:app &

caddy -agree --conf /etc/caddy/Caddyfile
