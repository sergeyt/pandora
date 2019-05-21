#!/usr/bin/env python3

import os
import audiosource
from langdetect import detect

import api
import resetdb as resetdb_impl

from functools import wraps
from flask import Flask, request
from werkzeug.contrib.fixers import ProxyFix

app = Flask(__name__)
app.debug = True
done = 'done!'


def auth(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        secret = request.args.get('secret')
        if secret is None or secret == '' or secret != os.getenv(
                'ADMIN_SECRET'):
            return 'bad secret!'
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
    return audiosource.find_audio(text, lang)


app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
