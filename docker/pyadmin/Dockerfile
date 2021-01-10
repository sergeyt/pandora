# see also https://github.com/wemake-services/wemake-django-template/blob/master/%7B%7Bcookiecutter.project_name%7D%7D/docker/django/Dockerfile
FROM python:3.9-alpine

RUN apk update \
    && apk --no-cache add \
    bash \ 
    build-base \
    tini

WORKDIR /pyadmin

COPY scripts /pyadmin
COPY schema.txt /schema.txt
COPY docker/pyadmin/main.sh /main.sh

RUN chmod +x "/main.sh" \
    && rm -f requirements.txt \
    && python requirements.py > requirements.txt \
    && pip install -r requirements.txt

ENTRYPOINT ["/sbin/tini", "--", "/main.sh"]
