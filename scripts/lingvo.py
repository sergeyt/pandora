#!/usr/bin/env python

import os
import json
import utils
import api

dir = os.path.dirname(os.path.realpath(__file__))

with open(os.path.join(dir, 'lingvo.txt'), 'r', encoding='utf-8') as f:
    data = f.read()


def init():
    api.set_nquads(data)


if __name__ == '__main__':
    init()
