FROM elasticsearch:latest
RUN /usr/share/elasticsearch/bin/plugin install elasticsearch/elasticsearch-mapper-attachments/2.4.3
RUN /usr/share/elasticsearch/bin/plugin install mobz/elasticsearch-head
