#!/usr/bin/env python

import sys
import os
import re
from datetime import datetime
import api
import macmillan
import forvo
import urllib

first = lambda a: next(iter(a or []), None)
testing = first(filter(lambda a: a.find('testing') >= 0, sys.argv)) is not None
dir = os.path.dirname(os.path.realpath(__file__))

# utils.enable_logging_with_headers()


def login():
    global user
    global user_id
    if testing:
        user = {'uid': '0x1', 'name': 'system'}
        user_id = user['uid']
        return
    api.login("system", os.getenv("SYSTEM_PWD"))
    user = api.current_user()
    user_id = user['uid']


def proxy_url(url):
    return url if testing else api.fileproxy(url)


def created_at(id):
    now = datetime.now().isoformat()
    return '{0} <created_at> "{1}"^^<xs:dateTime> .'.format(id, now)


def created_by(id):
    return '{0} <created_by> <{1}> .'.format(id, user_id)


def init():
    global user
    global user_id
    login()

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
                if kind == 'term':
                    add_audio(line, id, buf, audio)
                if line.find('<Term>') < 0:
                    buf.append('{0} <{1}> "" .'.format(id, kind.capitalize()))
                buf.append(created_at(id))
                buf.append(created_by(id))
                typed[id] = True

            if "<visual>" in line and id in audio:
                for t in audio[id]:
                    buf.append(t)

            buf.append(line)

    data = '\n'.join(buf)
    print(data)
    if not testing:
        api.set_nquads(data)


def map_type(line):
    if line.startswith('_:img_'):
        return 'file'
    if re.match(r'_:\w+_(en|ru)', line):
        return 'term'
    return None


def idof(line):
    return line.split(' ')[0]


def audio_nquads(term_id, url, i):
    term_id = term_id.lstrip('_').lstrip(':')
    u = urllib.parse.urlparse(url)
    src = '{0}://{1}'.format(u.scheme, u.netloc)
    id = '_:aud_{0}{1}'.format(term_id, i)
    nquads = [
        '{0} <File> "" .'.format(id),
        created_at(id),
        created_by(id),
        '{0} <url> "{1}" .'.format(id, url),
        '{0} <source> "{1}" .'.format(id, src),
    ]
    return id, nquads


def add_audio(line, id, buf, audio):
    m = re.match(r'_:(\w+)_(en|ru)\s*<text>\s*"([^"]+)"\s*\.', line)
    if m is None:
        return

    word = m.group(1)
    if word.find('_') >= 0:
        return

    lang = m.group(2)
    text = m.group(3)

    urls = []
    if lang == 'en':
        m = macmillan.find_audio(word)
        urls = ['https://howjsay.com/mp3/{0}.mp3'.format(word), m['mp3']]

    f = forvo.find_audio(text, lang)
    if f is None and len(urls) == 0:
        return
    urls.extend([t['url'] for t in f['mp3']])

    proxy_urls = []
    for url in urls:
        try:
            url2 = proxy_url(url)
            proxy_urls.append(url2)
        except:
            print('cannot proxy {0}'.format(url))

    audio[id] = []
    lines = []
    for i, url in enumerate(proxy_urls):
        (aud_id, nquads) = audio_nquads(id, url, i + 1)
        audio[id].append('{0} <audio> {1} .'.format(id, aud_id))
        lines.extend(nquads)
    for t in lines:
        buf.append(t)
    buf.append('')


def change_url(line):
    m = re.match(r'_:([\w_]+) <url> "([^"]*)"', line)
    if m is None:
        return line
    id = m.group(1)
    image_url = m.group(2)
    placeholder = 'https://imgplaceholder.com/420x320/ff7f7f/333333/fa-image'
    if image_url == '':
        image_url = placeholder
    try:
        image_url = proxy_url(image_url)
    except:
        try:
            image_url = proxy_url(placeholder)
        except:
            image_url = placeholder
    return '_:{0} <url> "{1}" .'.format(id, image_url)


if __name__ == '__main__':
    init()
