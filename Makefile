default: build

build:
	./make.sh

run: build
	PORT=4000 ./ratelimit
