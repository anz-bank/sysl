module.exports = {
  docs: {
    "Getting Started": [
      "introduction",
      "installation",
      "tutorial",
      "tutorial-codegen",
      // 'examples-dev',
      // 'examples-nondev',
    ],
    Features: [
      "features",
      "gen-diagram",
      "gen-code",
      // 'gen-docs',
      // 'gen-db',
      // 'gen-test',
    ],
    "Best Practices": ["best-practices/intro", "best-practices/formatting"],
    Reference: [
      "ref-summary",
      {
        Language: [
          "lang/intro",
          "lang/application",
          "lang/endpoint",
          "lang/statement",
          "lang/type",
          "lang/primitives",
          "lang/field",
          "lang/annotation",
          "lang/tag",
          "lang/table",
          "lang/enum",
          "lang/alias",
          "lang/union",
          "lang/view",
          "lang/mixin",
          "lang/module",
          "lang/import",
          "lang/comment",
          "lang/identifiers",
          "lang/source-context",
          "lang/project",
          "lang-spec",
        ],
      },

      // {
      //   "Supported Formats": [
      //     "formats-overview",
      //     "formats-avro",
      //     "formats-openapi",
      //     "formats-spanner",
      //     "formats-xsd",
      //   ]
      // },

      {
        "Sysl for TypeScript": [
          "ts/intro",
          "ts/pbmodel",
          "ts/model",
        ]
      },

      {
        "Sysl CLI": [
          "cmd/cmd",
          "cmd/common-flags",
          "cmd/cmd-validate",
          "cmd/cmd-diagram",
          "cmd/cmd-info",
          "cmd/cmd-env",
          "cmd/cmd-sd",
          "cmd/cmd-integrations",
          "cmd/cmd-datamodel",
          // 'cmd/cmd-db',
          // 'cmd/cmd-db-delta',
          // 'cmd/cmd-codegen',
          "cmd/cmd-import",
          "cmd/cmd-export",
          "cmd/cmd-protobuf",
          "cmd/cmd-repl",
          // 'cmd/cmd-template',
        ],
      },
      {
        "Sysl Catalog": [
          "catalog/sysl-catalog",
          "catalog/sysl-catalog-install",
          "catalog/sysl-catalog-cmd",
        ],
      },
    ],
  },
  community: ["discussions", "resources", "faq"],
};
