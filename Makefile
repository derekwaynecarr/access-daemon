build-daemon:
	make -C ./daemon build

build-client:
	make -C ./client build

build: build-daemon build-client

test: build
	go test ./...

clean:
	make -C ./daemon clean

run: test
	./daemon/daemon --config-dir=./config

docker: test
	docker build . -t eparis/access-daemon:latest

docker-run: docker
	docker run --rm --privileged --pid=host --network=host --log-driver=none eparis/access-daemon:latest

docker-clean:
	./docker-clean.sh
