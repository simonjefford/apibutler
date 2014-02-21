default: buildbin

buildbin:
	./make.sh

run: buildbin npm
	grunt serve

npm:
	npm install

buildweb: npm
	grunt

container: buildbin buildweb
	docker build -t simonj/ratelimit .
