#!/usr/bin/env python

import api
import requests


def test_upload():
    api.login("admin", "admin123")
    headers = api.headers()

    file_url = api.url('/api/file/schema.txt')

    files = {'file': open(api.schema_path(), 'rb')}
    resp = requests.post(file_url, headers=headers, files=files)
    resp.raise_for_status()

    files = {'file': open(api.schema_path(), 'rb')}
    resp = requests.put(file_url, headers=headers, files=files)
    resp.raise_for_status()

    resp = requests.delete(file_url, headers=headers)
    resp.raise_for_status()
