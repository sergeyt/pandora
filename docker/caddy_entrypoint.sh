#!/bin/sh -e

cd /pyadmin && gunicorn admin:app --timeout 60 --graceful-timeout 20 --log-level DEBUG &

caddy -email stodyshev@gmail.com -agree --conf /etc/caddy/Caddyfile
