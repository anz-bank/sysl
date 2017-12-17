Protobuf usage
==============

The [Protocol Buffer Compiler](https://developers.google.com/protocol-buffers/) `protoc >= 3.5.0` is used to generate `src/sysl/proto/sysl_pb2.py` from the protocol buffer definition file `src/proto/sysl.proto`: 

```
protoc --python_out=src/sysl/proto  --proto_path=src/proto sysl.proto
```
Currently the Python based importers, exporters and the sysl core use the generated `sysl_pb2.py` module.  The generated `sysl_pb2.py` file has been added to this repo in order to avoid the additional third party dependency on the `protoc` compiler.

If `sysl.proto` is updated the above command needs to be re-run to update the corresponding Python definitions in `sysl_pb2.py`.

Several other target languages beyond Python are also supported by the `protoc` compiler if required.

