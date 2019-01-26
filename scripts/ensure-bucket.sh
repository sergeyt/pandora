# export the vars in .env into your shell:
export $(egrep -v '^#' ../.env | xargs)

MINIO_NAME=localhost

while ! nc -z ${MINIO_NAME} 9000; do echo 'wait minio...' && sleep 0.1; done; \
  sleep 5 && \
  s3cmd mb s3://${MINIO_BUCKET}

# docker-compose run minio --entrypoint sh minio/mc -c "\
#   while ! nc -z minio 9000; do echo 'wait minio...' && sleep 0.1; done; \
#   sleep 5 && \
#   mc config host add myminio http://minio:9000 \$MINIO_ENV_MINIO_ACCESS_KEY \$MINIO_ENV_MINIO_SECRET_KEY && \
#   mc rm -r --force myminio/\$MINIO_BUCKET || true && \
#   mc mb myminio/\$MINIO_BUCKET && \
#   mc policy download myminio/\$MINIO_BUCKET \
# "
