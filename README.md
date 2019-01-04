# CQRS

## Start

```language-bash
./ci-cd/build-all.sh
./ci-cd/deploy.sh
```

### Test

```language-bash
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/users -d '{"username":"slavayssiere", "email":"slavayssiere@wescale.fr", "address":"23 rue taitbout 75009", "age":32}'
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/users -d '{"username":"alexis", "email":"alexis@wescale.fr", "address":"23 rue taitbout 75009", "age":22}'

curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/topics -d '{"topicname":"wespeakcloud"}'
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/topics -d '{"topicname":"perroquetGif"}'

curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/messages -d '{"userid":1, "topicid":2, "data":"no piaf here"}'
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/messages -d '{"userid":2, "topicid":2, "data":"some perroquet here"}'

curl -i -H "Host:cqrs.com" -X PUT http://localhost:8000/users/2 -d '{"age":33}'

curl -i -H "Host:cqrs.com" -X GET http://localhost:8000/users/2
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
