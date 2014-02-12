FROM ubuntu
RUN apt-get update
RUN apt-get install -y redis-server
ADD ratelimit /srv/ratelimit
ADD dockerstart /srv/dockerstart
ADD public /srv/public
ENTRYPOINT /srv/dockerstart
