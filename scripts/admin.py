#!/usr/bin/env python3

import sys
import os
import logging
import traceback
import api
import initdata
import words

from functools import wraps
from flask import Flask, request
from werkzeug.contrib.fixers import ProxyFix

app = Flask(__name__)
app.debug = True
done = 'done!'


def secret_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        secret = request.args.get('secret')
        if secret is None or secret == '' or secret != os.getenv(
                'ADMIN_SECRET'):
            return 'bad secret!'
        return f(*args, **kwargs)

    return decorated_function


@app.route('/api/pyadmin/dropall')
@secret_required
def dropall():
    api.drop_all()
    api.init_schema()
    return done


@app.route('/api/pyadmin/initschema')
@secret_required
def initschema():
    api.drop_all()
    api.init_schema()
    return done


@app.route('/api/pyadmin/resetdb')
@secret_required
def resetdb():
    api.drop_all()
    api.init_schema()
    initdata.init()
    words.init()
    return done


app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
