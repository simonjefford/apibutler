#!/usr/bin/env bash
set -e

apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	apt-utils \
	automake \
	build-essential \
	curl \
	git \
        tree \
	--no-install-recommends

curl -s https://go.googlecode.com/files/go1.2.src.tar.gz | tar -v -C /usr/local -xz

cd /usr/local/go/src && ./make.bash --no-clean 2>&1
echo 'export PATH=/usr/local/go/bin:$PATH' >> /home/vagrant/.profile
