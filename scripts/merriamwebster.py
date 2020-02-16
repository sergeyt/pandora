#!/usr/bin/env python

import requests
import json
import utils
from bs4 import BeautifulSoup
from models import File, Term
from utils import dump_json

NAME = 'merriam-webster'


def parse_btn(btn):
    lang = btn['data-lang'].replace('_', '/')
    dir = btn['data-dir']
    file = btn['data-file']
    return f'https://media.merriam-webster.com/audio/prons/{lang}/mp3/{dir}/{file}.mp3'


def stripped_text(node):
    if node is None:
        return None
    return node.get_text().strip()


def get_data(query, lang):
    if lang != 'en':
        return

    url = f'https://www.merriam-webster.com/dictionary/{query}'

    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')

    #find transcription and audio
    prs = soup.find('span', class_='prs')

    if prs is not None:
        transcription = prs.find('span', class_='pr')
        transcription = stripped_text(transcription)

        yield ('transcription',
               Term(text=transcription, lang='ipa', region=None))

        btns = prs.find_all('a', class_='play-pron')
        urls = [parse_btn(b) for b in btns]
        urls = [u for u in urls if utils.url_exists(u)]
        for url in urls:
            yield ('audio', File(url=url, region=None))

    #find definitions and 'in'

    vg = soup.find_all('div', class_='vg')

    for v in vg:
        definitions = v.find_all(class_='dt')
        for d in definitions:
            text = stripped_text(d)
            #all defenitions start with ':' with class mw_t_bc
            if (d.find(class_='mw_t_bc') is not None):
                text = text.lstrip(':').strip()
                #with defenitions we can take examples of text with class ex-sent, we need drop it
                if (d.find(class_='ex-sent') is not None):
                    text = text.split('\n')[0].strip()
                yield ('definition', Term(text=text, lang=lang, region=None))
    #parse examples
    data_in = soup.find_all(class_='ex-sent')
    for d in data_in:
        if ('t' in d['class']):
            yield ('in', Term(text=stripped_text(d), lang=lang, region=None))
    #parse related
    ure = soup.find_all(class_='ure')
    for d in ure:
        yield ('related', Term(text=stripped_text(d), lang=lang, region=None))
    #parse tags
    tag = soup.find_all('span', class_='fl')
    for d in tag:
        yield ('tag', Term(text=stripped_text(d), lang=lang, region=None))

    #add tag with name 'word', becouse our name is word
    yield ('tag', Term(text='word', lang=lang, region=None))

    #move to second page, in teasaurus
    url_t = f'https://www.merriam-webster.com/thesaurus/{query}'
    resp = requests.get(url_t, headers=headers)
    if resp.ok:
        for t in parse_thesaurus(lang, resp.text):
            yield t


def parse_thesaurus(lang, page):
    soup = BeautifulSoup(page, 'html.parser')

    dlist = soup.find_all('span', class_='syn-list')
    for d in dlist:
        synonyms = d.find_all('a')
        for s in synonyms:
            yield ('synonym',
                   Term(text=stripped_text(s), lang=lang, region=None))

    dlist = soup.find_all('span', class_='rel-list')
    for d in dlist:
        related = d.find_all('a')
        for r in related:
            yield ('related',
                   Term(text=stripped_text(r), lang=lang, region=None))

    dlist = soup.find_all('span', class_='ant-list')
    for d in dlist:
        antonyms = d.find_all('a')
        for r in antonyms:
            yield ('antonym',
                   Term(text=stripped_text(r), lang=lang, region=None))


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(dump_json(result))


if __name__ == '__main__':
    main()
