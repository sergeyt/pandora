#!/usr/bin/env python

import os
import api


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    print(api.access_token)


if __name__ == '__main__':
    main()
