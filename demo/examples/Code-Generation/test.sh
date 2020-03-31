#!/bin/bash

# Run your application
go run main.go

# Test the server with curl or open the browser to http://localhost:8080

# The response should be {"Content":"Hello World!"}
curl localhost:8080
