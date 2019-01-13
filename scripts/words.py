#!/usr/bin/env python

import json
import utils
import api

with open('scripts/words.json', 'r', encoding='utf-8') as f:
    words = json.loads(f.read())


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


def get_id(resp, id='x'):
    return resp['data']['uids'][id]


def pairs(list):
    for x in list:
        for y in list:
            yield (x, y)


# TODO find existing words and update their transcriptions, pronunciations, images
# TODO port this script to golang and provide an API endpoint
for word in words:
    # insert word nodes for all languages
    nodes = {}
    word_ids = []
    for key in word:
        if key.startswith('text@'):
            text = word[key]
            lang = key[key.index('@')+1:]
            props = {
                '_word': '',
                'text': text,
                'lang': lang,
            }
            trans_key = 'transcription@' + lang
            if trans_key in word:
                props['transcription'] = word[trans_key]
            data = nquads(props)
            resp = api.set_nquads(data)
            wid = get_id(resp)
            nodes[lang] = wid
            word_ids.append(wid)

    # link words together with translated_as predicate
    proc = {}
    for (w1, w2) in pairs(word_ids):
        if w1 == w2:
            continue
        k1 = "{0}-{1}".format(w1, w2)
        k2 = "{0}-{1}".format(w2, w1)
        if k1 in proc:
            continue
        data = "<{0}> <translated_as> <{1}> .\n".format(w1, w2)
        api.set_nquads(data)
        proc[k1] = 1
        proc[k2] = 1

    # todo upload images to S3 store
    # insert images
    image_nodes = []
    for url in word['images']:
        data = nquads({
            '_image': '',
            'url': url,
            'source': 'google',
        })
        resp = api.set_nquads(data)
        obj_id = get_id(resp)
        # link with connected words
        for lang in nodes:
            wid = nodes[lang]
            data = "<{0}> <relevant> <{1}> .\n".format(wid, obj_id)
            api.set_nquads(data)

    # todo upload sounds to S3 store
    for key in word:
        if key.startswith('pronunciation@'):
            lang = key[key.index('@')+1:]
            data = nquads({
                '_sound': '',
                'url': word[key],
                'source': 'https://howjsay.com',
            })
            resp = api.set_nquads(data)
            obj_id = get_id(resp)
            # link with word
            wid = nodes[lang]
            data = "<{0}> <pronounced_as> <{1}> .\n".format(wid, obj_id)
            api.set_nquads(data)
