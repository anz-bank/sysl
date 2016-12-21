py_binary(
    name = "sysl",
    srcs = ["sysl.py"],
    deps = [
        ":sysllib",
        "//src/util",
    ],
)

py_library(
    name = "sysllib",
    srcs = [
        "syslalgo.py",
        "sysldata.py",
        "syslints.py",
        "syslloader.py",
        "syslparse.py",
        "syslseqs.py",
        "syslx.py",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//src/proto:sysl",
        "//src/util",
    ],
)
