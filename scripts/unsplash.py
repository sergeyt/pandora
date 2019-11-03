#!/usr/bin/env python

import sys
import json
import requests
import utils
from bs4 import BeautifulSoup
from models import Term, File

headers = {
    'User-Agent': utils.CHROME_USER_AGENT,
    'Accept': 'text/html',
}


def get_data(text, lang='en'):
    base = 'https://unsplash.com'
    pat = '{0}/s/photos/{1}'
    url = pat.format(base, text.replace(' ', '-'))

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
                data['visual'].append(File(url=src, region=None))

    return data


def main():
    (text, lang) = utils.find_audio_args(sys.argv)
    result = get_data(text)
    print(json.dumps(result, sort_keys=True, indent='  ', ensure_ascii=False))


if __name__ == '__main__':
    main()