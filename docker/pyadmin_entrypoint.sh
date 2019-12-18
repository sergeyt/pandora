#!/bin/sh -e

cd /pyadmin

FLASK_PORT=${FLASK_PORT:-4211}
FLASK_ENV=${FLASK_ENV:-development}
FLASK_WORKERS=${FLASK_WORKERS:-2}

export FLASK_ENV=${FLASK_ENV}
export FLASK_APP=admin.py

pipenv install --deploy --system

if [[ "$FLASK_ENV" = "development" ]]; then
    python -m flask run --host 0.0.0.0 --port ${FLASK_PORT}
else
    gunicorn --bind :${FLASK_PORT} --workers ${FLASK_WORKERS} --timeout 900 --graceful-timeout 30 --log-level DEBUG admin:app
fi
