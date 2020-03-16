import os

broker_url = os.getenv('CELERY_BROKER_URL', 'pyamqp://guest@localhost//')
result_backend = os.getenv('CELERY_BACKEND_URL', 'redis://localhost:6379')

task_serializer = 'json'
result_serializer = 'json'
accept_content = ['json']
