#!/usr/bin/env python

import os
import utils
import api
from minio import Minio
import boto3
from botocore.client import Config

utils.enable_logging_with_headers()

filename = api.schema_path()

access_key = os.getenv('MINIO_ACCESS_KEY')
secret_key = os.getenv('MINIO_SECRET_KEY')
bucket_name = os.getenv('S3_BUCKET')

# mc = Minio(
#   # 'localhost:4200/api/s3',
#   'localhost:9000',
#   access_key=access_key,
#   secret_key=secret_key,
#   secure=False,
# )

# mc.fput_object(bucket_name, 'schema.txt', filename)

s3 = boto3.resource(
    's3',
    endpoint_url='http://localhost:9000',
    aws_access_key_id=access_key,
    aws_secret_access_key=secret_key,
    region_name='us-east-1',
    config=Config(signature_version='s3v4'),
)
s3.Bucket(bucket_name).upload_file(filename, 'schema.txt')
