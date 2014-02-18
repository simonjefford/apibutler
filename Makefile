default: build

build:
	./make.sh

run: build
	./ratelimit

web:
	npm install
	grunt

container: build web
	docker build -t simonj/ratelimit .
