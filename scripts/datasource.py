#!/usr/bin/env python

import json
import howjsay
import cambridge
import macmillan
import merriamwebster
import forvo
import unsplash
import utils
from utils import dump_json

sources = [howjsay, cambridge, macmillan, merriamwebster, forvo, unsplash]


def get_data(text, lang):
    for source in sources:
        data = source.get_data(text, lang)
        if data is None:
            continue
        for t in data:
            yield t


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(dump_json(result))


if __name__ == '__main__':
    main()
