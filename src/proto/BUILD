package(default_visibility = ["//visibility:public"])

genrule (
    name = "sysl_py",
    srcs = ["sysl.proto"],
    outs = ["sysl_pb2.py"],
    cmd = "protoc -Isrc/proto --python_out=$(@D) $<",
)

py_library(
  name = "sysl",
  srcs = [":sysl_py"],
)
