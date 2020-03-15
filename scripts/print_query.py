#!/usr/bin/env python

import os
import api
from utils import dump_json


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    q = """query doc() {
        doc(func: has(Document)) {
            uid
            tag
            title
        }
    }
    """
    resp = api.query(q)
    print(dump_json(resp))
    return


if __name__ == '__main__':
    main()
