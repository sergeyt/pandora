#!/usr/bin/env python

from tusclient import client

tc = client.TusClient('http://localhost:4200/api/files/')
uploader = tc.uploader('./schema.txt', chunk_size=1024)
uploader.upload()
