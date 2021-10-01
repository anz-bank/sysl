# Why are there bundled files?
There are certain core proto files that are required to implement protobufs in a language. In Java they are included in the [pom](https://github.com/protocolbuffers/protobuf/blob/master/java/core/pom.xml). For the Arr.ai implementation they are placed in `bundled_files` and used to populate `bundled_files/local_imports.arrai`.

# How do bundled files work?

When an import statement is being processed, the parser will first attempt to resolve it within the import paths provided in the command line.

`sysl import example/proto/foo.proto --import-paths example/proto example/external`

If `foo.proto` contains an import that cannot be found in `example/proto` or `example/external`, the parser will check `local_imports.arrai`

# How do I update the bundled files?
Any files with the `.proto` extension within the `bundled_files` directory will be picked up and added to `bundled_files/local_imports.arrai` when you run `make bundled-proto`. Any bundled files should be placed in folders that represent the package structure.
