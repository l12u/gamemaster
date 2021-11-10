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

# Data structure

When communicating with the gamemaster service, there are two data types that are especially important: the *board* and the *game*. A board specifies one specific type of game and looks
like this:

```json
{
  "type": "hangman",
  "url": "board.l12u.party/hangman"
}
```

The url references the board service(s) that are responsible for this board. A game for this specific board, could look like this:

```json
{
  "id": "729e4265-877d-4549-8064-e99d77f7295e",
  "type": "hangman",
  "players": [
    {
      "id": "4e6683b8-d7ac-4e6c-9682-88ba025e5286",
      "name": "Peter",
      "role": "host"
    }
  ],
  "state": "lobby",
  "createdAt": "1636535217809",
  "updatedAt": "1636535217809"
}
```