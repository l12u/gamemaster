# gamemaster

The gamemaster service is the main part when dealing with games and everything around it. It stores and manages all currently running games and the players that are associated with
them.

To be able to run multiple gamemaster services at the same time, it is highly recommended to use the Redis storage provider option.

# Usage

To start the gamemaster service, especially for local development, there is a dedicated script `run.sh` to use for that. If the Docker image is not built yet locally, the script will
automatically do that for you. If you want to force the rebuild regardless, you can execute:

```shell
sh run.sh --rebuild
```

This will of course only run the gamemaster, which means that you have to provide a Redis instance for yourself - or just use the local storage provider. If you're feeling lazy, you
could also just use:

```shell
docker-compose up -d 
```

# Configuration

A list of all available environment variables and their usage:

| Variable             | Default                       | Description |
| -------------------- | ----------------------------- | ----------- |
| ENABLE_REDIS_STORAGE | `false`                       | Enables the Redis backend storage provider. If `false`, then the local storage provider will be used, which means, that the data is lost when the service goes down.            |
| REDIS_ADDRESS        | `localhost:6379`              | Address of the Redis instance, if the Redis storage provider is enabled.            |
| REDIS_PASSWORD       | ` `                           | Password when authenticating to the Redis instance. Default for the authentication is an empty string meaning no authentication at all.            |
| REDIS_DATABASE       | `0`                           | Id of the database to use. Normally this setting doesn't need to be changed, but sometimes it can be useful, when keys are potentially overlapping.        |
| BOARD_CONFIG         | `/etc/gamemaster/boards.json` | Path to the config, where all the boards are configured.            |