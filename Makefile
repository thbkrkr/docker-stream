
build:
	doo b

push: build
	doo p

run:
	docker run -ti \
	  -e B=${B} \
	  -e U=${U} \
	  -e P=${P} \
	  -e T=${T} \
	  -v /var/run/docker.sock:/var/run/docker.sock \
	  krkr/docker-events-stream
