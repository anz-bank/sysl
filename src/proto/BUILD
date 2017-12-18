package(default_visibility = ["//visibility:public"])

load("@org_pubref_rules_protobuf//python:rules.bzl", "py_proto_compile")

py_proto_compile(
    name = "sysl_py",
    protos = ["sysl.proto"],
    with_grpc = False,
)

py_library(
    name = "sysl",
    srcs = [":sysl_py"],
)