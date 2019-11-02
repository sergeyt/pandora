#!/usr/bin/env python

import json
import utils
from models import File


def get_data(text, lang='en'):
    if lang != 'en':
        return None

    url = 'https://howjsay.com/mp3/{0}.mp3'.format(text)
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
