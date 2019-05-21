#!/usr/bin/env python

import re
import requests
import json
import utils
from localcache import Cache

first = lambda a: next(iter(a or []), None)
cache = Cache('macmillan')


def find_audio(text, lang='en'):
    if lang != 'en':
        return None

    result = cache.get(text)
    if result is not None:
        return result

    pat = 'https://www.macmillandictionary.com/dictionary/british/{0}'
    url = pat.format(text)
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    mp3 = first(re.findall('data-src-mp3="([^"]+)"', resp.text))
    ogg = first(re.findall('data-src-ogg="([^"]+)"', resp.text))
    if mp3 is None and ogg is None:
        return None

    result = {}
    if mp3:
        result['mp3'] = [{'url': mp3}]
    if ogg:
        result['ogg'] = [{'url': ogg}]

    cache.put(text, result)

    return result


def main():
    (text, lang) = utils.find_audio_args()
    result = find_audio(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
