import os
import requests
import json
import jwt
from dotenv import load_dotenv

load_dotenv(dotenv_path='.env')

host = 'http://localhost:4200'

jwt_secret = os.getenv('JWT_SECRET')
dgraph_token = os.getenv('DGRAPH_TOKEN')
system_token = jwt.encode({
    'user_id': 'system',
    'user_name': 'system',
    'email': 'stodyshev@gmail.com',
    'role': 'admin',
}, jwt_secret)

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
    t = access_token if access_token else system_token
    h['Authorization'] = 'Bearer ' + t
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
    headers = {'X-Dgraph-AuthToken': dgraph_token}
    resp = requests.post('http://localhost:8080/alter', headers=headers, data='{"drop_all": true}')
    dump_json(resp)
    resp.raise_for_status()

def init_schema():
    with open('schema.txt', 'r') as f:
        schema = f.read()
        headers = {'X-Dgraph-AuthToken': dgraph_token}
        resp = requests.post('http://localhost:8080/alter', headers=headers, data=schema)
        dump_json(resp)
        resp.raise_for_status()

def mutate(data):
    resp = requests.post('http://localhost:8080/mutate', data=data)
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()

def set_nquads(dataset):
    data = "{\nset{\n" + dataset + "}}"
    return mutate(data)
