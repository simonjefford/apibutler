default: buildbin

buildbin:
	./make.sh

run: buildbin npm bower
	grunt serve

npm:
	npm install

bower:
	bower install

buildweb: npm bower
	grunt

container: buildbin buildweb
	docker build -t simonj/ratelimit .
