# CQRS

## Start

```language-bash
cd iac
docker-compose up -d
```

```language-bash
docker run --rm \
    -e "KONG_DATABASE=postgres" \
    -e "KONG_PG_HOST=kong-database" \
    -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
    --network iac_default \
    kong kong migrations up
docker run -d --name kong \
    --link kong-database:kong-database \
    --link write:write \
    --link read:read \
    -e "KONG_DATABASE=postgres" \
    -e "KONG_PG_HOST=kong-database" \
    -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
    -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
    -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
    -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
    -e "KONG_ADMIN_LISTEN=0.0.0.0:8001" \
    -e "KONG_ADMIN_LISTEN_SSL=0.0.0.0:8444" \
    -p 8000:8000 \
    -p 8443:8443 \
    -p 8001:8001 \
    -p 8444:8444 \
    --network iac_default \
    kong
```

### Configure Kong

```language-bash
curl -i -X POST \
  --url http://localhost:8001/services/ \
  --data 'name=read-cqrs' \
  --data 'url=http://read:8080'

curl -i -X POST \
  --url http://localhost:8001/services/ \
  --data 'name=write-cqrs' \
  --data 'url=http://write:8080'

curl -i -X POST \
  --url http://localhost:8001/services/write-cqrs/routes \
  --data 'hosts[]=cqrs.com&methods[]=POST&methods[]=PUT'

curl -i -X POST \
  --url http://localhost:8001/services/read-cqrs/routes \
  --data 'hosts[]=cqrs.com&methods[]=GET'
```

### Test

```language-bash
curl -H "Host:cqrs.com" -X POST http://localhost:8000/users -d '{"username":"slavayssiere", "email":"slavayssiere@wescale.fr", "address":"23 rue taitbout 75009", "age":32}'

curl -H "Host:cqrs.com" -X GET http://localhost:8000/users
```

## Stop

```language-bash
cd iac
docker-compose down -v --remove-orphans
```

## Connect to

### MySQL

```language-bash
mysqlsh --sql root@localhost:3306
```

### Redis

```language-bash
redis-cli -h localhost
```
