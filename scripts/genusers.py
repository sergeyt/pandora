#!/usr/bin/env python

import os
import utils
import api
from faker import Faker
from initdata import ensure_user

utils.enable_logging()


def generate():
    fake = Faker()
    for i in range(100):
        name = fake.name()
        user = {
            'login': name,
            'name': name,
            'email': name + '123@gmail.com',
            'password': name + '123',
        }
        ensure_user(user)


if __name__ == '__main__':
    generate()
