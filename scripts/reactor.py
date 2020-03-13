import os
import asyncio
import json
from nats.aio.client import Client as NATS
from nats.aio.errors import ErrConnectionClosed, ErrTimeout, ErrNoServers
from tasks import index_file


async def run(loop):
    nc = NATS()

    NATS_URI = os.getenv('NATS_URI', 'nats://localhost:4222')
    await nc.connect(NATS_URI, loop=loop)

    async def message_handler(msg):
        subject = msg.subject
        reply = msg.reply
        data = json.loads(msg.data.decode())

        print(f"received a message on '${subject} ${reply}': ${data}")

        resouce_type = data.get('resource_type', '')

        if resouce_type == 'file':
            url = data.get('url')
            print(f'processing file ${url}')
            index_file.delay(url)

    await nc.subscribe("global", cb=message_handler)


if __name__ == '__main__':
    print('reactor started')
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    # todo graceful shutdown
    loop.run_forever()
    print('reactor exited')
