#!/usr/bin/env python

import os
import api

# todo check bad auth cases as separate test cases


def crud(resource):
    api.login("system", os.getenv("SYSTEM_PWD"))

    data = {
        'name': "bob",
        'age': 39,
    }

    print('CREATE')

    base = '/api/data/' + resource
    resp = api.post(base, data)

    id = resp['uid']

    data = {
        'name': "joe",
        'age': 40,
    }

    resp = api.post(base, data)
    id2 = resp['uid']

    print('GET LIST')

    resp = api.get(f'{base}/list')

    print('GET BY ID')

    resp = api.get(f'{base}/{id}')

    print('QUERY')

    query = """{{
  data(func: eq(name, "bob")) @filter(has({0})) {{
    uid
    name
    age
  }}
}}""".format(resource.capitalize())
    resp = api.post('/api/query', query, raw=True)

    print('search terms')

    resp = api.search_terms('abc', 'en', no_links=True)

    print('UPDATE')

    data = {
        'name': 'rob',
        'age': 42,
    }

    resp = api.put(f"{base}/{id}", data)

    print('GET BY ID')

    resp = api.get(f'{base}/{id}')

    print('DELETE')

    api.delete(f'{base}/{id}')
    api.delete(f'{base}/{id2}')


def test_crud_user():
    crud('user')


def test_crud_term():
    crud('term')


def test_crud_document():
    crud('document')


def test_graph_update():
    api.login("system", os.getenv("SYSTEM_PWD"))

    data = {
        'name': "bob",
        'age': 39,
    }

    resp = api.post('/api/data/user', data)

    id = resp['uid']
    user_url = f'/api/data/user/{id}'

    nquads = '\n'.join([f'<{id}> <first_lang> "ru" .'])
    api.post('/api/nquads', nquads, content_type='application/n-quads')

    resp = api.get(user_url)
    assert resp['first_lang'] == 'ru'

    mutation = {
        'set': '\n'.join([f'<{id}> <age> "38"^^<xs:int> .']),
        'delete': '\n'.join([f'<{id}> <first_lang> * .']),
    }
    api.post('/api/nquads', mutation)

    resp = api.get(user_url)
    assert resp['age'] == 38
    assert 'first_lang' not in resp

    api.delete(user_url)
