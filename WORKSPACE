git_repository(
    name = "io_bazel_rules_jsonnet",
    remote = "https://github.com/bazelbuild/rules_jsonnet",
    tag = "0.0.2",
)

maven_jar(
    name = "joda_time_joda_time",
    artifact = "joda-time:joda-time:2.9.6",
    sha1 = "e370a92153bf66da17549ecc78c69ec6c6ec9f41",
)

maven_jar(
    name = "org_apache_commons_commons_lang3",
    artifact = "org.apache.commons:commons-lang3:3.5",
    sha1 = "6c6c702c89bfff3cd9e80b04d668c5e190d588c6",
)

maven_jar(
    name = "com_fasterxml_jackson_core_jackson_core",
    artifact = "com.fasterxml.jackson.core:jackson-core:2.8.5",
    sha1 = "60d059f5d2930ccd1ef03535b713fd9f933d1ba7",
)

maven_jar(
    name = "com_fasterxml_jackson_core_jackson_databind",
    artifact = "com.fasterxml.jackson.core:jackson-databind:2.8.5",
    sha1 = "b3035f37e674c04dafe36a660c3815cc59f764e2",
)

load("@io_bazel_rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_repositories")

jsonnet_repositories()

git_repository(
    name = "org_pubref_rules_protobuf",
    remote = "https://github.com/pubref/rules_protobuf",
    tag = "v0.7.1",
)

load("@org_pubref_rules_protobuf//python:rules.bzl", "py_proto_repositories")

py_proto_repositories()
