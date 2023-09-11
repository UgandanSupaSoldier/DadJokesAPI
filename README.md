# DadJokesAPI

## Requirements
In order to build and run this project, you will require the following prerequisite programs:

* Docker Compose
* Go 1.20 or later

*<b>Note</b>: This project was designed to run primarily on Linux distributions. Although it is possible to compile the project into a Windows-compatible binary, this is currently untested.*

## Installation
This project is configured to run as a ```CLI``` tool. The CLI tool consists of three primary commands:
1. ```Setup``` is used to configure and install the server.
2. ```Serve``` is used to launch the API server.
3. ```Test``` is used to run the unit tests found within the tests folder.

Included in the root directory of the project is a shell script that will build an executable binary of the CLI. Alternatively the main.go can be run directly. As such, the server can be setup as follows:

1. Directly using Go:
```
    go run main.go setup
```
2. Using the binary
```
    ./build.sh
    ./CLI_linux_<ARCH> setup
```

This will first create a <b>config.toml</b> file with default server config (please refer to the config section for details on the config.toml) in the project root directory. It will then create and start a Postgres docker container called dad-jokes-postgres and populate its schema and data with the preset jokes.txt file.<

*<b>Note</b>: If you wish to create a database with custom a config, then you should create the config.toml before running setup.*<br>
*<b>Note</b>: For more information on the CLI commands, you can run the CLI with no arguments or with the -h help flag.*
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

*<b>Note</b>: The config.toml should always appear in the root of the project (after running setup)*
### Server Config
* ```port``` sets which port on localhost the API should listen on.
* ```debug``` toggles the log level between INFO and DEBUG.

### Postgres
* ```debug``` toggles database query debugging logs.
* ```port``` sets the port that will be exposed by the docker container.
* ```db``` sets the name of the database where the data will be stored.
* ```user``` sets the login name for the Postgres user.
* ```password``` sets the login password for the Postgres user.
* ```max_idle_connections``` sets the maximum number of idle connections
* ```max_open_connections``` sets the maximum number of open connections

### Rate Limiter
* ```rate``` sets the maximum rate at which requests may be made to the server in requests per second.
* ```burst``` sets the maximum concurrent API calls allowed to pass, and may allow exceeding the rate.
* ```timeout``` sets the timeout for clearing old rate_limiter cache data such as IP addresses
* ```enabled``` toggles whether to enable or disable the rate limiter

## API

### Authorization
For protected endpoints, an authorization token will be required. This token must be supplied using bearer auth, which can be specified in the requests headers as <b>"Authorization: bearer \<token\>"</b>. The server stores all auth tokens generated in a table called ```tokens```.

### Responses
All responses are structured as follows:
```
{
    "data": <interface>,
    "error": map[string]interface
}
```
Where the data is the object returned, and the error contains any errors that occurred in the form:
```
{
    "type": <string>,
    "message": <string>,
    "details": <interface>
}
```

### Resources
The API is made up of two resources, namely <b>tokens</b> and <b>jokes</b>.

Tokens are described with the following JSON struct, where ```expires_at``` is the token's optional expiry date.
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
Accepts an optional query string parameter ```expiry``` to set a token expiry, and must contain the timestamp in the format ```YYYY-MM-DDThh:mm:ss```.

Returns the newly generated token in the response's data section.

#### GET /v1/joke/random
Returns a random joke in the response's data section.

#### GET /v1/joke/search
Returns a paginated list of jokes, and accepts two query string parameters ```page``` and ```page_size```.
The ```page``` parameter sets the page offset from where to begin selecting jokes.
The ```page_size``` parameter sets the size of the page returned.

*<b>Note</b>: The offset is calculated as ```page * page_size```.*

Returns the list of jokes in the response's data section.

#### POST /v1/joke/create
Used to create new jokes, and accepts a JSON-encoded body with the <b>joke</b> in the same format as described in the resources section above. Note that the joke's text is required, and the maximum value of the rating is 10. This resource is protected and requires a valid auth token to access.

The created joke resource is returned in the response's data section.

### Misc

#### Logging
All requests made to the server as well as their responses are recorded to the database in a table called ```logs```. This is handled by a logging middleware in ```logging.go```, and can be modified to record any information needed about API calls.

#### Panics
All panics on the server are handled in the ```panic_handler.go``` middleware file. For external or incident logging and paging, this is where additional functionality can be injected.

#### Uninstallation
Currently, there is no uninstallation script or support. As such, in order to remove this service completely, the docker container and volume associated with this program will need to be removed manually. This can be done by navigating to the setup directory and running the following:
```
docker-compose down
docker volume rm setup_postgres-data

```

#### Missing Features and Improvements
These are some features or improvements that, given more time, would have been included:

* Including the token used in the request logs.
* Improving the setup script to take user input.
* Better handling of the setup scripts for already installed components.
* More detailed and secure token generation and management.
* Response caching.
* More verbose logging throughout the system for better visibility.
* Using an interface for the logging and database components, to allow for better flexibility.
* Splitting of view, system, and database models.
* Support for timezone in token expiry parameter
