#!/usr/bin/env python

import os
import api
import requests


def test_upload():
    api.login("system", os.getenv("SYSTEM_PWD"))
    headers = api.headers()
    params = {'key': api.API_KEY}

    file_url = api.url('/api/file/schema.txt')

    files = {'file': open(api.schema_path(), 'rb')}
    resp = requests.post(file_url, params=params, headers=headers, files=files)
    resp.raise_for_status()

    # node = resp.json()
    # file_id = node['uid']

    # files = {'file': open(api.schema_path(), 'rb')}
    # resp = requests.put(file_url, params=params, headers=headers, files=files)
    # resp.raise_for_status()
