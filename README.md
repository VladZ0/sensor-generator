I left my config.yml and .env varibales to simplify your app configuration process.
Anyway you can change app configs and .env as you want.

Start app using Make commands, all are descripted below.

If you want to test app and start it in docker, just use: make test ---> make build ---> make up.
But dont forget to add your configs (.env.postgres, .env.redis, config.yml)

make swagger --->
    This command will create swagger docs if you change it.

make test --->
    Test app.

make build --->
    Build docker-compose.

make up --->
    Run docker-compose. (with redis and postgres)

make down --->
    Stop containers. (if you started it with 'make up' command)

make migrate --->
    If you want to use postgres locally (not in container), use this command to create tables.
    It uses default host and port.
    But don't forget to create your database and set configs.

make start_redis --->
    Start redis.

make stop_redis --->
    Stop redis container. 
