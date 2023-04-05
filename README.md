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
Just run the start script from the root folder:
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

### Acknowledgments
Thanks Pismo and Leonardo for the opportunity to do this challenge! :)
