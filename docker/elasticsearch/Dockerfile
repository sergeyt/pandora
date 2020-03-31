FROM docker.elastic.co/elasticsearch/elasticsearch:7.6.2

COPY --chown=elasticsearch:elasticsearch config/elasticsearch/jvm.options /usr/share/elasticsearch/config/
COPY --chown=elasticsearch:elasticsearch config/elasticsearch/elasticsearch.yml /usr/share/elasticsearch/config/

# TODO configure plugins
# https://github.com/RD17/ambar/blob/master/ElasticSearch/Dockerfile
# RUN bin/elasticsearch-plugin install http://dl.bintray.com/content/imotov/elasticsearch-plugins/org/elasticsearch/elasticsearch-analysis-morphology/5.6.3/elasticsearch-analysis-morphology-5.6.3.zip
# RUN bin/elasticsearch-plugin install analysis-stempel
# RUN bin/elasticsearch-plugin install analysis-smartcn

EXPOSE 9200 9300
