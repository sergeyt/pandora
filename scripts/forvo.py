#!/usr/bin/env python

import sys
import os
import urllib
import requests
import base64
import json
from bs4 import BeautifulSoup

AUDIO_HOST = 'https://audio00.forvo.com/audios/mp3'

dir = os.path.dirname(os.path.realpath(__file__))
cachePath = os.path.join(dir, 'forvo.json')
cache = {}


def decode_base64(s):
    return base64.b64decode(s).decode('utf-8')


def unquote(s):
    if s.startswith('\''):
        return s[1:len(s) - 1]
    return s


def parse_fn(src):
    i = src.find('(')
    j = src.find(')')
    name = src[:i]
    args = [unquote(s) for s in src[i + 1:j].split(',')]
    return {'name': name, 'args': args}


def find_in_cache(word):
    global cache
    if len(cache) == 0:
        with open(cachePath, 'r') as f:
            cache = json.load(f)
    return cache[word] if word in cache else None


def find_audio(word):
    result = find_in_cache(word)
    if result is not None:
        return result

    pat = 'https://ru.forvo.com/word/{0}/#ru'
    url = pat.format(urllib.parse.quote(word))
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    btns = soup.find_all('span', class_="play")
    fns = [parse_fn(b['onclick']) for b in btns]
    rel = [decode_base64(f['args'][4]) for f in fns if f['name'] == 'Play']
    urls = ['{0}/{1}'.format(AUDIO_HOST, r) for r in rel]
    mp3 = [u for u in urls if u.endswith('.mp3')]

    result = {'mp3': mp3}

    cache[word] = result
    with open(cachePath, 'w') as f:
        f.write(json.dumps(cache, sort_keys=True, indent='  '))

    return result


def main():
    word = sys.argv[1]
    print(find_audio(word))


if __name__ == '__main__':
    main()
