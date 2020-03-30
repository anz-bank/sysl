## Simple REST endpoint generation in go


api/: contains all API specifications for the generated application

gen/: contains all the generated code for the service

[internal/server/server.go](./internal/server/server.go): The hand-written code that's written; Server config and such

pkg/defaultcallback: contains code that sets up the defaults for generated code. (This will no longer be necessary in future Sysl versions)

When new endpoints are added, they need to be added to the `simple.ServiceInterface` variable in [server.go](./server/server.go)

[main.go](./main.go): runs the actual server


run `make` to regenerate application code
First you need to edit the start of the Makefile:

```
input = your input sysl file
app = < the app you want to develop>
down = <downstreams in a list separated by spaces>
basepath = <Your current go module path>
```

so: `make input=model/simple.sysl app=Simple` for this example

run `go run main.go` to start the server
