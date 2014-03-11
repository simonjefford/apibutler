export GOPATH = $(shell pwd)/.gopath:$(shell pwd)/vendor

default: buildbin

configureenv:
	rm -rf .gopath
	mkdir -p .gopath/src/fourth.com/
	ln -sf ../../.. .gopath/src/fourth.com/ratelimit

buildbin: configureenv
	go build -o ratelimit

run: buildbin npm bower
	grunt serve

test: configureenv
	cd apiproxyserver; go test

bench: configureenv
	cd apiproxyserver; go test -bench .

npm:
	npm install

bower:
	bower install

buildweb: npm bower
	grunt

container: buildbin buildweb
	docker build -t simonj/ratelimit .
