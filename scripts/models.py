from utils import is_empty

class Term:
    def __init__(self, text, lang=None, region=None):
        if is_empty(text):
            raise Exception('text is not deined')
        self.text = text.strip()
        self.lang = lang
        self.region = region

class File:
    def __init__(self, url, region = None):
        if is_empty(url):
            raise Exception('url is not defined')
        self.url = url

class TermWithData:
    def __init__(self, term, data):
        if term is None:
            raise Exception('term is not defined')
        if data is None or len(data) == 0:
            raise Exception('data is not defined')
        self.term = term
        self.data = data
