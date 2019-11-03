#!/usr/bin/env python

import os
import api


def test_search_audio():
    api.login("system", os.getenv("SYSTEM_PWD"))
    api.search_audio('apple', 'en')


def test_create_term():
    api.login("system", os.getenv("SYSTEM_PWD"))
    api.post('/api/pyadmin/term', {'text': 'girl', 'lang': 'en'})
    api.post('/api/pyadmin/term', {'text': 'girl', 'lang': 'en'})
