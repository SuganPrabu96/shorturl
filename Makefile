VERSION="0.1.0"

# Download and make redis
.PHONY: download-and-make-redis
download-and-make-redis:
	mkdir -p dependencies
	wget http://download.redis.io/redis-stable.tar.gz
	tar xvzf redis-stable.tar.gz && mv redis-stable dependencies/redis && rm redis-stable.tar.gz
	make -C dependencies/redis && make -C dependencies/redis install

# Setup config for redis server
.PHONY: setup-redis-server
setup-redis-server:
	mkdir -p /etc/redis
	mkdir -p /var/redis
	cp dependencies/redis/utils/redis_init_script /etc/init.d/redis_6379
	cp dependencies/redis/redis.conf /etc/redis/6379.conf
	mkdir /var/redis/6379
	update-rc.d redis_6379 defaults

# Start redis server for dev
.PHONY: start-redis-server
start-redis-server:
	/etc/init.d/redis_6379 start

# Build
.PHONY: build
build:
	go build

# Build docker image
.PHONY: docker-build
docker-build:
	docker build -t shorturl:${VERSION} .

# Run docker image
.PHONY: run
run:
	docker run shorturl:${VERSION}

# Run docker image locally
.PHONY: run-local
run-local:
	docker run --env SERVER_HOST=localhost -p 5000:5000 shorturl:${VERSION}
