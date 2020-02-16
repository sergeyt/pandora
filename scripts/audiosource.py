#!/usr/bin/env python

import json
import howjsay
import cambridge
import macmillan
import merriamwebster
import forvo
import utils
from utils import dump_json


def find_audio(text, lang):
    result = {}
    sources = [howjsay, cambridge, macmillan, merriamwebster, forvo]
    for source in sources:
        data = source.get_data(text, lang)
        if data is None:
            continue
        for (k, file) in data:
            if k != 'audio':
                continue
            for fmt in ['mp3', 'ogg']:
                if not file.url.endswith(fmt):
                    continue
                if fmt not in result:
                    result[fmt] = []
                result[fmt].append(file.url)
    return None if len(result) == 0 else result


def main():
    (text, lang) = utils.find_audio_args()
    result = find_audio(text, lang)
    print(dump_json(result))


if __name__ == '__main__':
    main()
