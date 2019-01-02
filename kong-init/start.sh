#!/bin/sh

while ! nc -z $KONG_PG_HOST 5432; do   
  echo "wait for PG started..."
  sleep 1
done

echo "PG started !"

kong migrations up

echo "Kong migration done !"

/docker-entrypoint.sh kong docker-start

