FROM ubuntu:saucy
RUN apt-get update
RUN apt-get install -y redis-server curl
RUN curl -O http://nodejs.org/dist/v0.10.25/node-v0.10.25-linux-x64.tar.gz
RUN tar -zxvf node-v0.10.25-linux-x64.tar.gz
RUN cp -R node-v0.10.25-linux-x64/* /usr/local
ADD redis.conf /etc/redis/redis.conf
ADD ratelimit /srv/ratelimit
ADD dockerstart /srv/dockerstart
ADD public/dist /srv/public
ADD testbackend /srv/testbackend
RUN cd /srv/testbackend && npm install
ENTRYPOINT /srv/dockerstart
