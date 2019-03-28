#!/usr/bin/env python

import os
import re
from datetime import datetime
import api
import macmillan
import requests
import utils

dir = os.path.dirname(os.path.realpath(__file__))

# utils.enable_logging_with_headers()


def init():
    api.login("system", os.getenv("SYSTEM_PWD"))
    user = api.current_user()
    user_id = user['uid']
    now = datetime.now().isoformat()

    with open(os.path.join(dir, 'lingvo.txt'), 'r', encoding='utf-8') as f:
        src = f.read().split('\n')
        buf = []
        typed = {}
        audio = {}
        for line in src:
            kind = map_type(line)
            id = idof(line)

            if kind == 'file':
                line = change_url(line)

            if kind and id not in typed:
                add_audio(line, id, buf, audio)
                buf.append('{0} <{1}> "" .'.format(id, kind.capitalize()))
                buf.append('{0} <created_at> "{1}"^^<xs:dateTime> .'.format(
                    id, now))
                buf.append('{0} <created_by> <{1}> .'.format(id, user_id))
                typed[id] = True

            if "<visual>" in line and id in audio:
                for t in audio[id]:
                    buf.append(t)

            buf.append(line)

    data = '\n'.join(buf)
    print(data)
    api.set_nquads(data)


def map_type(line):
    if line.startswith('_:img_'):
        return 'file'
    if re.match(r'_:\w+_(en|ru)', line):
        return 'term'
    return None


def idof(line):
    return line.split(' ')[0]


def add_audio(line, id, buf, audio):
    m = re.match(r'_:(\w+)_en', line)
    if m is None:
        return

    word = m.group(1)
    if word.find('_') >= 0:
        return

    m = macmillan.find_audio(word)
    url1 = api.fileproxy('https://howjsay.com/mp3/{0}.mp3'.format(word))
    url2 = api.fileproxy(m['mp3'])
    src2 = 'https://www.macmillandictionary.com'

    lines = [
        '_:aud_{0}1 <url> "{1}" .'.format(word, url1),
        '_:aud_{0}1 <source> "https://howjsay.com" .'.format(word),
        '_:aud_{0}2 <url> "{1}" .'.format(word, url2),
        '_:aud_{0}2 <source> "{1}" .'.format(word, src2),
    ]
    for t in lines:
        buf.append(t)
    buf.append('')

    audio[id] = []
    for i in range(1, 3):
        audio[id].append('{0} <audio> _:aud_{1}{2} .'.format(id, word, i))


def change_url(line):
    m = re.match(r'_:([\w_]+) <url> "([^"]*)"', line)
    if m is None:
        return line
    id = m.group(1)
    image_url = m.group(2)
    if image_url == '':
        image_url = 'https://imgplaceholder.com/420x320/ff7f7f/333333/fa-image'
    image_url = api.fileproxy(image_url)
    return '_:{0} <url> "{1}" .'.format(id, image_url)


if __name__ == '__main__':
    init()
