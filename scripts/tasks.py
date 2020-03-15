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

    # TODO keyword -> tag

    doc = meta
    doc['url'] = url
    doc['text'] = result['text']
    doc['author'] = {'uid': author_id}

    if id is None:
        return api.post('/api/data/document', doc)
    return api.put(f'/api/data/document/{id}', doc)


def search_doc(url):
    q = """query doc($url: string) {
  doc(func: eq(url, $url)) @filter(has(Document)) {
	uid
    expand(_all_)
  }
}
"""
    resp = api.query(q, params={'$url': url})
    if len(resp['doc']) == 1:
        return resp['doc'][0]['uid']
    return None


def search_person(name):
    q = """query person($name: string) {
  person(func: eq(name, $name)) @filter(has(Person)) {
	uid
    expand(_all_)
  }
}
"""
    resp = api.query(q, params={'$name': name})
    if len(resp['person']) == 1:
        return resp['person'][0]['uid']
    return None
