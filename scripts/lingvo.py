#!/usr/bin/env python

import os
import re
import api
import macmillan

dir = os.path.dirname(os.path.realpath(__file__))


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
    
    m = macmillan.find_audio(word)
    src2 = 'https://www.macmillandictionary.com'
    lines = [
        '_:aud_{0}1 <url> "https://howjsay.com/mp3/{0}.mp3" .'.format(word),
        '_:aud_{0}1 <source> "https://howjsay.com" .'.format(word),
        '_:aud_{0}2 <url> "{1}" .'.format(word, m['mp3']),
        '_:aud_{0}2 <source> "{1}" .'.format(word, src2),
    ]
    for t in lines:
        buf.append(t)
    buf.append('')

    audio[id] = []
    for i in range(1, 3):
        audio[id].append('{0} <audio> _:aud_{1}{2} .'.format(id, word, i))


def init():
    with open(os.path.join(dir, 'lingvo.txt'), 'r', encoding='utf-8') as f:
        src = f.read().split('\n')
        buf = []
        typed = {}
        audio = {}
        for line in src:
            kind = map_type(line)
            id = idof(line)
            if kind and id not in typed:
                add_audio(line, id, buf, audio)
                buf.append('{0} <_{1}> "" .'.format(id, kind))
                typed[id] = True

            if "<visual>" in line and id in audio:
                for t in audio[id]:
                    buf.append(t)

            buf.append(line)

            if "<visual>" in line:
                buf.append(line.replace("<visual>", "<relevant>"))
                buf.append(line.replace("<visual>", "<related>"))

    data = '\n'.join(buf)
    print(data)
    api.set_nquads(data)


if __name__ == '__main__':
    init()
