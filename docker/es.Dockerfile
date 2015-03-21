FROM elasticsearch:latest
RUN /usr/share/elasticsearch/bin/plugin install mobz/elasticsearch-head
RUN /usr/share/elasticsearch/bin/plugin install royrusso/elasticsearch-HQ
RUN /usr/share/elasticsearch/bin/plugin install lmenezes/elasticsearch-kopf
RUN /usr/share/elasticsearch/bin/plugin install karmi/elasticsearch-paramedic
