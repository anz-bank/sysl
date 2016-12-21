java_library(
    name = "sysl",
    srcs = [":sources"],
    visibility = ["//visibility:public"],
    deps = ["@joda_time_joda_time//jar"],
)

filegroup(
    name = "sources",
    srcs = glob(["**/*.java"]),
    visibility = ["//util:__pkg__"],
)
