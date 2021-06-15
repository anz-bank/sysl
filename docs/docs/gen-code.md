---
id: gen-code
title: Code Generation
keywords:
  - code
  - go
  - grpc
  - rest
  - generation
---

import useBaseUrl from '@docusaurus/useBaseUrl';

<!-- :::caution
WIP

Resources:

- https://github.com/anz-bank/sysl-go
- https://github.com/anz-bank/protoc-gen-sysl
- https://github.com/anz-bank/sysl-template

**TODO:**

- Update and polish content.
- Move referenced assets to a permanent directory on GitHub and update links.

::: -->

---

## Introduction

Sysl can generate server and client code for specified applications. Currently, Go is the only language supported, but we have several on the roadmap including Swift, JavaScript, Java and Kotlin.

:::info If you have a language you'd like to see supported, please [raise an issue](https://github.com/anz-bank/sysl/issues/new?labels=enhancement&template=feature_request.md), so we can prioritise it.
:::

<!-- To get started with generating code with sysl, refer to [code-generation](/docs/byexample/code-generation/) -->

## Quick Start

You first need to write a `.sysl` file, then you can quickly bootstrap the application with a single command, which sets up everything automatically for you. To create a project, run:

```sh
docker run -it --rm -v `pwd`:/work anzbank/sysl-go:latest [SYSL_FILE] [GO_MOD_NAME] [APP:PKG...]
```

For example:

```sh
docker run -it --rm -v `pwd`:/work anzbank/sysl-go:latest project.sysl github.com/anz-bank/test-project MyApp1:myapp1 MyApp2:myapp2
```

If any of the arguments are not supplied, the CLI will provide prompts to gather the necessary information to generate an application.

After the bootstrapper is finished, to generate the server code and handle dependencies, run the following command:

```sh
make && go mod tidy
```

You can read the documentation for the config file needed by running the following command:

```sh
go run ./cmd/your_app_name/main.go --help
```

Assuming you have a config file `config.yml`, you can quickly run the following application using the following command:

```sh
go run ./cmd/your_app_name/main.go config.yml
```

## Folder Structure

Running the commands will generate the initial project structure:

```
my-app
├── .github/
│   └── workflows/
│       └── test.yml
├── cmd/
│   └── ...
├── gen/
│   └── pkg/
│       └── servers/
│           └── ...
├── .gitattributes
├── .gitignore
├── Dockerfile
├── README.md
├── Makefile
├── codegen.mk
├── go.mod
└── go.sum
```

- `.github/workflows/test.yml`: Contains a generic GitHub workflow: build, lint, test and coverage.
- `Dockerfile`: Builds the app and creates an image containing the binary.
- `README.md`: A README outlining the app, to which details can be added.
- `cmd/{apps}/main.go`: The entry point to the app. Provides the config and request handlers to the framework.
- `.gitignore`: Ignores relevant files.
- `.gitattributes`: Flags generated files so GitHub hides diffs by default.
- `gen/pkg/servers/{apps}/{files}`: The generated server and client code for each application.

## Features

### Architecture

A Sysl model can generate code for REST and gRPC services as follows:

<img alt="sysl go architecture diagram" src={useBaseUrl('img/sysl-go-architecture-diagram.png')}/>

The generated code makes use of the common code in the [sysl-go](https://github.com/anz-bank/sysl-go) library, and a `main` function that constructs the structs defined in the generated code and starts listening for requests.

- `Service` describes the API of a service. `Client` provides the means to call the API of a service. In that sense a `Client` implements the API of its `Service`.
- `ServiceHandler` wraps actual serving of a request. It parses data from the request, validates the inputs, constructs and configures a client to make any necessary downstream calls, then passes the client to the `ServiceInterface` to perform the business logic, including any downstream requests. Once the `ServiceInterface`'s work is complete, the `ServiceHandler` validates the return value and sends a response.

<img alt="Sysl go service diagram" src={useBaseUrl('img/sysl-go-service-diagram.png')}/>

#### Rest Endpoint Naming

When generating the names for the handlers in the service interface for a REST service it will base it off the HTTP method and the full path of the endpoint. For the enpoint `GET /foo` in the example above, the name generated is `GetFoo`.

Generally it will leave out any parameters that are in the path from the name. So the endpoint below would generate the name `GetGreeting`:

```
HelloService:
    /greeting/{userId <: int}:
           GET:
               ...
```

If you have 2 endpoints that generate to the same name, for example if you had another endpoint that was just `/greeting`, then you can add the `~vars_in_url_name` annotation to specify that you want the parameters to be include in the name. So the endpoint below would generate the name `GetGreetingUserid`:

```
HelloService:
    /greeting/{userId <: int} [~vars_in_url_name]:
           GET:
               ...
```

### Authorization

The `@authorization_rule` annotation can be used in the sysl specification of a service to permit or deny access to a method or endpoint.

By default, if no `@authorization_rule` annotation is present, then all requests to the method or endpoint are accepted.

Authorization rules are supported when:

- a code-generated application serves gRPC methods

Authorization rules are not yet supported when:

- a code-generated application serves REST endpoints

The following facts may be used in authentication rule logic:

- scopes in the claims of a verified JSON Web Token (JWT).

Example of an authorization rule:

```
GreetingService [package="greeting", ~gRPC]:

    Greet(Request <: GreetRequest):
        @authorization_rule = "any(jwtHasScope('hello'))"

        return ok <: GreetResponse
```

`@authorization_rule` is an authorization expression like `any(jwtHasScope('hello'))`.

An application server generated by `sysl-go` from this specification will include an authorization rule to check that any call to the `Greet` method has a JWT with verified claims containing the "hello" scope.

#### Run-time configuration of authorization logic

A code-generated application that uses the `@authorization_rule` annotation must be configured at runtime to define how to verify JSON Web Tokens (JWTs).

A list of trusted JSON Web Token issuers must be defined in the application config file, under the config key path `library.authentication.jwtauth.issuers`.

For example:

```
library:
  authentication:
    jwtauth:
      issuers:
        - name: "an-example-issuer"
          jwksURL: "https://example.com/.well-known/jwks.json"
          cacheTTL: 10m
```

The list of issuers can contain multiple entries if you need to configure your application to trust JWTs issued by different issuers. Tokens signed by any of the configured issuers are regarded as verified by any endpoint or method tagged with an `@authorization_rule` annotation.

Explanation of parameters:

- `name`: this parameter is used by [sysl-go/jwtauth](https://github.com/anz-bank/sysl-go/tree/master/jwtauth) library to determine which configured issuer will be used to verify a given JWT. The `name` parameter MUST be equal to the value of the "iss" (Issuer) Claim inside a JWT that should be verified using this issuer. `sysl-go/jwtauth` does not support verification of JWTs that lack an "iss" (Issuer) claim. The name parameter is case sensitive. This claim is defined by [RFC 7519 section 4.1.1](https://tools.ietf.org/html/rfc7519#section-4.1.1).
- `jkwsUrl`: the URL of a resource containing a [JSON Web Key Set (JKWS)](https://tools.ietf.org/html/rfc7517) object. This JWKS object encoding the public keys that will be trusted when the application validates a JWT. This URL will be accessed by the application using the HTTP GET method.
- `cacheTTL`: mandatory duration parameter giving the time-to-live (TTL) of a cached JWKS object. Must be a nonzero duration.
- `refreshCache`: optional duration parameter giving a period between refreshes of the JKWS cache. If set to a positive duration, the cached JWKS will be periodically refreshed using a new JWKS object obtained from the given `jkwsUrl`

For more details, please refer to the documentation of the [sysl-go/jwtauth](https://github.com/anz-bank/sysl-go/tree/master/jwtauth) library.

#### Disabling authorization logic in a development environment

All authorization logic can be disabled by setting the following configuration option in application configuration:

```
development:
  disableAllAuthorizationRules: true
```

Beware: disabling all authorization rules is insecure.

#### Authorization expressions:

##### Authorization expression grammar:

```
<expr>      := <op-expr> | <atom>
<op-expr>   := <ident> "(" <expr> ("," <expr>)* ")"
<atom>      := <ident> "(" [ <literal> ("," <literal>)* ] ")"
<literal>   := <string-literal>
```

String literals can be single or double quoted.

##### Authorization expression evaluation:

Each `op-expr` represents an operator that transforms one or more input argument boolean values into an output boolean value.

Each `atom` represents a fact, parametrized by 0 or more string literal parameters. It evaluates to a boolean value.

If an expression evaluates to true, this is interpreted as allowing access.

If an expression evaluates to false, this is interpreted as denying access.

##### Defined operators:

- `all`: conjunction (logical and) of the arguments. Takes one or more arguments.
- `any`: disjunction (logical or) of the arguments. Takes one or more arguments.
- `not`: negation (logical not) of the arguments. Takes exactly one argument.

##### Defined facts:

- `jwtHasScope`: is the gRPC method in question is presented with a JWT containing verified claims including the named scope. Takes exactly one string literal argument (the scope name).

##### Example authorization expression:

```
all(any(jwtHasScope("fizz"),jwtHasScope("buzz")),not(jwtHasScope("test")))
```

|         | Example claims                                                                  |
| ------- | ------------------------------------------------------------------------------- |
| Valid   | `{"scope": "fizz"}`<br/>`{"scope": "banana fizz"}`<br/>`{"scope": "buzz"}`      |
| Invalid | `{"scope": "banana"}`<br/>`{"scope": "test"}`<br/>`{"scope": "fizz buzz test"}` |

Note: the behaviour of `jwtHasScope` fact evaluation can be customised, the above example shows the default behaviour of how JWT claims for scopes are encoded [1](https://www.iana.org/assignments/jwt/jwt.xhtml) [2](https://tools.ietf.org/html/rfc8693).

### Context

#### Clock

You can use [anz/pkg/clock](https://github.com/anz-bank/pkg/tree/master/clock), which is a context-driven wrapper for the time library. It allows substitution of mock clocks via context.Context for testing and other purposes.

You can also use `TimeTravel` which is a mock Clock that offsets time by a given offset. Also, when a call is made that involves a time delay, it travels to that time and returns instantly.

#### Logging

The generated project uses [github.com/anz-bank/pkg/log](https://github.com/anz-bank/pkg/tree/master/log) as the logger. This library makes the logger available via context. It is designed to make development and testing easier by allowing loggers to be defined per context instead of a single global logger or passing a logger down through the call stack.

##### Initialise the logger

sysl-go always uses a pkg logger internally. If custom code passes in a [Logrus](https://github.com/sirupsen/logrus) logger (a mechanism which is deprecated), then a hook is added to the internal pkg logger that forwards logged events to the provided logger.

sysl-go can be requested to log in a verbose manner, including additional details within log events where appropriate. The mechanism to set this verbose manner is to either have a sufficiently high Logrus log level or the verbose mode set against the pkg logger.

### Configuration

YAML files configure the generated code, library and admin server. The config file can only contain valid Sysl-go configuration properties. If a config file contains properties unknown to Sysl-go then the application will terminate on launch. To see the complete set of valid Sysl-go configuration properties, view the help information for your application:

```
$ ./app --help
Configuration file YAML schema:
...
```

#### Generated code

##### Upstream

The `upstream` section drives the configuration of the generated upstream server. The `contextTimeout` serves as the access timeout of the server. Further configuration is split into `http` and/or a `grpc` sections which describe the whereabouts and access timeouts of the respective servers.

```
genCode:
  upstream:
    contextTimeout: 120s
    http:
      basePath: /
      readTimeout: 120s
      writeTimeout: 120s
      common:
        hostName: ""
        port: 8080
    grpc:
      hostName: ""
      port: 8081
```

##### Downstream

The `downstream` section drives the configuration of the generated downstream clients. The `contextTimeout` serves as the access timeout of all the downstream dependencies.

Further configuration is consolidated under a downstream client specific section. The example below has only one downstream called `flickr`.

- `serviceURL` specifies the coordinates of the downstream service which is incidentally a REST service here.
- `clientTimeout` gives an opportunity to further customize the access timeout for downstream `flickr`.

The `clientTLS` section is used to drive the TLS configuration of chosen downstream `flickr`. `minVersion` and `maxVersion` specify the version of TLS being used. `serverIdentity`.`certKeyPair` specify the directory where the certificates are stored.

```
genCode:
  downstream:
    contextTimeout: 120s
    flickr:
      serviceURL: http://localhost:6060
      clientTimeout: 59s
        clientTransport:
        useProxy: true
        clientTLS:
            insecureSkipVerify: false
            minVersion: "1.2"
            maxVersion: "1.3"
            selfSigned: true
            serverIdentity:
            certKeyPair:
                certPath: "etc/cert/cert"
                KeyPath: "etc/cert/key"
```

#### Library

The configuration parameters for all the libraries used in the generated application are consolidated under the `library` section. For example the following attributes are used to drive the configuration of the logging library.

```
library:
  log:
    format: text
    level: 5
    caller: false
```

To assist in tracing calls through the log, every request will automatically be assigned a `traceid` which will be added to all log messages associated with that request. This ID will either be copied from the incoming `RequestID` header (if it exists) otherwise it will be a randomly generated UUID.
If the ID you would like to use is in a different header, you can override the default location using the config:
```
library:
  trace:
    incomingHeaderForID: headerToUse
```

#### Admin Server

The attributes `hostName`, `port` and `basePath` represent the whereabouts of the admin server whereas the attributes `readTimeout` and `writeTimeout` represent the access timeouts for the read and the write operations.

```
custom:
  adminServer:
    common:
      hostName: ""
      port: 3332
    basePath: /admintest
    readTimeout: 1s
    writeTimeout: 3s
```

### Environment Variables

Environment Variables matching entries in the config YAML file will be read from the environment. By default, Sysl-go reads configuration values from a configuration file. However, in certain circumstances it is desirable to override some or all of these
configuration values with environment variables. To support this, Sysl-go allows for a prefix to be set so that configuration values will first be read from the environment variables before falling back to the configuration file. The attribute `envPrefix` determines the prefix to use when reading environment variables.

For example, given a configuration file with the following contents:

```
envPrefix: ENV
genCode:
    upstream:
        contextTimeout: 120s
```

If an environment variable named `ENV_GENCODE_UPSTREAM_CONTEXTTIMEOUT` exists then it will override the value found in the configuration file.
