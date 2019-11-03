#!/usr/bin/env python3

import api
import initdata


def run():
    api.drop_all()
    api.init_schema()
    initdata.init()


if __name__ == '__main__':
    run()
