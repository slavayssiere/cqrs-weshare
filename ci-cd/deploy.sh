#!/bin/bash

ABSPATH=$(cd "$(dirname "$0")"; pwd)
echo $ABSPATH

cd $ABSPATH/../iac
docker-compose up -d
cd -

while ! nc -z localhost 8001; do   
  echo "wait for Kong started..."
  sleep 1
done

echo "Kong started !"
sleep 10

curl -i -X POST \
  --url http://localhost:8001/services/ \
  --data 'name=read-cqrs' \
  --data 'url=http://loadbalancer:8081'

curl -i -X POST \
  --url http://localhost:8001/services/ \
  --data 'name=write-cqrs' \
  --data 'url=http://loadbalancer:8082'

curl -i -X POST \
  --url http://localhost:8001/services/write-cqrs/routes \
  --data 'hosts[]=cqrs.com&methods[]=POST&methods[]=PUT'

curl -i -X POST \
  --url http://localhost:8001/services/read-cqrs/routes \
  --data 'hosts[]=cqrs.com&methods[]=GET'

write_route_id=$(curl -X GET http://localhost:8001/services/write-cqrs/routes | jq -r .data[0].id)
read_route_id=$(curl -X GET http://localhost:8001/services/read-cqrs/routes | jq -r .data[0].id)

curl -X POST http://localhost:8001/services/read-cqrs/plugins \
    --data "name=zipkin"  \
    --data "config.http_endpoint=http://tracing:9411/api/v2/spans" \
    --data "config.sample_ratio=1"

curl -X POST http://localhost:8001/services/write-cqrs/plugins \
    --data "name=zipkin"  \
    --data "config.http_endpoint=http://tracing:9411/api/v2/spans" \
    --data "config.sample_ratio=1"

curl -X POST http://localhost:8001/routes/$write_route_id/plugins \
    --data "name=zipkin"  \
    --data "config.http_endpoint=http://tracing:9411/api/v2/spans" \
    --data "config.sample_ratio=1"

curl -X POST http://localhost:8001/routes/$read_route_id/plugins \
    --data "name=zipkin"  \
    --data "config.http_endpoint=http://tracing:9411/api/v2/spans" \
    --data "config.sample_ratio=1"
