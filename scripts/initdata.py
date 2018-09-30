import requests
import json
import utils
import api

utils.enable_logging()

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
        'login': 'admin',
        'name': 'admin',
        'email': 'stodyshev@gmail.com',
        'password': 'admin123',
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

# api.drop_all()
# api.init_schema()
init()
