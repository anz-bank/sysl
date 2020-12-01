---
id: cmd-protobuf
title: Protobuf
sidebar_label: protobuf
keywords:
  - command
---

## Summary

`sysl protobuf` converts a Sysl file to a [protobuf](https://developers.google.com/protocol-buffers) representation, which can be consumed by programs in other languages. The output is encoded using the [Sysl protobuf definition](https://github.com/anz-bank/sysl/blob/master/pkg/sysl/sysl.proto).

## Usage

```bash
sysl protobuf --mode=[textpb, json, pb] --output=<output file name> [<flags>]
```

## Optional flags

- `--mode` Output file format which can be `textpb` (the [prototext format](https://pkg.go.dev/google.golang.org/protobuf/encoding/prototext?tab=doc)), `json` (JSON with some embedding of types) or `pb` (binary proto encoding).
- `-o, --output` Output filename. If not provided, the protobuf output will be printed to `stdout`.

[More common optional flags](common-flags.md)

## Arguments

- `<MODULE>` Input sysl file that contains the system specifications. e.g `simple.sysl`. The `.sysl` file type is optional.

## Examples

### `textpb` mode

Command line

```bash
sysl protobuf --mode=textpb --output=simple.textpb simple.sysl
```

```sysl title="Input Sysl file: simple.sysl"
Simple "Simple":
    @description =:
        | Simple demo for protobuf export

    # definitions
    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
```

```sysl title="Output textpb file: simple.textpb"
apps: {
  key: "Simple"
  value: {
    name: {
      part: "Simple"
    }
    long_name: "Simple"
    attrs: {
      key: "description"
      value: {
        s: "Simple demo for protobuf export\n"
      }
    }
    types: {
      key: "SimpleObj"
      value: {
        tuple: {
          attr_defs: {
            key: "name"
            value: {
              primitive: STRING
              attrs: {
                key: "json_tag"
                value: {
                  s: "name"
                }
              }
              opt: true
              source_context: {
                file: "simple.sysl"
                start: {
                  line: 8
                  col: 16
                }
                end: {
                  line: 8
                  col: 30
                }
              }
            }
          }
        }
        source_context: {
          file: "simple.sysl"
          start: {
            line: 6
            col: 4
          }
          end: {
            line: 8
            col: 30
          }
        }
      }
    }
    source_context: {
      file: "simple.sysl"
      start: {
        line: 1
        col: 1
      }
      end: {
        line: 1
        col: 7
      }
    }
  }
}
```

### `json` mode

Command line

```bash
sysl protobuf --mode=json --output=simple.json simple.sysl
```

```sysl title="Input Sysl file: simple.sysl"
Simple "Simple":
    @description =:
        | Simple demo for protobuf export

    # definitions
    !type SimpleObj:
        name <: string?:
            @json_tag = "name"
```

```sysl title="Output json file: simple.json"
{
 "apps": {
  "Simple": {
   "name": {
    "part": [
     "Simple"
    ]
   },
   "longName": "Simple",
   "attrs": {
    "description": {
     "s": "Simple demo for protobuf export\n"
    }
   },
   "types": {
    "SimpleObj": {
     "tuple": {
      "attrDefs": {
       "name": {
        "primitive": "STRING",
        "attrs": {
         "json_tag": {
          "s": "name"
         }
        },
        "opt": true,
        "sourceContext": {
         "file": "simple.sysl",
         "start": {
          "line": 8,
          "col": 16
         },
         "end": {
          "line": 8,
          "col": 30
         }
        }
       }
      }
     },
     "sourceContext": {
      "file": "simple.sysl",
      "start": {
       "line": 6,
       "col": 4
      },
      "end": {
       "line": 8,
       "col": 30
      }
     }
    }
   },
   "sourceContext": {
    "file": "simple.sysl",
    "start": {
     "line": 1,
     "col": 1
    },
    "end": {
     "line": 1,
     "col": 7
    }
   }
  }
 }
}
```
