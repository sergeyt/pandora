import requests
import json

host = 'http://localhost:4200'
access_token = ''

stdHeaders = {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    'Authorization': 'local_admin',
}

def dump_json(resp):
    if resp.ok:
        print(json.dumps(resp.json(), sort_keys=True, indent=2))

def headers():
    h = stdHeaders.copy()
    if access_token:
        h['Authorization'] = 'Bearer ' + access_token
    return h

def get(path, data):
    url = host + path
    resp = requests.get(url, headers=headers())
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()

def post(path, payload, auth=None):
    url = host + path
    data = json.dumps(payload, sort_keys=True, indent=2)
    resp = requests.post(url, data=data, headers=headers(), auth=auth)
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()

def login(username, password):
    resp = post('/api/login', None, auth=(username, password))
    access_token = resp['token']
    return access_token

def drop_all():
    resp = requests.post('http://localhost:8080/alter', data='{"drop_all": true}')
    dump_json(resp)
    resp.raise_for_status()

def init_schema():
    with open('schema.txt', 'r') as f:
        schema = f.read()
        resp = requests.post('http://localhost:8080/alter', data=schema)
        dump_json(resp)
        resp.raise_for_status()
