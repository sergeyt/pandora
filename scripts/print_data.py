#!/usr/bin/env python

import os
import api
from utils import dump_json


def dump(q):
    resp = api.query(q)
    print(dump_json(resp))


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    dump("""query tag() {
        tag(func: has(Tag)) {
            uid
            text
        }
    }
    """)
    dump("""query doc() {
        doc(func: has(Document)) {
            uid
            tag {
                uid
                text
            }
            title
        }
    }
    """)


if __name__ == '__main__':
    main()
