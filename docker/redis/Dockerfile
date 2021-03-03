FROM redis:5.0.12-alpine

RUN apk add --update curl \
    && rm -rf /var/cache/apk/*

# COPY redis.conf /usr/local/etc/redis/redis.conf
# CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
CMD [ "redis-server" ]
