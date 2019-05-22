#!/usr/bin/env python3

import os
import api
import resetdb as resetdb_impl
import audiosource

from langdetect import detect
from functools import wraps
from flask import Flask, request, jsonify
from werkzeug.contrib.fixers import ProxyFix

app = Flask(__name__)
app.debug = True
done = 'done!'


def auth(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        auth_ok = False
        auth = request.headers.get('Authorization', '').split(' ')
        if len(auth) == 2 and auth[0].lower() == 'bearer':
            resp = api.check_token(auth[1])
            auth_ok = resp.ok
        if not auth_ok:
            secret = request.args.get('secret')
            if secret != os.getenv('ADMIN_SECRET'):
                return 'bad auth', 401
        return f(*args, **kwargs)

    return decorated_function


@app.route('/api/pyadmin/dropall')
@auth
def dropall():
    api.drop_all()
    api.init_schema()
    return done


@app.route('/api/pyadmin/initschema')
@auth
def initschema():
    api.drop_all()
    api.init_schema()
    return done


@app.route('/api/pyadmin/resetdb')
@auth
def resetdb():
    resetdb_impl.run()
    return done


@app.route('/api/pyadmin/search/audio/<text>')
@auth
def find_audio(text):
    lang = request.args.get('lang')
    if lang is None or lang == '':
        lang = detect(text)
        if lang != 'ru':
            lang = 'en'
    result = audiosource.find_audio(text, lang)
    return jsonify(result)


app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
