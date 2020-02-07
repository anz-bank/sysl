

all: sysl grpc server catalog
# Build sysl into GRPC
sysl: *.sysl
	sysl tmpl  --template grpc.sysl --app-name hello --start start --outdir gen  hello
# Execure proto to generate go code
grpc: *
	protoc -I gen/ gen/hello.proto --go_out=plugins=grpc:gen

# Run the server implementation
server: */*
	go run server/main.go &


