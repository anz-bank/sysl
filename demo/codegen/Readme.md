# demo/

## Warning
These examples have test case dependencies, think carefuly before changing them

# Generation of grpc client server
Using Sysl, you can specify endpoints, endpoint behaviour and the underlying data models. The Sysl compiler could automatically generate grpc client and server stubs from a common Sysl specification.

This is how you generate grpc code stubs, it assumes all the Sysl dependencies have been previously installed.
1. Execute the below command in the codegen directory:
    sysl tmpl --root ./AuthorisationAPI --root-template . --template grpc.sysl --app-name AuthorisationAPI --start start --outdir ./AuthorisationAPI authorisation
2. The output .proto is generated in the output folder.
3. Run Protoc to generate Go Code:
    protoc --go_out=plugins=grpc:. AuthorisationAPI/AuthorisationAPI.proto
