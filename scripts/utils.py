import sys
import logging
import requests
from langdetect import detect

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


def find_audio_args():
    text = sys.argv[1]
    lang = sys.argv[2] if len(sys.argv) >= 3 else detect(text)
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
    print('not found: {0}'.format(url))
    with open("notfound.txt", 'a') as f:
        f.write(url + '\n')
    return False
