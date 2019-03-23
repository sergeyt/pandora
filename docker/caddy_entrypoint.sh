#!/bin/sh -e

cd /pyadmin && gunicorn --timeout 300 --graceful-timeout 30 --log-level DEBUG admin:app &

caddy -email stodyshev@gmail.com -agree --conf /etc/caddy/Caddyfile
