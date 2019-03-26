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

    print('CREATE')

    resp = api.post('/api/data/user', data)

    id = resp['uid']

    nquads = '\n'.join(['<{0}> <first_lang> "ru" .'.format(id)])

    api.post('/api/nquads', nquads, content_type='application/n-quads')

    resp = api.get('/api/data/user/list')

    print('DELETE')

    api.delete('/api/data/user/' + id)
