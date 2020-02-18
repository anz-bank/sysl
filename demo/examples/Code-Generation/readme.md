## Simple REST endpoint generation in go


model/: contains the sysl application that's built

Simple/: contains all the generated code for the service

implementation/: The manual code that's written; Server config and such

gencallback.go: has the server config `Callback` and `Config` structs which are used in the `LoadServices` func that sets up the server. 

methods.go: Contains all these Functions, (yes, functions), that are then composed into the `ServiceInterface` struct in `LoadServices`

When new endpoints are added, they need to be added to the `simple.ServiceInterface` variable in `LoadServices`

main.go: runs the actual server


run `make input=<inputfile> app=<appname>` to regenerate application code
so: `make input=model/simple.sysl app=Simple` for this example

run `go run main.go` to start the server