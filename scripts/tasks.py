import os
import requests
from celery import Celery
from elasticsearch import Elasticsearch
from worker import app
import api

TIKA_HOST = os.getenv('TIKA_HOST', 'http://localhost:4219')


# just for testing purposes
@app.task
def add(x, y):
    return x + y


@app.task
def index_file(url):
    print(f'indexing {url}')
    # parse file by given URL
    resp = requests.get(url=TIKA_HOST + '/api/tika/parse', params={'url': url})
    resp.raise_for_status()
    result = resp.json()

    api.login("system", os.getenv("SYSTEM_PWD"))

    meta = result['metadata']

    id = search_doc(url)
    author = meta.get('author', '')
    if author == '':
        author = meta.get('creator', '')
    author_id = search_person(author)
    if author_id is None:
        person = {'name': author}
        person = api.post('/api/data/person', person)
        author_id = person['uid']

    tag = []
    keyword = meta.get('keyword', '')
    if isinstance(keyword, str):
        tag = make_tags(keyword.split(','))
    else:
        tag = make_tags(keyword)

    doc = meta
    doc['url'] = url
    doc['text'] = result['text']
    doc['author'] = {'uid': author_id}

    if len(tag) > 0:
        doc['tag'] = tag

    if id is None:
        return api.post('/api/data/document', doc)
    return api.put(f'/api/data/document/{id}', doc)


def make_tags(keywords):
    return [make_tag(k) for k in keywords]


def make_tag(text):
    id = find_tag(text)
    if id is not None:
        return {'uid': id}

    tag = {'text': text}
    tag = api.post('/api/data/term', tag)
    return {'uid': tag['uid']}


def find_tag(text):
    return find_node_by('text', text, 'Term')


def search_doc(url):
    return find_node_by('url', url, 'Document')


def search_person(name):
    return find_node_by('name', name, 'Person')


def find_node_by(predicate, value, resourceType):
    q = (
        f"query node($value: string) {{\n"
        f"node(func: eq({predicate}, $value)) @filter(has({resourceType})) {{\n"
        f"uid\n"
        f"expand(_all_)\n"
        f"}}\n"
        f"}}\n")

    resp = api.query(q, params={'$value': value})
    if len(resp['node']) > 0:
        return resp['node'][0]['uid']
    return None
