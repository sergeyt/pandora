#!/usr/bin/env python

import os
import re
import api

dir = os.path.dirname(os.path.realpath(__file__))


def map_type(line):
    if line.startswith('_:img_'):
        return 'file'
    if re.match("_:\w+_(en|ru)", line):
        return 'term'
    return None


def idof(line):
    return line.split(' ')[0]


def init():
    with open(os.path.join(dir, 'lingvo.txt'), 'r', encoding='utf-8') as f:
        src = f.read().split('\n')
        buf = []
        typed = {}
        for line in src:
            buf.append(line)
            kind = map_type(line)
            if kind:
                id = idof(line)
                if id not in typed:
                    buf.append('{0} <_{1}> "" .'.format(id, kind))
                    typed[id] = True
            if "<visual>" in line:
                buf.append(line.replace("<visual>", "<relevant>"))
                buf.append(line.replace("<visual>", "<related>"))

    data = '\n'.join(buf)
    api.set_nquads(data)


if __name__ == '__main__':
    init()
