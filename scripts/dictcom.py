#!/usr/bin/env python

import sys
import urllib
import requests
from bs4 import BeautifulSoup
from localcache import Cache

ru_en = 'russian-english'
en_ru = 'english-russian'
cache = Cache('dictcom')


def translate(text, lang=ru_en):
    result = cache.get(text)
    if result is not None:
        return result

    pat = 'https://www.dict.com/{0}/{1}'
    url = pat.format(lang, urllib.parse.quote(text))
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    trans = soup.find('span', class_='lex_ful_tran')
    if trans is None:
        return None
    result = {
        'lang': short_lang(lang),
        'tran': [s.strip() for s in trans.contents[0].split(',')]
    }
    cache.put(text, result)
    return result


def short_lang(lang):
    if lang == ru_en:
        return 'ru-en'
    if lang == en_ru:
        return 'en-ru'
    return lang


def main():
    lang = sys.argv[2] if len(sys.argv) >= 3 else ru_en
    print(translate(sys.argv[1], lang))


if __name__ == '__main__':
    main()
