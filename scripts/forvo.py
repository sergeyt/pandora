#!/usr/bin/env python

import sys
import os
import urllib
import requests
import base64
import json
from bs4 import BeautifulSoup

AUDIO_HOST = 'https://audio00.forvo.com/audios/mp3'

first = lambda a: next(iter(a or []), None)

dir = os.path.dirname(os.path.realpath(__file__))
cachePath = os.path.join(dir, 'forvo.json')
cache = {}
testing = False


def decode_base64(s):
    return base64.b64decode(s).decode('utf-8')


def unquote(s):
    if s.startswith('\''):
        return s[1:len(s) - 1]
    return s


def parse_fn(src):
    if not src:
        return None
    i = src.find('(')
    j = src.find(')')
    name = src[:i]
    args = [unquote(s) for s in src[i + 1:j].split(',')]
    return {'name': name, 'args': args}


def translate_gender(val):
    return val.strip()


def translate_counry(val):
    return val.strip()


def parse_from(s):
    if not s:
        return None
    s = s.strip('(').strip(')')
    a = s.split(',')
    if len(a) == 0:
        return None
    result = {'gender': translate_gender(a[0])}
    if len(a) == 2:
        result['country'] = translate_counry(a[1])
    return result


def parse_item(item):
    btn = item.find('span', class_='play')
    if btn is None:
        return None

    fn = parse_fn(btn['onclick'])
    if fn is None or fn['name'] != 'Play':
        return None
    rel = decode_base64(fn['args'][4])
    url = '{0}/{1}'.format(AUDIO_HOST, rel)
    if not url.endswith('.mp3'):
        return None

    result = {'url': url}
    author = item.find('span', class_='ofLink')
    if author and 'data-p2' in author.attrs:
        result['author'] = author.attrs['data-p2']

    from_tag = item.find('span', class_='from')
    if from_tag:
        d = parse_from(from_tag.contents[0])
        if d:
            for k, v in d.items():
                result[k] = v

    return result


def find_in_cache(word):
    global cache
    if len(cache) == 0:
        with open(cachePath, 'r') as f:
            cache = json.load(f)
    return cache[word] if word in cache else None


def find_audio(word, lang='ru'):
    result = find_in_cache(word)
    if not testing and result is not None:
        return result

    pat = 'https://ru.forvo.com/word/{0}/#{1}'
    url = pat.format(urllib.parse.quote(word), lang)
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    article = soup.find('article', class_='pronunciations')
    if article is None:
        return None
    ul = article.find('ul', class_="show-all-pronunciations")
    if ul is None:
        return None

    li = ul.find_all('li')
    parsed_items = [parse_item(t) for t in li]
    items = [t for t in parsed_items if t is not None]
    result = {'mp3': items}

    if not testing:
        cache[word] = result
        with open(cachePath, 'w') as f:
            f.write(json.dumps(cache, sort_keys=True, indent='  '))

    return result


def main():
    word = sys.argv[1]
    print(find_audio(word))


if __name__ == '__main__':
    main()
