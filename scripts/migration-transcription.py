#!/usr/bin/env python

import api
import json


def dump_json(data):
    print(json.dumps(data, sort_keys=True, indent='  '))


def main():
    data = api.query("""{
  terms(func: has(Term)) @filter(has(transcript@en)) {
    uid
    text
    lang
    transcript@en
    transcript@ru
  }
}""")
    for t in data['terms']:
        key_en = 'transcript@en'
        key_ru = 'transcript@ru'
        for k in [key_en, key_ru]:
            if k in t:
                tlang = 'ipa' if k == 'en' else 'ru'
                id = api.add_term(t[key_en], tlang)
                api.link_terms(t['uid'], id, 'transcription')
        api.delete_edge(t['uid'], 'transcript')


if __name__ == '__main__':
    main()
