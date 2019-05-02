#!/bin/sh -e

cd /pyadmin && gunicorn --timeout 900 --graceful-timeout 30 --log-level DEBUG admin:app &

caddy -email stodyshev@gmail.com -agree --conf /etc/caddy/Caddyfile
