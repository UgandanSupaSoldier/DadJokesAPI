# DadJokesAPI

## Requirements
In order to build and run this project, you will require the following prerequisite programs:

1. Docker Compose
2. Go 1.20 or later

## Installation
This project is configured to run as a CLI tool. Included in the root directory of the project is a shell script which will build an executable binary to launch the CLI. Alternatively the main.go can be run directly. The server can be setup as follows:

1. Directly using Go:
```go run main.go setup```
2. Using the binary
```
    ./build.sh
    ./CLI_linux_<ARCH> setup
```

Firstly, this will create a config.toml file with default server config (please refer to the donfig section for details on the config.toml). If you wish to create a database with custom config such as a custom user, password or port, then you should create the config.toml before running setup.
Secondly, this will create and start a postgres docker container called dad-jokes-postgres and populate it's schema and data with the preset jokes.txt file.

*Note: This project was designed to run primarily on linux distributions. Although it is possible to compile the project into a windows compatible binary, this is currently untested.

## CLI TOOL
The CLI tool consists of three primary commands. These are namely setup, serve and test.
1. ```Setup``` is used to configure and install the server
2. ```Serve``` is used to launch the API server
3. ```Test``` is used to run the unit tests found within the tests folder

Example Usage:
```
    go run main.go Serve
    ./CLI_linux_amd64 Test
```
For more information on the CLI, you can run the CLI with no arguments or with the -h help flag.

## Config
The default config.toml file is structured as follows:
```
[server]
port = 8080
debug = false

[postgres]
debug = false
port = 12345
db = "dad-jokes-api"
user = "db-username"
password = "db-password"
max_idle_connections = 8
max_open_connections = 64

[rate_limiter]
rate = 1
burst = 10
timeout = 180
enabled = false
```

### Server Config
```port``` sets which port on localhost the API should listen on.
```debug``` toggles the log level between INFO and DEBUG.

### Postgres
```debug``` toggles database query debugging logs.
```port``` sets the port which will be exposed by the docker container.
```db``` sets the name of the database where the data will be stored.
```user``` sets the login name for the postgres user.
```password``` sets the login password for the postgres user.
```max_idle_connections``` sets the maximum number of idle connections
```max_open_connections``` sets the maximum number of open connections

### rate_limiter
```rate``` sets the maximum rate at which requests may be made to the server in requests per second.
```burst``` sets the maximum concurrent api calls allowed to pass, and may allow exceeding the rate.
```timeout``` sets the timeout for clearing old rate_limiter cache data such as IP addresses
```enabled``` toggles whether to enable or disable the rate limiter

## API

### Authorization
For protected endpoints, an authorization token will be required. This server uses a bearer token, which can be supplied in the requests headers as "Authorization: bearer <token>". The serve rstores all auth tokens generated in a table called ```tokens```.

### Responses
All responses are structured as follows:
```
{
    "data": <interface>,
    "error": map[string]interface
}
```
Where the data is the object returned, and the error contains any errors that occured in the form:
```
{
    "type": <string>,
    "message": <string>,
    "details": <interface>
}
```

### Resources
The API is made up of two resources, namely tokens and jokes.

Tokens are described with the following JSON struct, where ```expires_at``` is the tokens optional expiry date.
```
{
    "token": <string>,
    "expires_at": <*time>
}
```

Jokes are described as follows:
```
{
    "id": <int>
    "author": <*string>,
    "category":<*string>,
    "rating": <*float>,
    "text": <string>
}
```

#### GET /v1/token/new
A utility endpoint to generate new tokens for testing purposes. In production, a token management resource or server would be needed.
Accepts an optional query string parameter ```expiry``` to set a token expiry, and must contain a timestamp in the format ```YYYY-MM-DDThh:mm:ss```.

Returns a newly generated token in the response's data section.

#### GET /v1/joke/random
Returns a random joke in the response's data section.

#### GET /v1/joke/search
Returns a paginated list of jokes, and accepts two query string parameters ```page``` and ```page_size```.
The ```page``` parameter sets the page offset from where to begin selecting.
The ```page_size``` paramter sets the size of each page returned, and affects the offset.

Returns the list of jokes in the response's data section.

#### POST /v1/joke/create
Used to create new jokes, and accepts a json body with the same format as described in the resources secton above. Note that the joke's text is required, and the maximum value of the rating is 10. This endpoint is protected and requires a valid auth token to access.

The created joke resource is returned in the response's data section.

### Misc

#### Logging
All requests made to the server as well as their responses are recorded to the datase in a table called ```logs```. This is handled by a logging middleware in ```logging.go```, and can be modified to record any information needed about api calls.

#### Panics
All panics on the server are handled in the middleware in a file called ```panic_handler.go```. For external logging or incident logging, this is where additional functionality can be injected.

#### Uninstallation
Currently there is no uninstallation scripts or support. As such, in order to remove this service completely, the docker image and volume associated with this program will need to be removed manually.

#### Missing Features and Improvements
These are some features or improvements that, given more time, would have been included:

1. Including the token used in the request logs.
2. Improving the setup script to take user input.
3. Better handling of the setup scripts for already installed components.
4. More detailed and secure token generation and management.
5. Response caching.
6. More verbose logging throughout the system for better visibility.
7. Using an interface for the logging and database components, to allow for better flexibility.
8. Splitting of view, system and database models.
9. Support for timezone in token expiry parameter
