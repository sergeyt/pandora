#!/usr/bin/env python

import re
import requests
import json
import utils
from models import File

first = lambda a: next(iter(a or []), None)


def get_data(text, lang='en'):
    if lang != 'en':
        return None

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

    data = {
      'audio': []
    }
    for url in [mp3, ogg]:
      data['audio'].append(File(url=url, region=None))

    return data


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
