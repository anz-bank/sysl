filegroup(
  name = "test_data_group",
  srcs = [x for x in glob(["**/*.sysl","**/*.xsd","**/*.xlsx"]) if '~$' not in x],
  visibility = ["//test/sysl:__pkg__","//test/exporters:__pkg__","//test/importers:__pkg__"]
  )
