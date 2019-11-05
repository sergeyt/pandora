#!/usr/bin/env python

import sys
import os
import cambridge
import unsplash
import multitran
import merriamwebster
import api
from models import Term, File, TermWithData
from multiprocessing import Process

# here you can temporarily remove sources that you don't need to test
sources = [cambridge, merriamwebster, unsplash, multitran]

reverse_edges = {
    'transcription': 'transcription_of',
    'definition': 'definition_of',
    'collocation': 'collocation_of',
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


def key_of(text, lang):
    return '{0}@{1}'.format(text, lang)


def define_term(data):
    print('TERM {0}'.format(data.text))
    key = key_of(data.text, data.lang)
    if key in TERMS:
        return TERMS[key]
    id = api.add_term(data.text, data.lang, data.region)
    TERMS[key] = id
    return id


def push_data(term_id, data):
    for k, a in data.items():
        edges = []
        for v in a:
            is_file = isinstance(v, File)
            if is_file:
                # TODO optimize adding file with region
                file = api.fileproxy(v.url, as_is=True)
                related_id = file['uid']
                if v.region:
                    edges.append([related_id, 'region', v.region])
            elif isinstance(v, TermWithData):
                related_id = define_term(v.term)
                push_data(related_id, v.data)
            else:
                related_id = define_term(v)
            edges.append([term_id, k, related_id])
            if not is_file:
                reverse_edge = reverse_edges[k] if k in reverse_edges else k
                edges.append([related_id, reverse_edge, term_id])
        if len(edges) > 0:
            api.update_graph(edges)


def define_word(text, lang='en', source_idx=-1):
    term_id = define_term(Term(text=text, lang=lang, region=None))
    source_list = sources if source_idx < 0 else [sources[source_idx]]
    for source in source_list:
        data = source.get_data(text, lang)
        push_data(term_id, data)


def define_words(source_idx):
    words = read_words()
    for word in words:
        define_word(word, source_idx=source_idx)


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    for i in range(len(sources)):
        p = Process(target=define_words, args=(i, ))
        p.start()
        p.join()


if __name__ == '__main__':
    main()
