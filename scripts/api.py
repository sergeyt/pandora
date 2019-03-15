import os
import requests
import json
import jwt
from dotenv import load_dotenv

dir = os.path.dirname(os.path.realpath(__file__))
dotenv_path = os.path.realpath(os.path.join(dir, '../.env'))
load_dotenv(dotenv_path=dotenv_path)

DGRAPH_URL = os.getenv('DGRAPH_URL', 'http://dgraph:8080')
API_GATEWAY_URL = os.getenv('API_GATEWAY_URL', 'http://localhost:4200')

jwt_secret = os.getenv('JWT_SECRET')
dgraph_token = os.getenv('DGRAPH_TOKEN')
system_token = jwt.encode({
    'user_id': 'system',
    'user_name': 'system',
    'email': 'stodyshev@gmail.com',
    'role': 'admin',
}, jwt_secret).decode('utf-8')

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


def url(path):
    return API_GATEWAY_URL + path


def get(path):
    resp = requests.get(url(path), headers=headers())
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()


def post(path, payload, auth=None, raw=False):
    data = payload if raw else json.dumps(payload, sort_keys=True, indent=2)
    resp = requests.post(url(path), data=data, headers=headers(), auth=auth)
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()


def put(path, payload, auth=None, raw=False):
    data = payload if raw else json.dumps(payload, sort_keys=True, indent=2)
    resp = requests.put(url(path), data=data, headers=headers(), auth=auth)
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()


def delete(path, auth=None):
    resp = requests.delete(url(path), headers=headers(), auth=auth)
    resp.raise_for_status()
    return resp


def login(username, password):
    resp = post('/api/login', None, auth=(username, password))
    access_token = resp['token']
    return access_token


def drop_all():
    headers = {'X-Dgraph-AuthToken': dgraph_token}
    resp = requests.post(
        DGRAPH_URL + '/alter', headers=headers, data='{"drop_all": true}')
    dump_json(resp)
    resp.raise_for_status()


def schema_path():
    filename = os.path.realpath(os.path.join(dir, '../schema.txt'))
    return filename


def init_schema():
    p = schema_path()
    with open(p, 'r') as f:
        schema = f.read()
        headers = {'X-Dgraph-AuthToken': dgraph_token}
        resp = requests.post(
            DGRAPH_URL + '/alter', headers=headers, data=schema)
        dump_json(resp)
        resp.raise_for_status()


def mutate(data):
    headers = {
        'X-Dgraph-AuthToken': dgraph_token,
        'X-Dgraph-CommitNow': 'true',
    }
    if getattr(data, 'encode', None):
        data = data.encode('utf-8')
    resp = requests.post(DGRAPH_URL + '/mutate', headers=headers, data=data)
    dump_json(resp)
    resp.raise_for_status()
    return resp.json()


def rdf_repr(v):
    if isinstance(v, str):
        return '"{0}"'.format(v)
    return v


def nquad(id, k, v):
    a = k.split('@')
    p = a[0]
    lang = a[1] if len(a) == 2 else ''
    s = rdf_repr(v)
    if len(lang) > 0:
        s += "@{0}".format(lang)
    return "_:{0} <{1}> {2} .\n".format(id, p, s)


def nquads(d, id='x'):
    result = ''
    for k, v in d.items():
        result += nquad(id, k, v)
    return result


def set_nquads(dataset):
    data = "{\nset {\n" + dataset + "}}"
    print(data)
    return mutate(data)
