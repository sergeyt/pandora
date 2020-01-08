#!/usr/bin/env python

import sys
import json
import requests
import utils
from bs4 import BeautifulSoup
from models import Term, File

NAME = 'unsplash'

headers = {
    'User-Agent': utils.CHROME_USER_AGENT,
    'Accept': 'text/html',
}


def get_data(text, lang='en'):
    base = 'https://unsplash.com'
    txt = text.replace(' ', '-')
    url = f'{base}/s/photos/{txt}'

    resp = requests.get(url, headers=headers)
    resp.raise_for_status()

    data = {'visual': []}

    soup = BeautifulSoup(resp.text, 'html.parser')
    container = soup.find('div', attrs={'data-test': 'search-photos-route'})
    for figure in container.find_all('figure'):
        for img in figure.find_all('img'):
            # srcset = img.attrs['srcset'].split(',')
            # print(srcset)
            # i = srcset.index('1000w')
            if 'src' in img.attrs:
                src = img.attrs['src']
                if not src.startswith('https://images.unsplash.com/profile-'):
                    data['visual'].append(File(url=src, region=None))

    return data


def main():
    (text, lang) = utils.find_audio_args(sys.argv)
    result = get_data(text)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()
