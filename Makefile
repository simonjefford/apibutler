GOPATH := $(CURDIR)/_vendor:$(GOPATH)

default: buildbin

buildbin:
	go build -o apibutler

run: buildbin npm bower
	grunt serve

test:
	go test fourth.com/apibutler/...

bench:
	cd apiproxyserver; go test -bench . 2> /dev/null

npm:
	npm install

bower:
	bower install

buildweb: npm bower
	grunt

container: buildbin buildweb
	docker build -t simonj/apibutler .
