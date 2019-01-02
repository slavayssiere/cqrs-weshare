# CQRS

## Start

```language-bash
cd iac
docker-compose up -d
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
curl -H "Host:cqrs.com" -X POST http://localhost:8000/users -d '{"username":"alexis", "email":"alexis@wescale.fr", "address":"23 rue taitbout 75009", "age":22}'

curl -H "Host:cqrs.com" -X POST http://localhost:8000/topics -d '{"topicname":"wespeakcloud"}'
curl -H "Host:cqrs.com" -X POST http://localhost:8000/topics -d '{"topicname":"perroquetGif"}'

curl -H "Host:cqrs.com" -X POST http://localhost:8000/messages -d '{"userid":1, "topicid":2, "data":"no piaf here"}'
curl -H "Host:cqrs.com" -X POST http://localhost:8000/messages -d '{"userid":2, "topicid":2, "data":"some perroquet here"}'

curl -H "Host:cqrs.com" -X PUT http://localhost:8000/users/2 -d '{"age":33}'

curl -H "Host:cqrs.com" -X GET http://localhost:8000/users/2
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
