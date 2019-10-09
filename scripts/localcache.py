import os
import json
import utils

TESTING = utils.TESTING
dir = os.path.dirname(os.path.realpath(__file__))


class Cache:
    def __init__(self, name, testing=None):
        self.name = name
        self.cache = None
        self.testing = TESTING

    def file_path(self):
        return os.path.join(dir, self.name + '.json')

    def get(self, word):
        if self.testing:
            return None
        if self.cache is None:
            self.cache = {}
            p = self.file_path()
            if os.path.isfile(p):
                with open(p, 'r') as f:
                    self.cache = json.load(f)
        return self.cache[word] if word in self.cache else None

    def put(self, word, record):
        if self.testing:
            return
        self.cache[word] = record
        with open(self.file_path(), 'w') as f:
            f.write(json.dumps(self.cache, sort_keys=True, indent='  '))
