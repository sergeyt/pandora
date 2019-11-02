#!/usr/bin/env python

import sys
import os
import cambridge
import api
from models import Term, File

reverse_edges = {
    'definition': 'definition_of',
}

__dir__ = os.path.dirname(os.path.realpath(__file__))
TERMS = {}  # key=text@lang, value=uid
__next_id__ = 1


def next_id(prefix):
    global __next_id__
    id = prefix + str(__next_id__)
    __next_id__ += 1
    return id


# TODO refactor as generator
def read_words():
    with open(os.path.join(__dir__, 'us1000.txt'), 'r', encoding='utf-8') as f:
        lines = f.read().split('\n')
        lines = [s.strip() for s in lines]
        return [s for s in lines if len(s) > 0]


def define_term(data):
    print('TERM {0}'.format(data.text))
    key = '{0}@{1}'.format(data.text, data.lang)
    if key in TERMS:
        return TERMS[key]
    id = api.add_term(data.text, data.lang, data.region)
    key = '{0}@{1}'.format(data.text, data.lang)
    TERMS[key] = id
    return id


def define_word(word, lang='en'):
    word_id = define_term(Term(text=word, lang=lang, region=None))
    data = cambridge.get_data(word, lang)
    for k, a in data.items():
        edges = []
        for v in a:
            is_file = isinstance(v, File)
            if is_file:
                file = api.fileproxy(v.url, as_is=True)
                related_id = file['uid']
                if v.region:
                    edges.append([related_id, 'region', v.region])
            else:
                related_id = define_term(v)
            edges.append([word_id, k, related_id])
            if k in reverse_edges:
                edges.append([related_id, reverse_edges[k], word_id])
        if len(edges) > 0:
            api.update_graph(edges)


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    if len(sys.argv) > 1:
        define_word(sys.argv[1])
        return
    words = read_words()
    for word in words[:10]:
        define_word(word)


if __name__ == '__main__':
    main()
