build:
	@docker build --file es.Dockerfile --tag es .

rebuild:
	@docker build --file es.Dockerfile --tag es --no-cache .

run: runfg

runfg:
	@docker run --publish 9200:9200 --name es01 --rm es elasticsearch -Des.node.name="TestNode01"

runbg:
	@docker run --detach --publish 9200:9200 --name es01 es elasticsearch -Des.node.name="TestNode01"

stop:
	@docker stop es01
