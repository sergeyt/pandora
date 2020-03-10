FROM docker.elastic.co/elasticsearch/elasticsearch:6.8.7

COPY --chown=elasticsearch:elasticsearch config/elasticsearch/jvm.options /usr/share/elasticsearch/config/
COPY --chown=elasticsearch:elasticsearch config/elasticsearch/elasticsearch.yml /usr/share/elasticsearch/config/

EXPOSE 9200 9300
