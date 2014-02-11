FROM ubuntu
RUN apt-get update
RUN apt-get install -y redis-server
ADD ratelimit /usr/local/bin/ratelimit
ADD dockerstart /usr/local/bin/dockerstart
ENTRYPOINT /usr/local/bin/dockerstart
