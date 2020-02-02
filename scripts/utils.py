import os
import sys
import re
import json
import urllib
import logging
import requests
from json import JSONEncoder
from langdetect import detect


def as_bool(s):
    return len(s) > 0 and (s == '1' or s.lower() == 'true')


TESTING = as_bool(os.getenv('TESTING', '0'))
CHROME_USER_AGENT = 'Mozilla/5.0 AppleWebKit/537.36 (KHTML like Gecko) Chrome/78.0.3883.121'

# https://stackoverflow.com/questions/16337511/log-all-requests-from-the-python-requests-module
# https://stackoverflow.com/questions/10588644/how-can-i-see-the-entire-http-request-thats-being-sent-by-my-python-application


def enable_logging(with_headers=False):
    if with_headers:
        enable_logging_with_headers()
    else:
        logging.basicConfig(level=logging.DEBUG)


def enable_logging_with_headers():
    # These two lines enable debugging at httplib level (requests->urllib3->http.client)
    # You will see the REQUEST, including HEADERS and DATA, and RESPONSE with HEADERS but without DATA.
    # The only thing missing will be the response.body which is not logged.
    try:
        import http.client as http_client
    except ImportError:
        # Python 2
        import httplib as http_client
    http_client.HTTPConnection.debuglevel = 1

    # You must initialize logging, otherwise you'll not see debug output.
    logging.basicConfig()
    logging.getLogger().setLevel(logging.DEBUG)
    requests_log = logging.getLogger("requests.packages.urllib3")
    requests_log.setLevel(logging.DEBUG)
    requests_log.propagate = True


def find_audio_args(argv=sys.argv):
    text = argv[1]
    lang = argv[2] if len(argv) >= 3 else detect(text)
    if lang != 'ru':
        lang = 'en'
    return (text, lang)


def url_exists(url):
    if not url:
        return False
    try:
        headers = {
            'User-Agent': 'script',
        }
        resp = requests.head(url, headers=headers)
        if resp.ok:
            return True
        with requests.get(url, headers=headers, stream=True) as resp:
            if resp.ok:
                return True
    except:
        pass
    print(f'not found: {url}')
    with open("notfound.txt", 'a') as f:
        f.write(url + '\n')
    return False


def url_quote(val):
    if isinstance(val, str):
        return urllib.parse.quote(val)
    return [urllib.parse.quote(s) for s in val]


def is_word(s):
    return s and re.match(r"^[^\s]+$", s) != None


def is_empty(val):
    return val is None or len(val.strip()) == 0


class JSONEncoderEx(JSONEncoder):
    def default(self, o):
        return o.__dict__


def dump_json(d):
    return json.dumps(d,
                      cls=JSONEncoderEx,
                      sort_keys=True,
                      indent='  ',
                      ensure_ascii=False)
