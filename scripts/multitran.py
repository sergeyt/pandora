#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import utils
import requests
import json
from bs4 import BeautifulSoup

headers = {
    'User-Agent': 'script',
    'Accept': 'text/html',
}

# TODO consider collecting automatically https://www.multitran.com/m.exe?s=place&l1=2&l2=1&fl=1
categories = [{
    'text@en': 'idiom',
    'text@ru': 'идиома',
    'id': 895,
}, {
    'text@en': 'proverb',
    'text@ru': 'пословица',
    'id': 310,
}, {
    'text@en': 'americanism',
    'text@ru': 'американизм',
    'id': 13,
}, {
    'text@en': 'bible',
    'text@ru': 'библия',
    'id': 66,
}]


def parse_phrase_row(row):
    def parse_td(td):
        # if 'class' not in td.attrs:
        #     return None
        # k = td.attrs['class']
        # if k not in ['phraselist1', 'phraselist2']:
        #     return None
        a = td.find('a')
        if a is None:
            return a
        return a.get_text().strip()

    result = [parse_td(t) for t in row.find_all('td')]
    if len(result) != 2:
        return None
    if any(t is None for t in result):
        return None
    return {'text': result[0], 'trans': result[1]},


def find_phrases_impl(text, lang, category):
    l1 = "1"
    l2 = "2"
    if lang == 'ru':
        l1 = "2"
        l2 = "1"
    pat = "https://www.multitran.com/m.exe?a=3&sc=895&s={0}&l1={1}&l2={2}"
    url = pat.format(*utils.url_quote([text, l1, l2]))

    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    rows = [parse_phrase_row(r) for r in soup.find_all('tr')]

    return {
        'text@en': category['text@en'],
        'text@ru': category['text@ru'],
        'multitran_id': category['id'],
        'phrases': [r for r in rows if r is not None],
    }


def find_phrases(text, lang):
    return {
        'text': text,
        'lang': lang,
        'categories': [find_phrases_impl(text, lang, c) for c in categories],
    }


def main():
    (text, lang) = utils.find_audio_args()
    result = find_phrases(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
