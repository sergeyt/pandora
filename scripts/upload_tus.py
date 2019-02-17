#!/usr/bin/env python

import os
from tusclient import client
import utils

utils.enable_logging_with_headers()

dir = os.path.dirname(os.path.realpath(__file__))
filename = os.path.realpath(os.path.join(dir, '../schema.txt'))

my_client = client.TusClient('http://localhost:1080/api/files')
# my_client.set_headers(
#     {'Upload-Metadata': 'filename {0}'.format(os.path.basename(fname))})

uploader = my_client.uploader(filename, chunk_size=1024)
uploader.upload()
