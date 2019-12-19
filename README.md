- [Introduction](#introduction)
- [Deliverables](#deliverables)
  - [HTTP Web Server](#http-web-server)
  - [tokengen](#tokengen)
  - [HTTP Server Service](#http-server-service)
- [Project Details](#project-details)
  - [JWT Authentication](#jwt-authentication)
  - [OpenAPI Interface](#openapi-interface)
  - [Layout](#layout)
  - [Building](#building)
  - [Dependencies](#dependencies)
  - [Testing](#testing)
- [References](#references)

# Introduction

This represent a learning project that explores:

* Creating a HTTP Server Service 
  * Using middleware wrappers around enpoint to add new functionality
* Using JWT Authentication to protect endpoints
* Using OpenAPI Interface to document the interfaces

# Deliverables

Three components are delivered as part of this project

* HTTP Web Server (cmd/server)
* JWT Token Generator (cmd/tokengen)
* HTTP Server Service (go module)

## HTTP Web Server

The sample web server includes the following endpoints. Some enpoints are protected using Bear Authorization using JTW Authentication.
See discussion under [JTW Authentication](#jw-authentication)

* [http://localhost:8080](http://localhost:8080) (OpenAPI Client interface)
* [http://localhost:8080/unprotectedAPI](http://localhost:8080/unprotectedAPI)
* [http://localhost:8080/protectedAPI](http://localhost:8080/protectedAPI) (Requires Bear JWT Token Authentication)
* [http://localhost:8080/admin](http://localhost:8080/admin)  (Requires Bear JWT Token Authentication with admin priv)

The web server includes the following options

```
$ cd cmd/server; ./server -h
Usage of ./server:
  -config string
    	The configuration file (default "Configuration.json")
  -port string
    	HTTP localhost service port to use (default ":8080")
```

**Note**: The configuration file used by tokengen and sever must be the same to ensure the same secret key is used for both. Otherwise JWT Authentication will fail.

```
$ cd cmd/server; ./server -config ../../Configuration.json
2019/12/17 12:59:00 Reading configuration file ../../Configuration.json
2019/12/17 12:59:00 Server started on localhost:8080
```

As the protected enpoints require authentication you will need to use the include OpenAPI static files to access these
via (http://localhost:8080)[http://localhost:8080] or your favorite testing tools like Postman or curl. 

**Note**: If you change the default port the static Open API file mechanism will not work.

The following illustrates how to perform a curl request 

```
$ export TOKEN="{your token generated from tokengen}"
curl -H 'Accept:application/json' -v -H "Authorization: Bearer ${TOKEN}" http://localhost:8080/protectedAPI
```

## Tokengen

Tokengen is a simple JWT Token tool that can be used to create and configure JWT Authorization token.
JWT Authentication relies on the use of a secret key to generate tokens. This secret key value is
stored in the _Configuration.json_ file. Before using this value should be changed to a unique secret
value. This value used here must be common between the tokengen cli tool and the web server. 

Tokengen supports the following options

```
$./tokengen -h
Usage of ./tokengen:
  -admin
    	Boolean flag to indicate whether the user is privilege admin user
  -config string
    	The configuration file (default "Configuration.json")
  -quite
    	No verbose output, only output token
  -time int
    	The amount of time in seconds before the token expires (default 1200)
  -username string
    	The username to use in the generation of the JWT authorization token (default "guest")
```

Example of creating an admin user token

```
$ cd cmd/tokengen; ./tokengen -config ../../Configuration.json -username janedoe -admin true
Reading ../../Configuration.json ...
Generating Token ...
JWT Authentication Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImphbmVkb2UiLCJpc19hZG1pbiI6dHJ1ZSwiaWF0IjoxNTc2NjE2MTE4LCJleHAiOjE1NzY2MTczMTh9.KAZ4xkwhLtODcURggKuEOzJK-SULK8OXtcgkmNwhdxI
```

## HTTP Server Service

The HTTP Server Service included in this code base follows the basic service design discussed by Mat Ryer in his Gophercon 2019 talk, [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831). 

See go docs for information on usage of this service

```
$ godoc -http=localhost:6060
```

# Project Details

## JWT Authentication

This project uses JWT Authentication is used to restrict access the endpoints 

* /protectedAPI
* /admin

The cli command tokengen is used to generate a Bear Authorizatio token. The username, admin privileges (true/false) information which is embedded in the JWT Claim object can be changed via options to the tokengen command.

This command generated an admin user token that will grant access to the endpoints /protectedAPI and /admin
```
$ cd cmd/tokengen; ./tokengen -config ../../Configuration.json -username admin -admin true
Reading ../../Configuration.json ...
Generating Token ...
JWT Authentication Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiaXNfYWRtaW4iOnRydWUsImlhdCI6MTU3NjYxNjE0MSwiZXhwIjoxNTc2NjE3MzQxfQ.XkuGWDFkufEDGlITMAenQ4wn83etO1TExX4y7ZnZDsw
```

This command generates a token that will grant access to the endpoint /protectedAPI. Note that since the user is not a admin privileged  they will not have access to the /admin enpoint.
```
$ cd cmd/tokengen; ./tokengen -config ../../Configuration.json -username averagejoe
Reading ../../Configuration.json ...
Generating Token ...
JWT Authentication Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImF2ZXJhZ2Vqb2UiLCJpc19hZG1pbiI6ZmFsc2UsImlhdCI6MTU3NjYxNjE5MiwiZXhwIjoxNTc2NjE3MzkyfQ.x8XLjenNF4jF6tKempqE7PZUruM-Mgopf0WE4092Nlk
```

The JWT Authorization Token in this project uses a custom Claims object that contains some addition fields data fields.

Using a [JWT decoder](https://jwt.io/) you can view the data embedded in the token. In our
first example above the JWT decoder will show

```
{
  "username": "admin",
  "is_admin": true,
  "iat": 1576616141,
  "exp": 1576617341
}
```

In our second it will show

```
{
  "username": "averagejoe",
  "is_admin": false,
  "iat": 1576616192,
  "exp": 1576617392
}
```

These username and is_admin fields match the options given to tokengen when the tokens were created.

**NOTE**: When using JWT Token in real life treat them as passwords. They are valid until they expire.


## OpenAPI Interface

The JSON endpoints exposed by this project can be viewed using the OpenAPI files included in _api.
In addition the server will serve up static html files generated from these files that can
be used to interact with the included server.

```
http://localhost:8080
```

For more information on OpenAPI API see [What is OpenAPI?](https://swagger.io/docs/specification/about/)

## Layout

This project conforms to the following layout

```
.
├── Configuration.json      - Configuration file, JWT Secret Key
├── Makefile                - Makefile for building project
├── README.md               - Project README
├── _api                    - OpenAPI files
│   ├── Makefile            - Makefile to regenerate OpenAPI static files
│   └── openapi.yaml        - OpenAPI 3.0 design
├── claims.go               - Helper file for JWT Authentication
├── cmd
│   ├── server
|   |   |── index.html      - OpenAPI interface served by server
│   │   ├── main.go         - Sample web server
│   └── tokengen
│       ├── main.go         - JWT Authentication Server
├── config.go               - Reads and processes configuration file
├── go.mod
├── go.sum
├── routes.go               - Defines the routes supported by Server Service
└── server.go               - Implements Server Service
```
## Building 

This project can be built by the included Makefile.  The following makefile options
are available

- all - performs the operations: clean, build, and test
- build - build the code base
- test - test the code base. See Testing.
- coverage - test coverage report
- clean - perform a clean operation
- run - start the included server

```
$ make
```

## Dependencies

The OpenAPI portion of this project that generates the static files has a dependency on
openapi. This can be installed via

```npm install @openapitools/openapi-generator-cli -g```

To regenerate the static files

```
$ cd _api
$ make
```

## Testing

All tests in this project can be run either via 

```
$ make test
```
or
```
$ go test
```

# References

* [Implementing JWT based authentication in Golang](https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/)
* [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)