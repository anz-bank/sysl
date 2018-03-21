---
title: "Output formats"
description: "Explore various output formats from sequence diagrams to Java code."
date: 2018-02-28T14:05:31+11:00
weight: 60
draft: false
bref: "Explore various output formats from diagrams to code"
toc: true
---

Sysl consists of two executables: `sysl` and `reljam`.</br> The **Sy**stem **S**pecification **L**anguage `sysl` is mainly concerned with diagram creation whereas the **Rel**ational **Ja**va **M**odel program `reljam` generates different types of source code output.

Sysl outputs
------------
| Command | Description |
|---------|-------------|
| data    | Data Model diagrams |
| ints    | Integration Diagrams |
| sd      | Sequence Diagrams |
| pb      | Binary Protocol Buffer files of the Sysl definitions (plugins)    |
| textpb  | Text based Protocol Buffer files of the Sysl definitions (plugins, debugging) |

Reljam outputs
--------------
| Command | Description |
|---------|-------------|
| model   | Java model implementation (in memory) |
| facade  | Java facade implementation (restricted access to creating and populating models) |
| view    | Java implementation of Sysl model transformations|
| xsd     | XSD representation of Sysl model |
| swagger | Swagger representation of REST APIs and models |
| spring-rest-service | Java Spring REST API implementation |

Sysl examples
-------------
`sysl` can generate diagrams - Data model diagrams, Integration Diagrams and Sequence Diagrams - and Protobuf intermediate representations from `*.sysl` input files.

### Integration Diagrams

Input file `sysl-ints.sysl`:

```
IntegratedSystem:
    integrated_endpoint_1:
        System1 <- endpoint
    integrated_endpoint_2:
        System2 <- endpoint

System1:
    endpoint: ...
System2:
    endpoint: ...

Project [appfmt="%(appname)"]:
    _:
        IntegratedSystem
        System1
        System2
```

Sysl command:

	sysl ints -o sysl-ints.svg sysl-ints.sysl --project Project

Output diagram:

![Integration diagram](/img/sysl/simple-sysl-ints.svg)

add the `--epa` (endpoint analysis) option for a more detailed diagrams.


### Data Model Diagrams

Input file `sysl-data.sysl`:

```
Project:
    _:
        App

App:
    !type Order:
        order_id <: int

    !type Customer:
        customer_id <: int
        orders <: set of Order
```

Sysl command:

	sysl data -o sysl-data.svg -j Project sysl-data.sysl

Output diagram:

![Data model diagram](/img/sysl/simple-sysl-data.svg)

### Sequence Diagrams

Input file `sysl-sd.sysl`:

```
Database[~db]:
    QueryUser (user_id):
        Return User

Api:
    /users/{user_id<:int}/profile:
        GET:
                Database <- QueryUser(user_id)
                Return UserProfile

WebFrontend:
    RequestProfile:
        Api <- GET /users/{user_id}/profile
        Return Profile Page

Project [seqtitle="Profile"]:
    _:
        WebFrontend <- RequestProfile
```

Sysl command:

	sysl sd -a Project -o "sysl-sd-%(epname)".png sysl-sd.sysl

Output diagram:

![Sequence diagram](/img/sysl/simple-sysl-sd.svg)

### Text based Protocol Buffer output
Protocol buffers is a "language-neutral, platform-neutral, extensible mechanism for serializing structured data – think XML, but smaller, faster, and simpler". It is a strongly typed binary format used as intermediate representations of Sysl definitions comparable to an abstract syntax tree. The strongly typed protocol buffers are supported in most major programming languages.

Please refer to our developer documentation on how to compile the Protobuf definitions to your preferred porgramming language in order to [create your own Sysl extension]
(https://github.com/anz-bank/sysl#extending-sysl). If you want to generate human readable, text-based Protobuf output use the `textpb` command.

For the following contents of `hello.sysl`

```
HelloWorld:
    !type Message:
        text <: string
```

the command

	sysl textpb hello.sysl --out hello.textpb

generates a `hello.textpb` file. Its contents are

```
apps {
  key: "HelloWorld"
  value {
    name {
      part: "HelloWorld"
    }
    types {
      key: "Message"
      value {
        tuple {
          attr_defs {
            key: "text"
            value {
              primitive: STRING
              source_context {
                start {
                  line: 4
                }
              }
            }
          }
        }
      }
    }
  }
}
```


