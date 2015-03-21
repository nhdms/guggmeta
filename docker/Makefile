build:
	@docker build --file es.Dockerfile --tag es .

rebuild:
	@docker build --file es.Dockerfile --tag es --no-cache .

run:
	@docker run --publish 9200:9200 --name es01 --rm es elasticsearch --env="ES_MIN_MEM=320m,ES_MAXM_MEM=640m" -Des.node.name="TestNode01" -Des.node.master=true -Des.node.data=true -Des.node.local=true -Des.discovery.zen.ping.multicast.enabled=false -Des.index.number_of_shards=1 -Des.index.number_of_replicas=0

stop:
	@docker stop es01
