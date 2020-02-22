FROM python:3.8-alpine

RUN apk update \
    && pip3 install pipenv \
    && pip3 install gunicorn

COPY scripts /pyadmin
COPY schema.txt /schema.txt
COPY docker/pyadmin/main.sh /main.sh

WORKDIR /pyadmin

ENTRYPOINT ["/main.sh"]
