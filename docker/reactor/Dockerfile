FROM python:3.8-alpine
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8 PYTHONUNBUFFERED=1

RUN apk add --no-cache bash inotify-tools

COPY ./scripts /app
WORKDIR /app
RUN pip install --no-cache-dir -r requirements.txt

CMD [ "bash", "/app/reload.sh", "python", "reactor.py" ]
