py_test(
  name = "reljam_test",
  srcs = ["reljam_test.py"],
  size = "small",
  deps = ["//src/exporters:reljam"],
  data = ["//test/data/petshop:petshopgroup"],
)

py_test(
  name = "xsd_test",
  srcs = ["xsd_test.py"],
  size = "small",
  deps = ["//src/exporters:reljam"],
  data = ["//test/data:test_data_group"],
)

java_test(
  name = "test_lazy_load",
  size = "small",
  srcs = [
    "TestLazyLoad.java",
  ],
  test_class = "TestLazyLoad",
  deps = [
    "//test/data/petshop:api",
    "//test/data/petshop:api_to_model",
    "//test/data/petshop:facade",
    "//test/data/petshop:model",
    "//test/data/petshop:model_to_api",
    "//java/io/sysl",
  ],
)