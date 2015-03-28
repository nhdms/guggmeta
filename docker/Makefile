default: build run

build:
	@docker build --file Dockerfile.elasticsearch --tag ggm-elasticsearch .
	@docker build --file Dockerfile.redis --tag ggm-redis .

rebuild:
	@docker build --file Dockerfile.elasticsearch --no-cache --tag ggm-elasticsearch .
	@docker build --file Dockerfile.redis --no-cache --tag ggm-redis .

run:
	-@docker stop es01 redis01
	-@docker rm es01 redis01
	@docker run --detach --publish 9200:9200 --name es01 \
		ggm-elasticsearch \
		elasticsearch --env="ES_MIN_MEM=320m,ES_MAXM_MEM=640m" -Des.node.name="TestNode01" -Des.node.master=true -Des.node.data=true -Des.node.local=true -Des.discovery.zen.ping.multicast.enabled=false -Des.index.number_of_shards=1 -Des.index.number_of_replicas=0
	@docker run --detach --publish 6379:6379 --name redis01 \
		ggm-redis \
		redis-server

stop:
	@docker stop es01 redis01
