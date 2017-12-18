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
  deps = [
    "//src/proto:sysl",
    "//src/util",
  ],
  visibility = ["//visibility:public"],
)
