#!/usr/bin/env python

import os
import api

# todo check bad auth cases as separate test cases


def test_crud():
    api.login("system", os.getenv("SYSTEM_PWD"))

    data = {
        'name': "bob",
        'age': 39,
    }

    print('CREATE')

    resp = api.post('/api/data/user', data)

    id = resp['uid']

    data = {
        'name': "joe",
        'age': 40,
    }

    resp = api.post('/api/data/user', data)
    id2 = resp['uid']

    print('GET LIST')

    resp = api.get('/api/data/user/list')

    print('GET BY ID')

    resp = api.get('/api/data/user/' + id)

    print('QUERY')

    query = """{
  data(func: eq(name, "bob")) @filter(has(User)) {
    uid
    name
    age
  }
}"""
    resp = api.post('/api/query', query, raw=True)

    print('search terms')

    resp = api.search_terms('abc', 'en', no_links=True)

    print('UPDATE')

    data = {
        'name': 'rob',
        'age': 42,
    }

    resp = api.put("/api/data/user/" + id, data)

    print('GET BY ID')

    resp = api.get('/api/data/user/' + id)

    print('DELETE')

    api.delete('/api/data/user/' + id)
    api.delete('/api/data/user/' + id2)


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
