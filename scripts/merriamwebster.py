#!/usr/bin/env python

import requests
import json
import utils
from bs4 import BeautifulSoup
from models import File,Term


def parse_btn(btn):
    lang = btn['data-lang'].replace('_', '/')
    dir = btn['data-dir']
    file = btn['data-file']
    pat = 'https://media.merriam-webster.com/audio/prons/{0}/mp3/{1}/{2}.mp3'
    return pat.format(lang, dir, file)

def get_data(text, lang):
    if lang != 'en':
        return None
        
        
    data = {
        'audio': [],
        #'visual': [],
        'tag': [],
        'transcription': [],
        'definition': [],
        'in': [],
        'synonym':[],
        'antonym':[],
        'related':[]
    }
    
    
    pat = 'https://www.merriam-webster.com/dictionary/{0}'
    pat_t= 'https://www.merriam-webster.com/thesaurus/{0}'
    url = pat.format(text)
    url_t = pat_t.format(text)
    
    headers = {
        'User-Agent': 'script',
        'Accept': 'text/html',
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()
    
    soup = BeautifulSoup(resp.text, 'html.parser')
    
    #find transcription and audio
    prs=soup.find('span',class_='prs')
    
    transcription=prs.find('span',class_='pr').text
    
    btns=prs.find_all('a',class_='play-pron')
    urls = [parse_btn(b) for b in btns]
    urls = [u for u in urls if utils.url_exists(u)]
    for url in urls:
        data['audio'].append(File(url=url, region=None))
    data['transcription'].append(transcription)
    
    #find definitions and in
    
    vg=soup.find_all('div',class_='vg')
    
    for v in vg:
        definitions=v.find_all(class_='dt')
        for d in definitions:
            text=d.get_text().strip()
            if(d.find(class_='mw_t_bc') is not None):
                text=text[2:]
                if(d.find(class_='ex-sent') is not None):
                    text=text.split('\n')[0].strip()
                data['definition'].append(Term(text=text,lang=None,region=None))

    data_in=soup.find_all(class_='ex-sent')
    for d in data_in:
            if('t' in d['class']):
                data['in'].append(Term(text=d.get_text().strip(),lang=None,region=None))

    ure=soup.find_all(class_='ure')
    for d in ure:
        data['related'].append(Term(text=d.text,lang=None,region=None))
        
    tag=soup.find_all('span',class_='fl')
    for t in tag:
        data['tag'].append(Term(text=t.text,lang=None,region=None))    
    
    data['tag'].append(Term(text='word',lang=None,region=None))    
        
        
        
    resp = requests.get(url_t, headers=headers)
    resp.raise_for_status()
    
    soup = BeautifulSoup(resp.text, 'html.parser')
    
    slist=soup.find_all('span',class_='syn-list')
    for d in slist:
        #list=slist[i].find('div',class_='thes-list-content')
        synonyms=d.find_all('a')
        for s in synonyms:
            data['synonym'].append(Term(text=s.text,lang=None,region=None))
    
    rlist=soup.find_all('span',class_='rel-list')
    for d in rlist:
        related=d.find_all('a')
        for r in related:
            data['related'].append(Term(text=r.text,lang=None,region=None))
            
    rlist=soup.find_all('span',class_='ant-list')
    for d in rlist:
        antonyms=d.find_all('a')
        for r in antonyms:
            data['antonym'].append(Term(text=r.text,lang=None,region=None))
            
    return data
        
    
        
        
        
        
def main():
    (text, lang) = utils.find_audio_args()
    result = get_data(text, lang)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
