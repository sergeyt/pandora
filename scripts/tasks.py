import os
import requests
from celery import Celery
from elasticsearch import Elasticsearch
from worker import app

TIKA_HOST = os.getenv('TIKA_HOST', 'http://localhost:4219')
ES_HOSTS = os.getenv('ES_HOSTS', 'localhost:9200')
ES_INDEX = os.getenv('ES_INDEX_DOCS', 'docs')

es = Elasticsearch(hosts=ES_HOSTS)

# ensure ES index
if not es.indices.exists(ES_INDEX):
    # TODO configure mappings
    es.indices.create(ES_INDEX)


# just for testing purposes
@app.task
def add(x, y):
    return x + y


# TODO update existing doc found by url
@app.task
def index_doc(url):
    # parse file by given URL
    resp = requests.get(url=TIKA_HOST + '/api/tika/parse', params={'url': url})
    resp.raise_for_status()
    result = resp.json()

    # push metadata and text to elasticsearch
    doc = result['metadata']
    doc['source_url'] = url
    doc['text'] = result['text']
    es.index(ES_INDEX, body=doc)

    # TODO push dgraph data
    return result
