#!/bin/sh

docker-compose exec provider \
  /go/bin/creamy-home-auth-provider add username password
