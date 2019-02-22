#!/usr/bin/env python3

import sys
import os
import logging
import traceback
import api
import initdata

from flask import Flask
from werkzeug.contrib.fixers import ProxyFix

app = Flask(__name__)


@app.route('/api/pyadmin/resetdb')
def resetdb():
    api.drop_all()
    api.init_schema()
    initdata.init()
    return 'done!'


app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
