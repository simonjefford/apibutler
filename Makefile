GOPATH := $(CURDIR)/_vendor:$(GOPATH)

default: buildbin

buildbin:
	go build -o apibutler

run: buildbin npm bower
	forego start

builddeps:
	go install fourth.com/apibutler

test:
	go test fourth.com/apibutler/...

testpublish:
	go test -v fourth.com/apibutler/... | go2xunit -output tests.xml
	ember test > embertest.tap

bench:
	cd apiproxyserver; go test -bench . 2> /dev/null

npm:
	npm install

bower:
	bower install

buildweb: npm bower
	ember build --environment=production

container: buildbin buildweb
	docker build -t simonj/apibutler .
