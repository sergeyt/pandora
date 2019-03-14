#!/usr/bin/env python

import sys
import re
import requests

first = lambda a: next(iter(a or []), None)


def find_audio(word):
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
    return {'mp3': first(mp3), 'ogg': first(ogg)}


def main():
    word = sys.argv[1]
    print(find_audio(word))


if __name__ == '__main__':
    main()
