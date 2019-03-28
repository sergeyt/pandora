#!/usr/bin/env python

import os
import utils
import api
# from faker import Faker

utils.enable_logging()

SYSTEM_PWD = os.getenv('SYSTEM_PWD')
ADMIN_PWD = os.getenv('ADMIN_PWD')

if not SYSTEM_PWD:
    raise Exception('SYSTEM_PWD is not defined')
if not ADMIN_PWD:
    raise Exception('ADMIN_PWD is not defined')


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
        'login': 'system',
        'name': 'system',
        'email': 'stodyshev@gmail.com',
        'password': SYSTEM_PWD,
        'role': 'admin',
    },
    {
        'login': 'admin',
        'name': 'admin',
        'email': 'stodyshev@gmail.com',
        'password': ADMIN_PWD,
        'role': 'admin',
    },
    {
        'login': 'sergeyt',
        'name': 'sergeyt',
        'email': 'stodyshev@gmail.com',
        'password': 'sergeyt123',
    },
]


def init():
    for user in users:
        ensure_user(user)


# def generate():
#     fake = Faker()
#     for i in range(100):
#         name = fake.name()
#         user = {
#             'login': name,
#             'name': name,
#             'email': name + '123@gmail.com',
#             'password': name + '123',
#         }
#         ensure_user(user)

if __name__ == '__main__':
    init()
    # generate()
