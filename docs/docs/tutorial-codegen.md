---
id: tutorial-codegen
title: Codegen Tutorial
sidebar_label: Codegen Tutorial
---

import useBaseUrl from '@docusaurus/useBaseUrl';

:::info
For more detailed information about code generation, refer to [Code Generation](./gen-code.md)
:::

In this tutorial, we're going to generate code for a simple web application that returns a random pet.

## Start

Create a new folder to build our project in. Now create a `specs` folder to put our specification files. All Sysl code generation projects start from a specification file.

## Applications

Create a new file called `petdemo.sysl` in `specs` with the following content:

```sysl
# backend specifications
import petstore.yaml as petstore.Petstore

Petdemo "Petdemo":
    @package="Petdemo"
    /pet:
        GET:
            | Get a random pet
            Petstore <- GET /pet
            return ok <: Pet

    !type Pet:
        breed <: string


```

Create another file called `petstore.yaml` with the following content:

```yaml
openapi: "3.0.3"
info:
  version: 1.0.0
  title: Pet Service
  license:
    name: MIT
servers:
  - url: https://australia-southeast1-innate-rite-238510.cloudfunctions.net/pet-demo
paths:
  /pet:
    get:
      summary: Get random pet
      operationId: getPet
      tags:
        - pet
      responses:
        "200":
          description: A random pet
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Pet"
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"
components:
  schemas:
    Pet:
      type: string
    Error:
      type: object
      required:
        - code
        - message
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string
```

## Generate

First, run the bootstrap script using the command below

```bash
docker run -it --rm -v `pwd`:/work anzbank/sysl-go:latest specs/petdemo.sysl github.com/anz-bank/sysl-go-demo Petdemo:Petdemo
```

You should see the initial files generated including `Makefile`, `go.mod` and `codegen.mk`

Now run make to generate the skeleton code for your application.

```bash
make
```

Don't forget to run `go mod tidy` to generate the `go.sum` file

```bash
go mod tidy
```

This can now be committed if you're using a version control system.

## Handler

Before you can successfully run your application, you need to add some handler functions for the endpoints you've defined in your specification.

To do this, create the file in the path `internal\handlers\pet.go`

```go
package handlers

import (
	"context"
  "net/http"

	petdemo "github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo"
	"github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo/petstore"
	"github.com/anz-bank/sysl-go/common"
)

// GetRandomPetPicListRead reads random pic from downstream
func GetRandomPetPicListRead(ctx context.Context,
	getRandomPetPicListRequest *petdemo.GetPetListRequest,
	client petdemo.GetPetListClient) (*petdemo.Pet, error) {

	// Set response content type to JSON
	headers := common.RequestHeaderFromContext(ctx)
	headers.Set("Content-Type", "application/json; charset=utf-8")

	reqPetstore := petstore.GetPetListRequest{}
	pet, err := client.PetstoreGetPetList(ctx, &reqPetstore)
	if err != nil {
		return nil, err
	}

	// return the result
	return &petdemo.Pet{
		Breed: string(*pet),
	}, nil
}


```

Now we need to register the handler

Add the following line to L21 of `cmd/Petdemo/main.go`

```go
GetPetList: handlers.GetRandomPetPicListRead,
```

Great, now try and run it with

```
go run cmd/Petdemo/main.go
```

You'll notice that we get the following error message:

```bash
configuration is empty
```

Now we need to add our config file!

## Config

Add the following config file in `config/config.yaml`

```yaml
library:
  log:
    format: text
    level: info
    caller: false

genCode:
  upstream:
    contextTimeout: 120s
    http:
      basePath: /
      readTimeout: 120s
      writeTimeout: 120s
      common:
        hostName: ""
        port: 6060
  downstream:
    contextTimeout: 120s
    petstore:
      serviceURL: https://australia-southeast1-innate-rite-238510.cloudfunctions.net/pet-demo
      clientTimeout: 59s
```

Great, now you can run the application, passing in the config file

```bash
go run cmd/Petdemo/main.go config/config.yaml
```

Now open your browser to `http://localhost:6060/pet` and you should see

Note: The request may timeout the first time since the downstream is a Google Cloud Function. Try a few times to get a successful response

## Multiple Calls

Now let's add another API.

This integrates the Pokemon API, a test API that doesn't require authentication which can be found at [https://pokeapi.co/](https://pokeapi.co/):

```yaml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: PokeAPI Service
  license:
    name: MIT
servers:
  - url: https://pokeapi.co/api/v2
paths:
  /pokemon/{id}:
    get:
      summary: Get a pokemon
      operationId: getPokemon
      tags:
        - pokemon
      parameters:
        - name: id
          description: The identifier for this resource.
          in: path
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: A pokemon
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Pokemon"
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Error"
components:
  schemas:
    Pokemon:
      type: object
      properties:
        id:
          type: integer
          description: The identifier for this resource.
        name:
          type: string
          description: The name for this resource.
        height:
          type: integer
          description: The height of this Pokémon in decimetres.
    Error:
      type: object
      required:
        - code
        - message
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string
```

Add the following `import` statement to the top of your Sysl file:

```sysl
import pokeapi.yaml as pokeapi.PokeAPI
```

And the following call statement:

```sysl
PokeAPI <- GET /pokemon/{id}
```

This declares that we are calling the GET pokemon endpoint on the PokeAPI service.

```sysl
# import downstream specifications
import petstore.yaml as petstore.Petstore
import pokeapi.yaml as pokeapi.PokeAPI


Petdemo "Petdemo":
    @package="Petdemo"
    /pet:
        GET:
            | Get a random pet
            Petstore <- GET /pet
            PokeAPI <- GET /pokemon/{id}
            return ok <: Pet

    !type Pet:
        breed <: string

```

Now run `make` and some new PokeAPI client code is generated in the `gen` folder.

Add the following YAML block to your `config.yaml` file:

```yaml
pokeapi:
  serviceURL: https://pokeapi.co/api/v2
  clientTimeout: 59s
```

Your resultant `config.yaml` file should look like this:

```yaml
library:
  log:
    format: text
    level: info
    caller: false

genCode:
  upstream:
    contextTimeout: 120s
    http:
      basePath: /
      readTimeout: 120s
      writeTimeout: 120s
      common:
        hostName: ""
        port: 6060
  downstream:
    contextTimeout: 120s
    petstore:
      serviceURL: https://australia-southeast1-innate-rite-238510.cloudfunctions.net/pet-demo
      clientTimeout: 59s
    pokeapi:
      serviceURL: https://pokeapi.co/api/v2
      clientTimeout: 59s
```

Great, now we'll call the PokeAPI service in our handler code.

Let's define a new request type called `reqPokemon`.

If you type `pokeapi.` your IDE should offer `GetPokemonRequest` as an autocomplete option. Let's use the response of the random pet we received as the `ID` parameter.

`reqPokemon := pokeapi.GetPokemonRequest{ID: int64(len(*pet))}`

We can make a call to the endpoint by typing `client.` and we'll see the autocomplete option `PokeapiGetPokemon`. We simply need to pass in the `context` and `GetPokemonRequest` object we defined previously.

We'll be appending the result to the `Breed` field for now, but we'll fix this up later.

We also need to set the `Accept-Encoding` header as sysl-go currently only supports `gzip` encoding. Put this before your `reqPokemon` definition.

```
	headers = common.RequestHeaderFromContext(ctx)
	headers.Set("Accept-Encoding", "gzip")
```

Your `pet.go` file should now look like this:

```go
package handlers

import (
	"context"
    "net/http"

	petdemo "github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo"
	"github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo/petstore"
	"github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo/pokeapi"
	"github.com/anz-bank/sysl-go/common"
)

// GetRandomPetPicListRead reads random pic from downstream
func GetRandomPetPicListRead(ctx context.Context,
	getRandomPetPicListRequest *petdemo.GetPetListRequest,
	client petdemo.GetPetListClient) (*petdemo.Pet, error) {

	// Set response content type to JSON
	headers := common.RequestHeaderFromContext(ctx)
	headers.Set("Content-Type", "application/json; charset=utf-8")

	reqPetstore := petstore.GetPetListRequest{}
	pet, err := client.PetstoreGetPetList(ctx, &reqPetstore)
	if err != nil {
		return nil, err
	}

	// Set response encoding type to gzip
	// This is required as only gzip encoding is currently supported by sysl-codegen
	headers = common.RequestHeaderFromContext(ctx)
	headers.Set("Accept-Encoding", "gzip")

	reqPokemon := pokeapi.GetPokemonRequest{ID: int64(len(*pet))}
	pokemon, err := client.PokeapiGetPokemon(ctx, &reqPokemon)
	if err != nil {
		return nil, err
	}

	// return the result
	return &petdemo.Pet{
		Breed: string(*pet) + " " + *pokemon.Name,
	}, nil
}
```

Great, now run `go run cmd/Petdemo/main.go config/config.yaml` and navigate your browser to `http://localhost:6060/pet`.

You should see the name of a pokemon appended to the random pet.
e.g

```json
{ "breed": "guinea pig caterpie" }
```

Let's put this in its own field now.
Add the `pokemon` field to the response type.

Your Sysl file should now look like the following:

```sysl
# import downstream specifications
import petstore.yaml as petstore.Petstore
import pokeapi.yaml as pokeapi.PokeAPI


Petdemo "Petdemo":
    @package="Petdemo"
    /pet:
        GET:
            | Get a random pet
            Petstore <- GET /pet
            PokeAPI <- GET /pokemon/{id}
            return ok <: Pet

    !type Pet:
        breed <: string
        pokemon <: string
```

Now run `make` and modify the return statement in `pet.go` so that we return the Pokemon in a separate field.

```go
package handlers

import (
	"context"
    "net/http"

	petdemo "github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo"
	"github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo/petstore"
	"github.com/anz-bank/sysl-go-demo/gen/pkg/servers/Petdemo/pokeapi"
	"github.com/anz-bank/sysl-go/common"
)

// GetRandomPetPicListRead reads random pic from downstream
func GetRandomPetPicListRead(ctx context.Context,
	getRandomPetPicListRequest *petdemo.GetPetListRequest,
	client petdemo.GetPetListClient) (*petdemo.Pet, error) {

	// Set response content type to JSON
	headers := common.RequestHeaderFromContext(ctx)
	headers.Set("Content-Type", "application/json; charset=utf-8")

	reqPetstore := petstore.GetPetListRequest{}
	pet, err := client.PetstoreGetPetList(ctx, &reqPetstore)
	if err != nil {
		return nil, err
	}

	// Set response encoding type to gzip
	// This is required as only gzip encoding is currently supported by sysl-codegen
	headers = common.RequestHeaderFromContext(ctx)
	headers.Set("Accept-Encoding", "gzip")

	reqPokemon := pokeapi.GetPokemonRequest{ID: int64(len(*pet))}
	pokemon, err := client.PokeapiGetPokemon(ctx, &reqPokemon)
	if err != nil {
		return nil, err
	}

	// return the result
	return &petdemo.Pet{
		Breed:   string(*pet),
		Pokemon: *pokemon.Name,
	}, nil
}

```

Great, now run `go run cmd/Petdemo/main.go config/config.yaml` and navigate your browser to `http://localhost:6060/pet`
You should now see something similar to the following response

```json
{ "breed": "guinea pig", "pokemon": "caterpie" }
```

## Conditional Logic

Let's add some conditional logic to our specification file to specify that we only call the PokeAPI when the length of the response is less than 50.

Your `petdemo.sysl` file should look like the following:

```
# import downstream specifications
import petstore.yaml as petstore.Petstore
import pokeapi.yaml as pokeapi.PokeAPI


Petdemo "Petdemo":
    @package="Petdemo"
    /pet:
        GET:
            | Get a random pet
            Petstore <- GET /pet
            if len(Pet) < 50:
                PokeAPI <- GET /pokemon/{id}
            return ok <: Pet

    !type Pet:
        breed <: string
        pokemon <: string


```

Run `make` and notice how the generated code hasn't changed. This is because Sysl doesn't generate business logic, this must be done in the handler code written by you.

However, the conditional logic in the Sysl file does show up in the diagrams generated for you.

Run the command below and open your browser to `localhost:6900`. We're interested in the `Petdemo` service so navigate to there and expand the sequence diagram for the `GET /pet` endpoint.

```bash
  docker run -p 6900:6900 -v $(pwd):/usr:ro anzbank/sysl-catalog:v1.4.199 run --serve -v usr/specs/petdemo.sysl
```

You can see here that the sequence diagram shows the conditional logic defined in your Sysl file. These automatically-generated diagrams help you to easily understand the behavior of your application.

## Next steps

TODO: We'll be adding more instructions soon to use the response we received to fetch an image of a random pet.
