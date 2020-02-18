#!/bin/bash

# Make sure that GOPRIVATE is set so go can get private repos
export GOPRIVATE="github.service.anz"

go run main.go

curl localhost:8080/foobar
#{"completed":false,"id":1,"title":"delectus aut autem","userId":1}
