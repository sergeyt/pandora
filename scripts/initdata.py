#!/usr/bin/env python

import os
import api


def user_exists(user):
    try:
        api.login(user['login'], user['password'])
    except:
        return False
    return True


def ensure_user(user):
    if not user_exists(user):
        api.post('/api/data/user', user)


users = [
    {
        'login': 'sergeyt',
        'name': 'sergeyt',
        'email': 'stodyshev@gmail.com',
        'password': 'sergeyt123',
    },
]


def main():
    api.login("system", os.getenv("SYSTEM_PWD"))
    for user in users:
        ensure_user(user)


if __name__ == '__main__':
    main()
