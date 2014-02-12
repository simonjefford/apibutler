default: build

build:
	./make.sh

run: build
	./ratelimit
