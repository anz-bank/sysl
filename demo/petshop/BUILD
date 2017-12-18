load("//util/bazel/rules:sysl.bzl",
  "sysl_facade",
  "sysl_model",
  "sysl_xsd",
)

sysl_model(
  name = "model",
  srcs = ["petshop.sysl"],
  root = "demo/petshop",
  module = "/petshop",
  app = "PetShopModel",
  package = "io.sysl.demo.petshop.model",
  entities = [
    "Employee",
    "Breed",
    "Pet",
    "EmployeeTendsPet",
  ],
  visibility = ["//visibility:public"],
)

sysl_facade(
  name = "facade",
  srcs = ["petshop.sysl"],
  root = "demo/petshop",
  module = "/petshop",
  app = "PetShopFacade",
  model = ":model",
  package = "io.sysl.demo.petshop.facade",
  visibility = ["//visibility:public"],
)

sysl_model(
  name = "api",
  srcs = ["petshop.sysl"],
  root = "demo/petshop",
  module = "/petshop",
  app = "PetShopApi",
  package = "io.sysl.demo.petshop.api",
  entities = [
    "PetShop",
    "Employee",
    "Breed",
    "Pet",
  ],
  serializers = ["xml_*"],
  visibility = ["//visibility:public"],
)

sysl_model(
  name = "model_to_api",
  srcs = ["petshop.sysl"],
  root = "demo/petshop",
  module = "/petshop",
  app = "PetShopModelToApi",
  package = "io.sysl.demo.petshop.views",
  entities = [
    "PetRankedByLeg",
  ],
  deps = [
    ":api",
    ":model",
  ],
  serializers = [],
  visibility = ["//visibility:public"],
)

sysl_model(
  name = "api_to_model",
  srcs = ["petshop.sysl"],
  root = "demo/petshop",
  module = "/petshop",
  app = "PetShopApiToModel",
  package = "io.sysl.demo.petshop.views",
  entities = [],
  deps = [
    ":api",
    ":model",
  ],
  serializers = [],
  visibility = ["//visibility:public"],
)