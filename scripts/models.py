from collections import namedtuple

Term = namedtuple('Term', ['text', 'lang', 'region'], defaults=(None))
File = namedtuple('File', ['url', 'region'], defaults=(None))
