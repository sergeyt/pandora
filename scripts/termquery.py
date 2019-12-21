import re
from utils import is_word


def is_empty(v):
    return len(v) == 0


def keys(d):
    return d.keys()


relation_map = {
    'translated_as': {
        'label': 'Translations',
    },
    'transcription': {
        'label': 'Transcriptions',
        'reverse_edge': 'transcription_of',
    },
    'transcription_of': {
        'label': 'Transcription for',
        'reverse_edge': 'transcription',
    },
    'definition': {
        'label': 'Definitions',
        'reverse_edge': 'definition_of',
    },
    'definition_of': {
        'label': 'Definition for',
        'reverse_edge': 'definition',
    },
    'in': {
        'label': 'Used in',
    },
    'related': {
        'label': 'Related Terms',
    },
    'synonym': {
        'label': 'Synonyms',
    },
    'antonym': {
        'label': 'Antonyms',
    },
}

KIND = ['term', 'terms', 'audio', 'visual'] + list(keys(relation_map))

META = """
    created_at
    created_by {
      uid
      name
    }
"""

TAG = """tag {{
  uid
  text
  lang
  {0}
}}""".format(META)

TERM_BODY = """{{
  uid
  text
  lang
  {0}
  {1}
}}""".format(META, TAG)

FILE_BODY = """{{
  uid
  url
  source
  content_type
  views: count(see)
  likes: count(like)
  dislikes: count(dislike)
  {0}
  {1}
}}""".format(META, TAG)


def make_term_query(kind='terms',
                    term_id='',
                    offset=0,
                    limit=10,
                    lang='',
                    search_string='',
                    tags=[],
                    only_tags=False,
                    exact_match=False,
                    no_links=False):
    if not kind or kind not in KIND:
        raise Exception(f"invalid kind {kind}")
    if kind == 'term' and not term_id:
        raise Exception('termUid is required')
    if search_string is None:
        search_string = ''

    has_term_type = 'has(Term)'
    match_fn = f"uid({term_id})" if term_id else has_term_type
    is_term_list = kind == 'terms'
    is_term = kind == 'term'
    has_tag_type = 'has(Tag)' if is_term_list and only_tags else ''
    params = {}

    def make_search_filter():
        if not is_term_list:
            return ''

        str = search_string.strip()
        if not str:
            return ''

        params['$searchString'] = str
        if exact_match: return 'eq(text, $searchString)'

        # too small word fails with 'regular expression is too wide-ranging and can't be executed efficiently'
        use_regexp = is_word(str) and len(str) >= 3

        if use_regexp:
            params['$regexp'] = f"/{str}.*/i"

        regexp = "regexp(text, $regexp)" if use_regexp else ''
        anyoftext = "anyoftext(text, $searchString)"
        exprs = [s for s in [anyoftext, regexp] if len(s) > 0]
        if len(exprs) > 1:
            s = ' or '.join(exprs)
            return f"({s})"
        return exprs[0]

    range = f"offset: {offset}, first: {limit}"
    term_range = f", {range}" if is_term_list else ''

    brace = lambda s: f"({s})"
    search_filter = make_search_filter()
    lang_filter = f'eq(lang, "{lang}")' if lang else ''
    tag_filter = brace(' or '.join(
        f"uid_in(tag, {t['uid']})"
        for t in tags)) if not is_empty(tags) else ''

    filter_expr = ' and '.join([
        f for f in
        [has_term_type, has_tag_type, lang_filter, tag_filter, search_filter]
        if len(f) > 0
    ])
    term_filter = f"@filter({filter_expr})" if is_term_list else ''

    args = ", ".join([f"{k}: string" for k in keys(params)])
    param_query = f"query terms({args}) " if args else ''

    file_edges = ['audio', 'visual']

    def make_edge(name):
        is_file = name in file_edges
        myrange = f"({range})" if kind == name else "(first: 10)"
        body = FILE_BODY if is_file else TERM_BODY
        return f"{name} {myrange} {body}"

    all_edge_keys = list(keys(relation_map)) + file_edges
    edges = '\n'.join([make_edge(k) for k in all_edge_keys])

    def make_total(pred, name=''):
        if not name:
            name = f"{pred}_count"
        return f"{name}: count({pred})"

    totals = [make_total(k) for k in all_edge_keys] if is_term else []
    if not is_term:
        totals.insert(0, make_total('uid' if is_term_list else kind, 'total'))

    totals = '\n'.join(totals)

    count_query = """
    count(func: {match_fn}) {term_filter} {{
      {totals}
    }}
    """.format(match_fn=match_fn, term_filter=term_filter, totals=totals)

    text = """{param_query}{{
  terms(func: {match_fn}{term_range}) {term_filter} {{
    Tag
    uid
    text
    lang
    {meta}
    {tag}
    {edges}
  }}
  {count}
}}""".format(param_query=param_query,
             match_fn=match_fn,
             term_range=term_range,
             term_filter=term_filter,
             meta=META,
             tag='' if no_links else TAG,
             edges='' if no_links else edges,
             count='' if no_links else count_query)

    p = re.compile(r"^\s*[\r\n]", flags=re.MULTILINE)
    text = p.sub('', text)

    return {'text': text, 'params': params}
