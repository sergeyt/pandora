#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import utils
import requests
import json
from bs4 import BeautifulSoup
from models import Term, TermWithData
from utils import is_empty

NAME = 'multitran'

headers = {
    'User-Agent': utils.CHROME_USER_AGENT,
    'Accept': 'text/html',
}

# TODO consider collecting automatically https://www.multitran.com/m.exe?s=place&l1=2&l2=1&fl=1
categories = [{
    'tag': [
        Term(text='idiom', lang='en', region=None),
        Term(text='идиома', lang='ru', region=None),
    ],
    'id':
    895,
}, {
    'tag': [
        Term(text='proverb', lang='en', region=None),
        Term(text='пословица', lang='ru', region=None),
    ],
    'id':
    310,
}, {
    'tag': [
        Term(text='americanism', lang='en', region=None),
        Term(text='американизм', lang='ru', region=None),
    ],
    'id':
    13,
}, {
    'tag': [
        Term(text='bible', lang='en', region=None),
        Term(text='библия', lang='ru', region=None),
    ],
    'id':
    66,
}]


def stripped_text(node):
    if node is None:
        return None
    return node.get_text().strip()


def parse_phrase_row(row, lang, trans_lang, tags):
    def parse_td(td):
        # if 'class' not in td.attrs:
        #     return None
        # k = td.attrs['class']
        # if k not in ['phraselist1', 'phraselist2']:
        #     return None
        a = td.find('a')
        return stripped_text(a)

    result = [parse_td(t) for t in row.find_all('td')]
    if len(result) != 2:
        return []
    if any(is_empty(t) for t in result):
        return []
    term = Term(text=result[0], lang=lang)
    trans = Term(text=result[1], lang=trans_lang)
    return [TermWithData(term, {'tag': tags, 'translated_as': [trans]})]


def find_phrases(data, text, lang, category):
    trans_lang = 'ru'
    l1 = "1"
    l2 = "2"
    if lang == 'ru':
        l1 = "2"
        l2 = "1"
        trans_lang = 'en'
    pat = "https://www.multitran.com/m.exe?a=3&sc=895&s={0}&l1={1}&l2={2}"
    url = pat.format(*utils.url_quote([text, l1, l2]))

    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    for tr in soup.find_all('tr'):
        for term in parse_phrase_row(tr, lang, trans_lang, category['tag']):
            data['in'].append(term)


def get_data(text, lang):
    data = {'in': []}

    for c in categories:
        find_phrases(data, text, lang, c)

    return data


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
