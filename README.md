# FIFA World Cup Winners

This project exposes a Web API for accessing historic data from
the FIFA World Cup championship.

## Running tests

A proper Go environment is required in order to run this project.
Once setup, tests can be run with the `make test` or simply `make` command.

## Running the server

Once all tests are passing, the server can be started with
the `make start` command.

## Testing the API manually

Start the server with `make start` and then
read the _Access Token_ printed to standard output.
This token will be used for POST requests.

### GETting

`curl -i http://localhost:8000/`  
`curl -i http://localhost:8000/winners`  
`curl -i http://localhost:8000/winners?year=1970`  
`curl -i http://localhost:8000/winners?year=banana`

### POSTin with no access token

```
curl -i -X POST \
-d '{"country":"Croatia", "year": 2030}' http://localhost:8000/winners
```

### POSTing with valid access token

First, start the sever and read the value for the Access Token printed
to standard output.

```
curl -i -X POST \
-H "X-ACCESS-TOKEN: 5577006791947779410" \
-d '{"country":"Croatia", "year": 2030}' http://localhost:8000/winners
```

Then check for the newly added winner

`curl -i http://localhost:8000/winners`

### POSTing with invalid data

```
curl -i -X POST \
-H "X-ACCESS-TOKEN: 5577006791947779410" \
-d '{"country":"Russia", "year": 1984}' http://localhost:8000/winners
```

### POSTing with invalid method

`curl -i -X PUT -d '{"country":"Russia", "year": 2030}' http://localhost:8000/winners`

### Running with Docker

To build the image from the Dockerfile, run:

`docker build -t project-fifa-world-cup .`

To start an interactive shell, run:

`docker run -it --rm --name run-fifa project-fifa-world-cup`

From inside the shell, run the tests with:

`go test handlers/*`
