#!/usr/bin/env python

import requests
import json
import utils
from bs4 import BeautifulSoup
from models import File


def parse_btn(btn):
    lang = btn['data-lang'].replace('_', '/')
    dir = btn['data-dir']
    file = btn['data-file']
    pat = 'https://media.merriam-webster.com/audio/prons/{0}/mp3/{1}/{2}.mp3'
    return pat.format(lang, dir, file)


def get_data(text, lang):
    if lang != 'en':
        return None

    pat = 'https://www.merriam-webster.com/dictionary/{0}'
    url = pat.format(text)
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')
    prs = soup.find('span', class_='prs')
    if prs is None:
        return None
    btns = prs.find_all('a', class_='play-pron')
    if len(btns) == 0:
        return None

    urls = [parse_btn(b) for b in btns]
    urls = [u for u in urls if utils.url_exists(u)]
    data = {
      'audio': []
    }
    for url in urls:
        data['audio'].append(File(url=url, region=None))

    return data


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
