#!/usr/bin/env python

import json
import howjsay
import cambridge
import macmillan
import merriamwebster
import forvo
import utils


def find_audio(text, lang):
    result = {}
    sources = [
        howjsay.find_audio, cambridge.find_audio, macmillan.find_audio,
        merriamwebster.find_audio, forvo.find_audio
    ]
    for source in sources:
        a = source(text, lang)
        if a is None:
            continue
        for fmt in ['mp3', 'ogg']:
            if fmt not in a:
                continue
            if fmt not in result:
                result[fmt] = []
            result[fmt].extend(a[fmt])
    return None if len(result) == 0 else result


def main():
    (text, lang) = utils.find_audio_args()
    result = find_audio(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  '))


if __name__ == '__main__':
    main()
