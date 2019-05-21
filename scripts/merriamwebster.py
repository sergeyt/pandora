#!/usr/bin/env python

import requests
import json
import utils
from bs4 import BeautifulSoup
from localcache import Cache

cache = Cache('merriamwebster')


def parse_btn(btn):
    lang = btn['data-lang']
    dir = btn['data-dir']
    file = btn['data-file']
    pat = 'https://media.merriam-webster.com/audio/prons/{0}/{1}/{2}.mp3'
    return pat.format(lang.replace('_', '/'), dir, file)


def find_audio(text, lang):
    if lang != 'en':
        return None

    result = cache.get(text)
    if result is not None:
        return result

    pat = 'https://www.merriam-webster.com/dictionary/{0}'
    url = pat.format(text)
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    prs = soup.find('span', class_='prs')
    if prs is None:
        return None
    btns = prs.find_all('a', class_='play-pron')
    if len(btns) == 0:
        return None

    urls = [parse_btn(b) for b in btns]
    result = {}
    for url in urls:
        fmt = 'mp3'
        if fmt not in result:
            result[fmt] = []
        result[fmt].append({'url': url})

    cache.put(text, result)

    return result


def main():
    (text, lang) = utils.find_audio_args()
    result = find_audio(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
