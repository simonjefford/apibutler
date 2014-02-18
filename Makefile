default: build

build:
	./make.sh

run: build
	./ratelimit

container: build
	grunt
	docker build -t simonj/ratelimit .
