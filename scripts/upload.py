#!/usr/bin/env python

import os
from tusclient import client

dir = os.path.dirname(os.path.realpath(__file__))

tc = client.TusClient('http://localhost:4200/api/files/')
uploader = tc.uploader(os.path.join(dir, '../schema.txt'), chunk_size=1024)
uploader.upload()
