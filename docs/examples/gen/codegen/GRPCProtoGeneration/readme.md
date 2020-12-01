#### GRPC example

Make sure to run install.sh in order to get the required git modules

### Sysl grpc examples

GRPCProtoGeneration:

Makefile: has sysl command to build grpc sysl into protos

grpc.sysl: has transform needed to generate .proto file

hello.sysl: has the sysl that's compiled into .proto in the hello directory

hello: Has the .proto and the pb.go files generated from sysl and protoc respectively
