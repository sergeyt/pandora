#!/usr/bin/env python

import os
import api
import requests
from os.path import isfile, join

BOOK_DIR = '/Users/admin/Dropbox/books'


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    headers = api.headers()
    params = {'key': api.API_KEY}

    files = [f for f in os.listdir(BOOK_DIR) if isfile(join(BOOK_DIR, f))]
    for f in files:
        url = api.url(f'/api/file/{f}')
        fullpath = join(BOOK_DIR, f)
        files = {'file': open(fullpath, 'rb')}
        resp = requests.post(url, params=params, headers=headers, files=files)
        resp.raise_for_status()


if __name__ == '__main__':
    main()
