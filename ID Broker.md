# ID Broker Quickstart

### Environment Variables

For working docker compose runs, you need to export some environment variables.

```bash
$ # In order to use Docker and Compose with buildkit enabled, export two environment variables for your current shell
$ export DOCKER_BUILDKIT=1 
$ export COMPOSE_DOCKER_CLI_BUILD=1

$ # in order to run containers as the currently logged in user, export his user and group ids
$ export UID=$(id -u) 
$ export GID=$(id -g)
```


## Starting ZITADEL
You can connect to [ZITADEL on localhost:4200](http://localhost:4200) after the frontend compiled successfully. Initially it takes several minutes to start all containers.

```bash
$ docker compose -f ./build/local/docker-compose-local.yml --profile backend --profile frontend up --detach
```

### Initial login credentials

**username**: `zitadel-global-org-admin@caos.ch`
**password**: `Password1!`  