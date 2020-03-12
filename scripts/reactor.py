import os
import asyncio
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
        data = msg.data.decode()
        print("received a message on '{subject} {reply}': {data}".format(
            subject=subject, reply=reply, data=data))

        if 'resource_type' in data and 'url' in data and data[
                'resource_type'] == 'file':
            url = data['url']
            index_file.delay(url)

    await nc.subscribe("global", cb=message_handler)


if __name__ == '__main__':
    print('reactor started')
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    # todo graceful shutdown
    loop.run_forever()
    print('reactor exited')
