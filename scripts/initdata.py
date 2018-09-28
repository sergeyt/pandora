import requests
import json
import utils

utils.enable_logging()

base = 'http://localhost:4200'
access_token = ''

stdHeaders = {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    'Authorization': 'local_admin',
}

def headers():
    h = stdHeaders.copy()
    if access_token:
        h['Authorization'] = 'Bearer ' + access_token
    return h

def dump_json(resp):
    if resp.ok:
        print(json.dumps(resp.json(), sort_keys=True, indent=2))

def get(path, data):
    url = base + path
    resp = requests.get(url, headers=headers())
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()

def post(path, payload, auth=None):
    url = base + path
    data = json.dumps(payload, sort_keys=True, indent=2)
    resp = requests.post(url, data=data, headers=headers(), auth=auth)
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()

def login(username, password):
    resp = post('/api/login', None, auth=(username, password))
    access_token = resp['token']
    return access_token

def user_exists(user):
    try:
        login(user['login'], user['password'])
    except:
        return False
    return True

def ensure_user(user):
    if not user_exists(user):
        post('/api/data/user', user)

users = [
    {
        'login': 'admin',
        'name': 'admin',
        'email': 'stodyshev@gmail.com',
        'password': 'admin123',
    },
    {
        'login': 'sergeyt',
        'name': 'sergeyt',
        'email': 'stodyshev@gmail.com',
        'password': 'sergeyt123',
    },
]

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

def init():
    for user in users:
        ensure_user(user)

# drop_all()
# init_schema()
init()
