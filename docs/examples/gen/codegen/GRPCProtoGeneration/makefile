

all: sysl grpc server
# Build sysl into GRPC
sysl: *.sysl
	sysl tmpl --template grpc.sysl --app-name hello --start start --outdir hello hello
# Execure proto to generate go code
grpc: *
	protoc -I hello/ hello/hello.proto --go_out=plugins=grpc:hello