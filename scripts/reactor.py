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
        print("Received a message on '{subject} {reply}': {data}".format(
            subject=subject, reply=reply, data=data))
        if data['type'] == 'file':
            index_file.delay(data['url'])

    await nc.subscribe("*", cb=message_handler)

    await nc.drain()


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    loop.close()
