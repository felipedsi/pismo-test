# Pismo Test
This application runs an API to handle financial transactions.

### Requirements
- [Go](https://go.dev/)

Or you can just use Docker:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/compose-file/)

### Running with docker
First, you need to build the Docker container:
```bash
./script/build
```

This command is going to create an image called `pismo-test:dev` which can be used to run the tests and the application itself.

Then, you can just run using the docker-compose:
```bash
docker-compose up
```

The application will be acessible through `http://localhost:3000` endpoint.

### Running without docker
First, you need to have a postgres instance running and set the `POSTGRESQL_URL` environment variable with the connection info:
```bash
export POSTGRESQL_URL=postgres://<username>:<password>@<host>:<port>/<databaseName>
```

Then, just run the start script from the root folder:
```bash
./script/start
```

### Testing
You can run the tests with docker by running:
```bash
./script/test
```

Otherwise, you can just use the regular go command to run them:
```bash
go test ./...
```

### Migrations
This project uses the [golang-migration](https://github.com/golang-migrate/migrate) tool to track changes to the database schema.

To create a new migration, just run the following command and update the SQL up and down files accordingly:
```bash
migrate create -ext sql -dir db/migrations -seq <migration_name>
```

### Acknowledgments
Thanks Pismo and Leonardo for the opportunity to do this challenge! :)
