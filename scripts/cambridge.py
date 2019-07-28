#!/usr/bin/env python

import sys
import requests
import json
import utils
from bs4 import BeautifulSoup
from localcache import Cache

cache = Cache('cambridge')

headers = {
    'User-Agent': 'script',
    'Accept': 'text/html',
}


def parse_btn(btn):
    mp3 = btn['data-src-mp3']
    ogg = btn['data-src-ogg']
    region = btn.parent['class'][0]
    return {'mp3': mp3, 'ogg': ogg, 'region': region}


def find_audio(text, lang):
    if lang != 'en':
        return None

    result = cache.get(text)
    if result is not None:
        return result

    pat = 'https://dictionary.cambridge.org/dictionary/english/{0}'
    url = pat.format(text)

    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    header = soup.find('div', class_='pos-header')
    if header is None:
        return None

    btns = header.find_all('span', class_='audio_play_button')
    data = [parse_btn(b) for b in btns]
    result = {}
    for d in data:
        for fmt in ['mp3', 'ogg']:
            if fmt in d:
                if fmt not in result:
                    result[fmt] = []
                url = 'https://dictionary.cambridge.org' + d[fmt]
                if utils.url_exists(url):
                    result[fmt].append({'url': url, 'region': d['region']})

    cache.put(text, result)

    return result


def stripped_text(node):
    return node.get_text().strip()


def parse_def(d):
    text = stripped_text(d.find('b', class_='def'))
    trans = d.find('span', class_='trans')
    examples = [stripped_text(t) for t in d.find_all('span', class_='eg')]
    return {
        'text': text,
        'trans': {
            'text': stripped_text(trans),
            'lang': trans.attrs['lang'],
        },
        'examples': examples,
    }


def parse_sense(block, text):
    phrase = text
    d = block.find('div', class_='phrase-block')
    if d is not None:
        phrase = stripped_text(d.find('span', class_='phrase'))
    defs = [parse_def(t) for t in block.find_all('div', class_='def-block')]
    return {
        'phrase': phrase,
        'defs': defs,
    }


def translate(text, lang):
    if lang != 'en':
        return None

    pat = 'https://dictionary.cambridge.org/ru/{0}/{1}/{2}'
    url = pat.format(*utils.url_quote(['словарь', 'англо-русский', text]))

    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    ipa = soup.find('span', class_='ipa').string
    blocks = soup.find_all('div', class_='sense-block')
    phrases = [parse_sense(t, text) for t in blocks]
    return {
        'text': text,
        'lang': lang,
        'ipa': ipa,
        'phrases': phrases,
    }


def main():
    cmd = sys.argv[1]
    (text, lang) = utils.find_audio_args(sys.argv[1:])
    if cmd == 'audio':
        result = find_audio(text, lang)
    else:
        result = translate(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
