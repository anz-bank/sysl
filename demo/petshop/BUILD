load(
    "//util/bazel/rules:sysl.bzl",
    "sysl_facade",
    "sysl_model",
    "sysl_xsd",
)

sysl_model(
    name = "model",
    srcs = ["petshop.sysl"],
    app = "PetShopModel",
    entities = [
        "Employee",
        "Breed",
        "Pet",
        "EmployeeTendsPet",
    ],
    module = "/petshop",
    package = "io.sysl.demo.petshop.model",
    root = "demo/petshop",
    visibility = ["//visibility:public"],
)

sysl_facade(
    name = "facade",
    srcs = ["petshop.sysl"],
    app = "PetShopFacade",
    model = ":model",
    module = "/petshop",
    package = "io.sysl.demo.petshop.facade",
    root = "demo/petshop",
    visibility = ["//visibility:public"],
)

sysl_model(
    name = "api",
    srcs = ["petshop.sysl"],
    app = "PetShopApi",
    entities = [
        "PetShop",
        "Employee",
        "Breed",
        "Pet",
    ],
    module = "/petshop",
    package = "io.sysl.demo.petshop.api",
    root = "demo/petshop",
    serializers = ["xml_*"],
    visibility = ["//visibility:public"],
)

sysl_model(
    name = "model_to_api",
    srcs = ["petshop.sysl"],
    app = "PetShopModelToApi",
    entities = [
        "PetRankedByLeg",
    ],
    module = "/petshop",
    package = "io.sysl.demo.petshop.views",
    root = "demo/petshop",
    serializers = [],
    visibility = ["//visibility:public"],
    deps = [
        ":api",
        ":model",
    ],
)

sysl_model(
    name = "api_to_model",
    srcs = ["petshop.sysl"],
    app = "PetShopApiToModel",
    entities = [],
    module = "/petshop",
    package = "io.sysl.demo.petshop.views",
    root = "demo/petshop",
    serializers = [],
    visibility = ["//visibility:public"],
    deps = [
        ":api",
        ":model",
    ],
)
