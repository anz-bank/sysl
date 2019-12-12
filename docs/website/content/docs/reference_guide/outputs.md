---
title: "Output formats"
description: "Explore various output formats from sequence diagrams to Java code."
date: 2018-02-28T14:05:31+11:00
weight: 60
draft: false
bref: "Explore various output formats from diagrams to code"
toc: true
---

## Sysl outputs

| Command | Description |
|---------|-------------|
| data    | Data Model diagrams |
| ints    | Integration Diagrams |
| sd      | Sequence Diagrams |
| pb      | Binary Protocol Buffer files of the Sysl definitions |
| protobuf  | Text based Protocol Buffer files of the Sysl definitions |
| export  | Export sysl to Swagger/Open API specification |
| codegen | Generate code with sysl transform models | 
| datamodel| ... | 


## Sysl examples

`sysl` can generate diagrams - Data model diagrams, Integration Diagrams and Sequence Diagrams - and Protobuf intermediate representations from `*.sysl` input files.

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

## Reljam examples


## Java Model
In the following example we will use `reljam model` to generate Java source code for a sysl data model.

The input file `reljam-model.sysl` contains:

```
HelloWorld [package="io.sysl.demo"]:
    !table Message:
        id <: int [~pk, ~autoinc]
        text <: string(50)
```
When executing

    reljam model reljam-model.sysl HelloWorld

the directory `io/sysl/demo` is created. It contains the following Java source files:

    HelloWorld.java
    HelloWorldException.java
    HelloWorldJsonDeserializer.java
    HelloWorldJsonSerializer.java
    HelloWorldXmlDeserializer.java
    HelloWorldXmlSerializer.java
    Message.java

### XSD

In this example we will create an XSD file from a sysl data model with `reljam xsd`.

The content of the input file `reljam-xsd.sysl` is:

```
Model:
    !table Element1:
        attr <: int [~xml_attribute]
        element2 <: Element2.key

    !table Element2:
        key <: int [~pk]
        field <: string
```
When executing

    reljam xsd reljam-xsd.sysl

the following `Model.xsd` file is created:

```
<?xml version="1.0" encoding="UTF-8"?>
<xs:schema version="1.0" [...] >
  <xs:element name="Model">
    <xs:complexType>
      <xs:sequence maxOccurs="1" minOccurs="1">
        <xs:element type="Element1List" name="Element1List" [...] />
        <xs:element type="Element2List" name="Element2List" [...] />
      </xs:sequence>
    </xs:complexType>
 [...]
```


