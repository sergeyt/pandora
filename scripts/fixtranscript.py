#!/usr/bin/env python

import os
import re
import api

dir = os.path.dirname(os.path.realpath(__file__))

lines = []
current_id = None
props = {}


def fix_transcripts(props, id):
    k = 'transcript'
    if k not in props:
        return []

    buf = []
    for key, val in props[k].items():
        tid = id + '_transcription_' + key
        lang = 'ipa'
        region = 'us'
        if key == 'ru':
            lang = 'ru'
            region = 'ru'
        buf.extend(
            api.nquads({
                'text': val,
                'lang': lang,
                'region': region,
            }, tid))
        buf.append('')
        t = props.get('transcription', [])
        t.append('_:' + tid)
        props['transcription'] = t

    props.pop(k, None)
    return buf


def commit():
    global current_id
    global props
    if current_id is None:
        return
    trans = fix_transcripts(props, current_id)
    lines.extend(trans)
    buf = api.nquads(props, current_id)
    lines.extend(buf)
    current_id = None
    props = {}
    return None


with open(os.path.join(dir, 'lingvo.txt'), 'r', encoding='utf-8') as f:
    src = f.read().split('\n')

    for line in src:
        m = re.match(r'_:(\w+_(en|ru))\s*<(\w+)>\s*(.*)\s*\.', line)
        if m is None:
            lines.append(line)
            continue

        id = m.group(1)
        pred = m.group(3)
        val = m.group(4)

        m = re.match(r'\s*"([^"]*)"(@(\w+))?\s*', val)
        if m is not None:
            val = m.group(1)
            lang = m.group(3)
            if lang is not None:
                t = props.get(pred, {})
                t[lang] = val
                val = t
        else:
            val = val.strip()

        if current_id is not None and id != current_id:
            commit()

        current_id = id
        props[pred] = val

    commit()
    out = '\n'.join(lines)
    print(out)
