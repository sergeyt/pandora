#!/usr/bin/env python

import json
from urllib.parse import quote
import utils
from models import File

NAME = 'howjsay'


def get_data(text, lang='en'):
    if lang != 'en':
        return None

    url = f'https://howjsay.com/mp3/{quote(text)}.mp3'
    if not utils.url_exists(url):
        return None

    data = {'audio': [File(url=url, region=None)]}

    return data


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
