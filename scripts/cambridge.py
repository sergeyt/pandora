#!/usr/bin/env python

import sys
import requests
import json
import utils
from bs4 import BeautifulSoup
from models import Term, File

headers = {
    'User-Agent':
    'Mozilla/5.0 AppleWebKit/537.36 (KHTML like Gecko) Chrome/78.0.3883.121',
    'Accept': 'text/html',
}


def stripped_text(node):
    if node is None:
        return None
    return node.get_text().strip()


def find_strip(container, tag, class_):
    node = container.find(tag, class_=class_)
    return stripped_text(node)


def get_data(text, lang):
    if lang != 'en':
        return None

    base = 'https://dictionary.cambridge.org'
    pat = '{0}/dictionary/english/{1}'
    url = pat.format(base, text)

    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    data = {
        'tag': [],
        'audio': [],
        'transcription': [],
        'definition': [],
        'in': [],
        'visual': [],
    }

    if utils.is_word(text):
        data['tag'].append(Term(text='word', lang=lang, region=None))

    soup = BeautifulSoup(resp.text, 'html.parser')
    page = soup.find('div', class_='page')
    for dictionary in page.find_all('div', class_='dictionary'):
        header = dictionary.find('div', class_='pos-header')
        body = dictionary.find('div', class_='pos-body')

        posgram = header.find('div', class_='posgram')
        if posgram:
            pos = find_strip(posgram, 'span', class_='pos')
            data['tag'].append(Term(text=pos, lang=lang, region=None))
        # TODO parse codes like countable, etc

        # parse pronunciations
        for dpron in header.find_all('span', class_='dpron-i'):
            region = find_strip(dpron, 'span', 'region')
            amp = header.find('amp-audio')
            for source in amp.find_all('source'):
              file = File(url=base + source.attrs['src'], region=region)
              data['audio'].append(file)

            ipa = find_strip(dpron, 'span', class_='ipa')
            data['transcription'].append(
                Term(text=ipa, lang=lang, region=region))

        for dblock in body.find_all('div', class_='def-block'):
            def_text = stripped_text(dblock.find('div', class_='def'))
            data['definition'].append(
                Term(text=def_text, lang=lang, region=None))
            img = dblock.find('amp-img')
            if img is not None:
                file = File(url=base + img.attrs['src'], region=None)
                data['visual'].append(file)
            for eg in dblock.find_all('span', 'eg'):
                term = Term(text=stripped_text(eg), lang=lang, region=None)
                data['in'].append(term)

    for dataset in page.find_all('div', class_='dataset'):
        for eg in dataset.find_all('span', class_='deg'):
            term = Term(text=stripped_text(eg), lang=lang, region=None)
            data['in'].append(term)

    # TODO parse collocations

    return data


def find_audio(text, lang):
    data = get_data(text, lang)
    if data is None:
        return None
    return [a._asdict() for a in data['audio']]


def main():
    (text, lang) = utils.find_audio_args(sys.argv)
    result = get_data(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
