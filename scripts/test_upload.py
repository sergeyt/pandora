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

    node = resp.json()
    file_id = node['uid']

    files = {'file': open(api.schema_path(), 'rb')}
    resp = requests.put(file_url, headers=headers, files=files)
    resp.raise_for_status()

    # download by path
    resp = requests.get(file_url, headers=headers)
    resp.raise_for_status()
    print(resp.text)

    # download by id
    url1 = api.url('/api/file/{0}'.format(file_id))
    resp = requests.get(url1, headers=headers)
    resp.raise_for_status()
    print(resp.text)

    resp = requests.delete(file_url, headers=headers)
    resp.raise_for_status()


def test_delete_file_node():
    api.login("admin", "admin123")
    headers = api.headers()

    file_url = api.url('/api/file/schema2.txt')

    files = {'file': open(api.schema_path(), 'rb')}
    resp = requests.post(file_url, headers=headers, files=files)
    resp.raise_for_status()

    node = resp.json()
    file_id = node['uid']

    # download by path
    resp = requests.get(file_url, headers=headers)
    resp.raise_for_status()
    print(resp.text)

    # download by id
    url1 = api.url('/api/file/{0}'.format(file_id))
    resp = requests.get(url1, headers=headers)
    resp.raise_for_status()
    print(resp.text)

    url2 = api.url('/api/data/file/{0}'.format(file_id))
    resp = requests.delete(url2, headers=headers)
    resp.raise_for_status()
