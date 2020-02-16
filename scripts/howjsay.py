#!/usr/bin/env python

import json
from urllib.parse import quote
import utils
from models import File
from utils import dump_json

NAME = 'howjsay'


def get_data(text, lang='en'):
    if lang != 'en':
        return None

    url = f'https://howjsay.com/mp3/{quote(text)}.mp3'
    if utils.url_exists(url):
        yield ('audio', File(url=url, region=None))


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(dump_json(result))


if __name__ == '__main__':
    main()
