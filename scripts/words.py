#!/usr/bin/env python
# -*- coding: utf-8 -*-

import utils
import api

words = [
    {
        'text@en': 'house',
        'text@ru': 'дом',
        'transcription@en': 'hause',
        'transcription@ru': "'хаус",
        'pronunciation@en': 'https://howjsay.com/mp3/house.mp3',
        'images': [
            'http://epicpix.com/wp-content/uploads/2016/04/ff_3280.jpg',
        ],
    },
    {
        'text@en': 'lake',
        'text@ru': 'озеро',
        'transcription@en': 'leik',
        'transcription@ru': "'лэйк",
        'pronunciation@en': 'https://howjsay.com/mp3/lake.mp3',
        'images': [
            'https://github.com/flutter/website/blob/master/src/_includes/code/layout/lakes/images/lake.jpg?raw=true',
        ],
    },
    {
        "text@en": "gum",
        "text@ru": "жвачка",
        "transcription@ru": "гам",
        'pronunciation@en': 'https://howjsay.com/mp3/gum.mp3',
        "images": [
            "https://i.pinimg.com/originals/39/10/30/3910305670e1fb3f0584e998e30f1b71.jpg",
            "https://images-na.ssl-images-amazon.com/images/I/81F-d7PAZQL._SL1500_.jpg",
            "https://cdn.shopify.com/s/files/1/0004/8132/9204/products/double-mint-gum_1024x1024.jpg?v=1522355731",
        ],
    },
    {
        "text@en": "recreation",
        "text@ru": "отдых",
        "transcription@ru": "рикри`эйшен",
        'pronunciation@en': 'https://howjsay.com/mp3/recreation.mp3',
        "images": [
            "http://www.bridgewaternj.gov/wp-content/uploads/Images/recphoto.jpg",
            "https://vmcdn.ca/f/files/halifaxtoday/images/sports/071818-sports-equipment-recreation-gym-fitness-adobestock_190038155.jpeg;w=630",
            "http://www.worldwar-collectibles.com/wp-content/uploads/2017/03/recreation.jpg",
        ],
    },
    {
        "text@en": "apple",
        "text@ru": "яблоко",
        "transcription@ru": "эпл",
        'pronunciation@en': 'https://howjsay.com/mp3/apple.mp3',
        "images": [
            "http://static1.squarespace.com/static/5849b12a2e69cf47aecece6b/584ebb9646c3c416aac4f2b5/5b830f4dcd8366d1f2d15de9/1535382024204/apple.jpg?format=1500w",
        ],
    },
    {
        "text@en": "table",
        "text@ru": "стол",
        "transcription@ru": "тэйбл",
        'pronunciation@en': 'https://howjsay.com/mp3/table.mp3',
        "images": [
            "https://cdn.shopify.com/s/files/1/2660/5106/products/wisz8mrpd67l6pss3crw_2b63262e-9b7a-4bbb-8f4f-d2da5d3cc57f_800x.jpg?v=1539039199",
        ],
    },
    {
        "text@en": "spoon",
        "text@ru": "ложка",
        "transcription@ru": "спун",
        'pronunciation@en': 'https://howjsay.com/mp3/spoon.mp3',
        "images": [
            "https://coubsecure-s.akamaihd.net/get/b45/p/coub/simple/cw_timeline_pic/f357228b51f/972920dc0d51d3ed34e9d/big_1409081702_1382451563_att-migration20121219-1328-1jpmbyq.jpg",
        ],
    },
]


def rdf_repr(v):
    if isinstance(v, str):
        return '"{0}"'.format(v)
    return v


def nquad(id, k, v):
    a = k.split('@')
    p = a[0]
    lang = a[1] if len(a) == 2 else ''
    s = rdf_repr(v)
    if len(lang) > 0:
        s += "@{0}".format(lang)
    return "_:{0} <{1}> {2} .\n".format(id, p, s)


def nquads(d, id):
    result = ''
    for k, v in d.iteritems():
        result += nquad(id, k, v)
    return result


for word in words:
    # insert images
    images = []
    for url in word['images']:
        data = nquads({
            'url': url,
            'source': 'google',
        }, "img")
        resp = api.set_nquads(data)
        images.append(resp['data']['uids']['img'])

    # insert word
    word.pop('images', None)
    data = nquads(word, 'w')
    resp = api.set_nquads(data)
    wid = resp['data']['uids']['w']
    
    # link word with images
    for img in images:
        data = "<{0}> <image> <{1}> .\n".format(wid, img)
        data += "<{0}> <word> <{1}> .\n".format(img, wid)
        api.set_nquads(data)
