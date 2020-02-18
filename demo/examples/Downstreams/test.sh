#!/bin/bash

# Make sure that GOPRIVATE is set so go can get private repos
export GOPRIVATE="github.service.anz"

go run main.go

# {"Content":"Hello World!"}
curl localhost:8080/stuff
