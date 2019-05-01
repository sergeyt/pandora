#!/usr/bin/env python

import os
import re
from datetime import datetime
import api
import macmillan
import forvo
import requests
import utils
import urllib

testing = False
now = datetime.now().isoformat()
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
                if kind != 'tag':
                    buf.append('{0} <{1}> "" .'.format(id, kind.capitalize()))
                created_at = '{0} <created_at> "{1}"^^<xs:dateTime> .'
                created_by = '{0} <created_by> <{1}> .'
                buf.append(created_at.format(id, now))
                buf.append(created_by.format(id, user_id))
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
    if line.find('<Tag>') >= 0:
        return 'tag'
    return None


def idof(line):
    return line.split(' ')[0]


def aud_nquads(id, url, i):
    id = id.lstrip('_').lstrip(':')
    u = urllib.parse.urlparse(url)
    src = '{0}://{1}'.format(u.scheme, u.netloc)
    aud = '_:aud_{0}{1}'.format(id, i)
    t1 = '{0} <url> "{1}" .'.format(aud, url)
    t2 = '{0} <source> "{1}" .'.format(aud, src)
    return aud, [t1, t2]


def add_audio(line, id, buf, audio):
    m = re.match(r'_:(\w+)_(en|ru)\s*<text>\s*"([^"]+)"\s*\.', line)
    if m is None:
        return

    word = m.group(1)
    if word.find('_') >= 0:
        return

    lang = m.group(2)
    text = m.group(3)

    if lang == 'ru':
        f = forvo.find_audio(text)
        if f is None:
            return
        urls = [t['url'] for t in f['mp3']]
    else:
        m = macmillan.find_audio(word)
        urls = ['https://howjsay.com/mp3/{0}.mp3'.format(word), m['mp3']]

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
        (aud_id, nquads) = aud_nquads(id, url, i + 1)
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
