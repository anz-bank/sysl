# demo/

## Warning
These examples have test case dependencies, think carefuly before changing them

# Generation of grpc client server
Using Sysl, you can specify endpoints, endpoint behaviour and the underlying data models. The Sysl compiler could automatically generate grpc client and server stubs from a common Sysl specification.

This is how you generate grpc code stubs, it assumes all the Sysl dependencies have been previously installed.
1. Change into <SYSL_GIT_REPO_ROOT>.
2. Execute the below command:
    sysl tmpl --root demo/codegen --root-template demo/codegen --template grpc.sysl --app-name AuthorisationAPI --start start --outdir demo/codegen authorisation
3. The output .proto is generated in the output folder.
4. Run Protoc to generate Go Code:
    protoc --go_out=plugins=grpc:. demo/codegen/AuthorisationAPI/AuthorisationAPI.proto

SYSL_GIT_REPO_ROOT = Directory where Sysl repository has been cloned.
