#!/bin/sh -e

cd /pyadmin && gunicorn admin:app &

caddy -email stodyshev@gmail.com -agree --conf /etc/caddy/Caddyfile
