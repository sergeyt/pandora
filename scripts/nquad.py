import re


def is_uid(s):
    return len(s) > 0 and re.match(r"^0x[a-f0-9]+$", s) != None


def is_rdf_id(s):
    return len(s) > 0 and re.match(r"^_:([\w_]+)$", s) != None


def rdf_repr(v):
    if isinstance(v, str):
        if v == '*' or is_uid(v) or is_rdf_id(v):
            return v
        return "<{0}>".format(v) if is_uid(v) else '"{0}"'.format(v)
    return v


def format(id, k, v):
    a = k.split('@')
    p = a[0]
    lang = a[1] if len(a) == 2 else ''
    s = rdf_repr(v)
    if len(lang) > 0:
        s += "@{0}".format(lang)
    id = "<{0}>".format(id) if is_uid(id) else "_:{0}".format(id)
    return "{0} <{1}> {2} .".format(id, p, s)


# TODO refactor as generator
# converts dictionary to list of edges
def kv_list(d, id='x'):
    result = []
    for k, v in d.items():
        if type(v) is list:
            for t in v:
                result.append(format(id, k, t))
            continue
        result.append(format(id, k, v))
    return result


def format_edge(v):
    if isinstance(v, str):
        return v
    return format(v[0], v[1], v[2])


def format_edges(v):
    if isinstance(v, str):
        return v
    if isinstance(v, dict):
        return '\n'.join(kv_list(v))
    return '\n'.join([format_edge(t) for t in v])
