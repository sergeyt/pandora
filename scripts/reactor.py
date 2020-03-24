import os
import asyncio
import json
import traceback
import logging
import server_reloader
from nats.aio.client import Client as NATS
from nats.aio.errors import ErrConnectionClosed, ErrTimeout, ErrNoServers
from tasks import index_file
from worker import app


async def reactor(loop):
    nc = NATS()
    NATS_URI = os.getenv('NATS_URI', 'nats://localhost:4222')
    await nc.connect(NATS_URI, loop=loop)
    await nc.subscribe("global", cb=message_handler)


async def message_handler(msg):
    try:
        data = json.loads(msg.data.decode())
        method = data.get('method', '')
        resource_type = data.get('resource_type', '')

        # print(
        #     f"msg chan='{msg.subject}', type={resource_type}, data={data}")

        if resource_type == 'file':
            # use file id to cancel by same url
            file_id = data.get('resource_id')
            url = f'/api/file/{file_id}'
            if method == 'DELETE':
                print(f'cancelling tasks for file: {url}')
                cancel_tasks(url)
            else:
                print(f'processing file {url}')
                index_file.delay(url)
    except:
        traceback.print_exc()


# cancels active and scheduled related tasks
def cancel_tasks(url):
    state = app.control.inspect()

    def is_related(t):
        return t['name'] == 'tasks.index_file' and t['args'][0] == url

    cancel = []
    for _, a in state.active().items():
        cancel.extend([t for t in a if is_related(t)])

    for _, a in state.scheduled().items():
        cancel.extend([t for t in a if is_related(t['request'])])

    for t in cancel:
        print(f"cancelling task {t['id']} {t['name']}")
        app.control.revoke(t['id'], terminate=True)

    return len(cancel) > 0


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
