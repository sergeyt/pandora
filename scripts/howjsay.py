#!/usr/bin/env python

import json
import utils
from localcache import Cache

cache = Cache('howjsay')


def find_audio(text, lang='en'):
    if lang != 'en':
        return None

    result = cache.get(text)
    if result is not None:
        return result

    url = 'https://howjsay.com/mp3/{0}.mp3'.format(text)
    if not utils.url_exists(url):
        return None

    result = {'mp3': [{'url': url}]}
    cache.put(text, result)
    return result


def main():
    (text, lang) = utils.find_audio_args()
    result = find_audio(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
