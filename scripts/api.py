import os
import requests
import json
import jwt
import re
import string
import random
import urllib
import termquery

dir = os.path.dirname(os.path.realpath(__file__))

VERBOSE = os.getenv('PYADMIN_VERBOSE', '') == '1'
DGRAPH_URL = os.getenv('DGRAPH_URL', 'http://dgraph:8080')
HTTP_PORT = os.getenv('HTTP_PORT', 80)
DEFAULT_SERVER_URL = 'http://localhost:{0}'.format(HTTP_PORT)
API_GATEWAY_URL = os.getenv('API_GATEWAY_URL', DEFAULT_SERVER_URL)
API_KEY = os.getenv('API_KEY')

print('VERBOSE: {0}'.format(VERBOSE))
print('DGRAPH_URL: {0}'.format(DGRAPH_URL))
print('API_GATEWAY_URL: {0}'.format(API_GATEWAY_URL))

jwt_secret = os.getenv('JWT_SECRET')
dgraph_token = os.getenv('DGRAPH_TOKEN')
system_token = jwt.encode(
    {
        'user_id': 'system',
        'user_name': 'system',
        'email': 'stodyshev@gmail.com',
        'role': 'admin',
    }, jwt_secret).decode('utf-8')

access_token = ''

MIME_JSON = 'application/json'

stdHeaders = {
    'Accept': MIME_JSON,
    'Content-Type': MIME_JSON,
    'Authorization': 'local_admin',
}


def is_json(resp):
    t = resp.headers.get('Content-Type')
    return t.startswith(MIME_JSON)


def dump_response(resp):
    if VERBOSE and is_json(resp):
        try:
            print(json.dumps(resp.json(), sort_keys=True, indent=2))
        except:
            print(resp.text)


def headers(content_type=MIME_JSON):
    h = stdHeaders.copy()
    t = access_token if access_token else system_token
    h['Authorization'] = 'Bearer ' + t
    if content_type is not None:
        h['Content-Type'] = content_type
    return h


def url(path):
    return API_GATEWAY_URL + path


def get(path):
    params = {'key': API_KEY}
    resp = requests.get(url(path), params=params, headers=headers(None))
    dump_response(resp)
    resp.raise_for_status()
    return resp.json() if is_json(resp) else resp


def jsonstr(data):
    return json.dumps(data, sort_keys=True, indent=2)


def post(path,
         payload,
         params={},
         auth=None,
         raw=False,
         content_type=MIME_JSON):
    params_all = {'key': API_KEY, **params}
    h = headers(content_type)
    data = payload if raw or content_type != MIME_JSON else jsonstr(payload)
    resp = requests.post(url(path),
                         data=data,
                         params=params_all,
                         headers=h,
                         auth=auth)
    dump_response(resp)
    resp.raise_for_status()
    return resp.json()


def put(path, payload, auth=None, raw=False):
    params = {'key': API_KEY}
    data = payload if raw else json.dumps(payload, sort_keys=True, indent=2)
    resp = requests.put(url(path),
                        data=data,
                        params=params,
                        headers=headers(),
                        auth=auth)
    dump_response(resp)
    resp.raise_for_status()
    return resp.json()


def delete(path, auth=None):
    params = {'key': API_KEY}
    resp = requests.delete(url(path),
                           params=params,
                           headers=headers(),
                           auth=auth)
    resp.raise_for_status()
    return resp


def login(username, password):
    global access_token
    resp = post('/api/login', None, auth=(username, password))
    access_token = resp['token']
    return access_token


def check_token(token):
    headers = {'Authorization': 'Bearer ' + token}
    resp = requests.get(url('/api/token'), headers=headers)
    dump_response(resp)
    return resp


def query(text, params={}):
    return post('/api/query', text, params=params, raw=True)


def current_user():
    return get('/api/me')


def search_terms(text, lang):
    q = termquery.make_term_query(searchString=text, lang=lang)
    return query(q['text'], q['params'])


# todo split text into words
def add_term(text, lang):
    resp = search_terms(text, lang)
    if len(resp['terms']) > 0:
        return resp['terms'][0]['uid']
    resp = post('/api/data/term', {'text': text, 'lang': lang})
    return resp['uid']


def link_terms(source_id, target_id, edge):
    rel = termquery.relation_map[edge]
    reverse_edge = edge
    if 'reverse_edge' in rel:
        reverse_edge = rel['reverse_edge']
    q = '\n'.join([
        nquad(source_id, edge, target_id),
        nquad(target_id, reverse_edge, source_id)
    ])
    return post('/api/nquads', {'set': q})


def delete_edge(id, edge):
    q = nquad(id, edge, '*')
    return post('/api/nquads', {'delete': q})


def search_audio(text, lang):
    txt = urllib.parse.quote(text)
    url = '/api/pyadmin/search/audio/{0}?lang={1}'.format(txt, lang)
    return get(url)


def fileproxy(url):
    def path_from_url():
        return re.sub(r'https?://', '', url)

    resp = get('/api/fileproxy/{0}'.format(url))
    host = os.getenv('SERVER_URL', 'http://lingvograph.com')
    path = resp['path'] if 'path' in resp else path_from_url()
    result = '{0}/api/file/{1}'.format(host, path)
    return result


def drop_all():
    headers = {'X-Dgraph-AuthToken': dgraph_token}
    data = '{"drop_all": true}'
    url = DGRAPH_URL + '/alter'
    resp = requests.post(url, headers=headers, data=data)
    dump_response(resp)
    resp.raise_for_status()


def schema_path():
    filename = os.path.realpath(os.path.join(dir, '../schema.txt'))
    return filename


def init_schema():
    p = schema_path()
    with open(p, 'r') as f:
        schema = f.read()
        headers = {'X-Dgraph-AuthToken': dgraph_token}
        url = DGRAPH_URL + '/alter'
        resp = requests.post(url, headers=headers, data=schema)
        dump_response(resp)
        resp.raise_for_status()


def mutate(data):
    headers = {
        'X-Dgraph-AuthToken': dgraph_token,
        'X-Dgraph-CommitNow': 'true',
    }
    if getattr(data, 'encode', None):
        data = data.encode('utf-8')
    url = DGRAPH_URL + '/mutate'
    resp = requests.post(url, headers=headers, data=data)
    dump_response(resp)
    resp.raise_for_status()
    return resp.json()


# utils
def is_uid(s):
    return len(s) > 0 and re.match(r"^0x[a-f0-9]+$", s) != None


def rdf_repr(v):
    if isinstance(v, str):
        if v == '*':
            return v
        return "<{0}>".format(v) if is_uid(v) else '"{0}"'.format(v)
    return v


def nquad(id, k, v):
    a = k.split('@')
    p = a[0]
    lang = a[1] if len(a) == 2 else ''
    s = rdf_repr(v)
    if len(lang) > 0:
        s += "@{0}".format(lang)
    id = "<{0}>".format(id) if is_uid(id) else "_:{0}".format(id)
    return "{0} <{1}> {2} .".format(id, p, s)


# TODO refactor as generator
def nquads(d, id='x'):
    result = []
    for k, v in d.items():
        result.append(nquad(id, k, v))
    return result


def set_nquads(dataset):
    data = "{\nset {\n" + dataset + "}}"
    print(data)
    return mutate(data)


def generate_string(n):
    return ''.join(random.SystemRandom().choice(string.ascii_uppercase +
                                                string.digits)
                   for _ in range(n))


def generate_api_key(app_id='', app_secret=''):
    if len(app_id) == 0:
        app_id = generate_string(8)
    if len(app_secret) == 0:
        app_secret = generate_string(16)
    key = os.getenv('API_KEY_SECRET')
    headers = {'app_secret': app_secret}
    payload = {'app_id': app_id}
    api_key = jwt.encode(payload, key + app_secret,
                         headers=headers).decode('utf-8')
    result = {'app_id': app_id, 'app_secret': app_secret, 'api_key': api_key}
    return result
