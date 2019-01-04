# CQRS

L'objectif de ce travail est de tester la mise en place d'une architecture CQRS.
Les pré-requis pour le lancement en local sont Docker et Docker Compose.

## Lancement

Pour lancer l'infrastructure localement:

```language-bash
./ci-cd/build-all.sh
./ci-cd/deploy.sh
```

### Test

Création des utilisateurs:

```language-bash
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/users -d '{"username":"slavayssiere", "email":"slavayssiere@wescale.fr", "address":"23 rue taitbout 75009", "age":32}'
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/users -d '{"username":"alexis", "email":"alexis@wescale.fr", "address":"23 rue taitbout 75009", "age":22}'
```

Création des topics de discussion

```language-bash
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/topics -d '{"topicname":"wespeakcloud"}'
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/topics -d '{"topicname":"perroquetGif"}'
```

Lancement d'une conversation

```language-bash
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/messages -d '{"userid":1, "topicid":2, "data":"no piaf here"}'
curl -i -H "Host:cqrs.com" -X POST http://localhost:8000/messages -d '{"userid":2, "topicid":2, "data":"some perroquet here"}'
```

Modification d'un utilisateur

```language-bash
curl -i -H "Host:cqrs.com" -X PUT http://localhost:8000/users/2 -d '{"age":33}'
curl -i -H "Host:cqrs.com" -X GET http://localhost:8000/users/2
```

Récupération d'une discussion

```language-bash
curl -i -H "Host:cqrs.com" -X GET http://localhost:8000/topics
curl -i -H "Host:cqrs.com" -X GET http://localhost:8000/topics/2
curl -H "Host:cqrs.com" -X GET http://localhost:8000/topics/2/complete | jq .
```

Modification d'un message et relecture de la conversation

```language-bash
curl -i -H "Host:cqrs.com" -X GET http://localhost:8000/messages/1
curl -i -H "Host:cqrs.com" -X PUT http://localhost:8000/messages/1 -d '{"data":"no perroquet here"}'
curl -H "Host:cqrs.com" -X GET http://localhost:8000/topics/2/complete | jq .
```

## Stop

Pour supprimer les conteneurs lancés localement.

```language-bash
cd iac
docker-compose down -v --remove-orphans
```

## Connect to

### MySQL

Se connecter à MySQL - le mot de passe est "my-secret-pw" par défaut.

```language-bash
mysqlsh --sql root@localhost:3306
```

### Redis

Se connecter à Redis.

```language-bash
redis-cli -h localhost
```
