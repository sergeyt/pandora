#!/usr/bin/env python

import requests
import json
from urllib.parse import quote
import utils

from bs4 import BeautifulSoup
from models import File, Term

from utils import dump_json
NAME = 'macmillan'

first = lambda a: next(iter(a or []), None)


def stripped_text(node):
    if node is None:
        return None
    return node.get_text().strip()


def get_data(query, lang):
    if lang != 'en':
        return None

    data = {
        'audio': [],
        #'visual': [],
        'tag': [],
        'transcription': [],
        'definition': [],
        'in': [],
        'synonym': [],
        #'antonym': [],
        #'related': []
    }

    pat = 'https://www.macmillandictionary.com/dictionary/british/{0}'
    url = pat.format(query)
    headers = {
            'User-Agent': 'script',
            'Accept': 'text/html',
        }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, 'html.parser')

    #get transcription
    transcriptions=soup.find_all(class_='PRON')
    for t in transcriptions:
      data['transcription'].append(Term(text=stripped_text(t).replace('/',''),lang=lang, region=None))

    #get tags
    crop_text=stripped_text(soup.find(class_='zwsp'))
    part_speech=stripped_text(soup.find(class_='PART-OF-SPEECH')).replace(crop_text,'')
    syntax_coding=stripped_text(soup.find(class_='SYNTAX-CODING'))
    
    data['tag'].append(Term(text=part_speech,lang=lang, region=None))
    data['tag'].append(Term(text=syntax_coding,lang=lang, region=None))

    #get defenition
    defenitions=soup.find_all(class_='DEFINITION')
    for d in defenitions:
      data['definition'].append(Term(text=stripped_text(d),lang=lang, region=None))

    #get examples
    examples=soup.find_all(class_='EXAMPLES')
    for e in examples:
      data['in'].append(Term(text=stripped_text(e),lang=lang, region=None))
    examples=soup.find_all(class_='PHR-XREF')
    for e in examples:
      data['in'].append(Term(text=stripped_text(e),lang=lang, region=None))

    #get synonyms
    synonyms=soup.find_all(class_='synonyms')
    for allsyn in synonyms:
      subsynonyms=allsyn.find_all(class_='theslink')
      for syn in subsynonyms:
        if(not '...' in syn.text):
          data['synonym'].append(Term(text=stripped_text(syn),lang=lang, region=None))

    #get audio
    audio=soup.find(class_='audio_play_button')
    data['audio'].append(File(url=audio['data-src-mp3'], region=None))
    data['audio'].append(File(url=audio['data-src-ogg'], region=None))
    
    return data


def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(dump_json(result))


if __name__ == '__main__':
    main()
