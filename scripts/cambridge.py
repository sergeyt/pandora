#!/usr/bin/env python

import requests
import json
import utils
from bs4 import BeautifulSoup
from localcache import Cache

cache = Cache('cambridge')


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
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
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
                result[fmt].append({'url': url, 'region': d['region']})

    cache.put(text, result)

    return result


def main():
    (text, lang) = utils.find_audio_args()
    result = find_audio(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
