#!/usr/bin/env python

import os
import re

dir = os.path.dirname(os.path.realpath(__file__))

with open(os.path.join(dir, 'lingvo.txt'), 'r', encoding='utf-8') as f:
    src = f.read().split('\n')
    lines = []
    for line in src:
        m = re.match(r'(_:\w+_(en|ru))\s*<tag>\s*_:(\w+)\s*\.', line)
        if not m:
            lines.append(line)
            continue
        id = m.group(1)
        tag = m.group(3)
        if tag.endswith('_en') or tag.endswith('_ru'):
            if id[-2:] != tag[-2:]:
                continue
            lines.append(line)
            continue
        lines.append('{0} <tag> _:{1}_en .'.format(id, tag))
        lines.append('{0} <tag> _:{1}_ru .'.format(id, tag))
    out = '\n'.join(lines)
    print(out)
