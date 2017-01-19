package(default_visibility = ["//visibility:public"])

load(
  "//util/bazel/rules:sysl.bzl",
  "sysl_facade",
  "sysl_model",
)

sysl_model(
  name = "model",
  srcs = ["petshop.sysl"],
  root = "test/data/petshop",
  module = "/petshop",
  app = "PetShopModel",
  package = "io.sysl.test.data.petshop.model",
  entities = [
    "Employee",
    "Breed",
    "Pet",
    "EmployeeTendsPet",
  ],
)

sysl_facade(
  name = "facade",
  srcs = ["petshop.sysl"],
  root = "test/data/petshop",
  module = "/petshop",
  app = "PetShopFacade",
  model = ":model",
  package = "io.sysl.test.data.petshop.facade",
)

sysl_model(
  name = "api",
  srcs = ["petshop.sysl"],
  root = "test/data/petshop",
  module = "/petshop",
  app = "PetShopApi",
  package = "io.sysl.test.data.petshop.api",
  entities = [
    "PetShop",
    "Employee",
    "Breed",
    "Pet",
  ],
  serializers = ["xml_*"],
)

sysl_model(
  name = "model_to_api",
  srcs = ["petshop.sysl"],
  root = "test/data/petshop",
  module = "/petshop",
  app = "PetShopModelToApi",
  package = "io.sysl.test.data.petshop.views",
  entities = [
    "PetRankedByLeg",
  ],
  deps = [
    ":api",
    ":model",
  ],
  serializers = [],
)

sysl_model(
  name = "api_to_model",
  srcs = ["petshop.sysl"],
  root = "test/data/petshop",
  module = "/petshop",
  app = "PetShopApiToModel",
  package = "io.sysl.test.data.petshop.views",
  entities = [],
  deps = [
    ":api",
    ":model",
  ],
  serializers = [],
)

# export for use in the "exporters" tests
filegroup(
  name = "petshopgroup",
  srcs = ["petshop.sysl"],
  visibility = ["//test/exporters:__pkg__"],
)
