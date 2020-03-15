import os
import asyncio
import json
import traceback
import logging
import server_reloader
from urllib.parse import urlparse, urlunparse, ParseResult
from nats.aio.client import Client as NATS
from nats.aio.errors import ErrConnectionClosed, ErrTimeout, ErrNoServers
from tasks import index_file


# removes query from given url
def clean_url(url):
    a = urlparse(url)
    b = ParseResult(scheme=a.scheme,
                    netloc=a.netloc,
                    path=a.path,
                    params='',
                    query='',
                    fragment='')
    return urlunparse(b)


async def reactor(loop):
    nc = NATS()

    NATS_URI = os.getenv('NATS_URI', 'nats://localhost:4222')
    await nc.connect(NATS_URI, loop=loop)

    async def message_handler(msg):
        try:
            data = json.loads(msg.data.decode())
            method = data.get('method', '')
            resource_type = data.get('resource_type', '')

            print(
                f"msg chan='{msg.subject}', type={resource_type}, data={data}")

            if resource_type == 'file':
                url = clean_url(data.get('url'))
                print(f'processing file {url}')
                index_file.delay(url)
        except:
            traceback.print_exc()

    await nc.subscribe("global", cb=message_handler)


def run():
    print('reactor started')
    logging.basicConfig(level=logging.DEBUG)
    loop = asyncio.new_event_loop()
    loop.set_debug(True)
    loop.run_until_complete(reactor(loop))
    # todo graceful shutdown
    loop.run_forever()
    print('reactor exited')


def main():
    server_reloader.main(run, before_reload=lambda: print('Reloading code…'))


if __name__ == '__main__':
    main()
