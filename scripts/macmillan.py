#!/usr/bin/env python

import sys
import os
import re
import requests
import json

first = lambda a: next(iter(a or []), None)
dir = os.path.dirname(os.path.realpath(__file__))
cachePath = os.path.join(dir, 'macmillan.json')
cache = {}


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

    pat = 'https://www.macmillandictionary.com/dictionary/british/{0}'
    url = pat.format(word)
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    mp3 = re.findall('data-src-mp3="([^"]+)"', resp.text)
    ogg = re.findall('data-src-ogg="([^"]+)"', resp.text)
    result = {'mp3': first(mp3), 'ogg': first(ogg)}

    cache[word] = result
    with open(cachePath, 'w') as f:
        f.write(json.dumps(cache, sort_keys=True, indent='  '))

    return result


def main():
    word = sys.argv[1]
    print(find_audio(word))


if __name__ == '__main__':
    main()
