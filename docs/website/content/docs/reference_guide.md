---
title: "Quick reference guide"
date: 2018-02-28T10:11:24+11:00
description: ""
weight: 50
bref: ""
toc: true
---


Sysl allows you to specify Application behaviour and Data Models that are shared between your applications. This is useful in many use cases, especially in software projects with many moving parts which need their documentation kept up to date.

To explain these concepts, we will design an application called `MobileApp` which interacts with another application called `Server`.

## Applications


An __application__ is an independent entity that provides services via its various __endpoints__.

Here is how an application is defined in sysl.
```
MobileApp:
  ...

Server:
  ...
```
`MobileApp` and `Server` are user-defined Applications that do not have any endpoints yet. We will design this app as we move along.

Notes about sysl syntax:

  * `:` and `...` have special meaning. `:` followed by an `indent` is used to create a parent-child relationship.
    * All lines after `:` should be indented. The only exception to this rule is when you want to use the shortcut `...`.
    * The `...` (aka shortcut) means that we don't have enough details yet to describe how this endpoint behaves. Sysl allows you to take an iterative approach in documenting the behaviour. You add more as you know more.

## Endpoints


Endpoints are the services that an application offers. Let's add endpoints to our `MobileApp`.
```
MobileApp:
  Login: ...
  Search: ...
  Order: ...
```
Now, our `MobileApp` has three `endpoints`: `Login`, `Search` and `Orders`.

Notes about sysl syntax:

 * Again, `...` is used to show we don't have enough details yet about each endpoint.
 * All endpoints are indented. Use a `tab` or `spaces` to indent.
 * These endpoints can also be REST api's. See section on [Rest](#rest) below on how to define rest api endpoints.

Each endpoint should have statements that describe its behaviour. Before that lets took at data types and how it can used in sysl.

You will have various kinds of data passing through your systems. Sysl allows you to express ownership, information classification and other attributes of your data in one place.

Continuing with the previous example, let's define a `Server` that expects `LoginData` for the `Login` Flow.
```
Server:
  Login (request <: Server.LoginData): ...

  !type LoginData:
    username <: string
    password <: string
```
In the above example, we have defined another application called `Server` that has an endpoint called `Login`. It also defines a new data type called `LoginData` that it expects callers to pass in the login call.

Notes about sysl syntax:
  * `<:` is used to define the arguments to `Login` endpoint.
  * `!type` is used to define a new data type `LoginData`.
    * Note the indent to create fields `username` and `password` of type `string`.
    * See [Data Types](#data-types) to see what all the supported data types.
  * Data types (like `LoginData`) belong to the app under which it is defined.
  * Refer to the newly defined type by its fully qualified name. e.g. `Server.LoginData`.

## Data Types


Sysl supports following data types out of the box.
```
double 
int64 
float64 
string 
bool 
date.Date 
time.Time 
```
We can Also define our own datatypes using the `!type` keyword within an application.

```
!type response:
  data <: string
  type <: int
```

Now, we have two apps `MobileApp` and `Server`, but they do not interact with each other. Time to add some statements.


#### Return response
An endpoint can return response to the caller. Everything after `return` keyword till the end-of-line is considered response payload. You can have:
  * string - a description of what is returned, or
  * Sysl type - formal type to return to the caller

```
MobileApp:
  Login: 
    return string
  Search: 
    return string
  Order: 
    return sequence of string
```



## Statements

Our `MobileApp` does not have any detail yet on how it behaves. Let's use sysl statements to describe behaviour. Sysl supports following types of statements:
  * [Text](#text)
  * [Call](#Call)
  * [Return](#return-response)
  * [Control Statements](#control-statements)
  * [Arguments](#arguments)

#### Text
Use simple text to describe behaviour. See below for examples of text statements:
```
Server:
  Login:
    do input validation
    "Use special characters like \n to break long text into multiple lines"
    'Cannot use special characters in single quoted statements'
```

#### Call
A standalone service that does not interact with anybody is not a very useful service. Use the `call` syntax to show interaction between two services.

In the below example, MobileApp makes a call to backend Server which in turn calls database layer.
  
```
MobileApp:
  Login:
    Server <- Login

Server:
  Login(data <: LoginData):
    build query
    DB <- Query
    check result
    return Server.LoginResponse

  !type LoginData:
    username <: string
    password <: string

  !type LoginResponse:
    message <: string

DB:
  Query:
    lookup data
    return data
  Save:
    ...
```
See [/assets/call.sysl](/assets/call.sysl) for complete example.

Now we have all the ingredients to draw a sequence diagram. Here is one generated by sysl for the above example:

![](/assets/call-Seq.png)

See [Generate Diagrams](#generate-diagrams) on how to draw sequence and other types of diagrams using sysl.


## Control flows


## If/else

Sysl allows you to express high level of detail about your design. You can specify decisions, processing loops etc.

##### If, else
You can express an endpoint's critical decisions using IF/ELSE statement:
```
Server:
  HandleFormSubmit:
    validate input
    IF session exists:
      use existing session
    Else:
      create new session
    process input
```
See [/assets/if-else.sysl](/assets/if-else.sysl) for complete example.

`IF` and `ELSE` keywords are case-insensitive. Here is how sysl will render these statements:

![](/assets/if-else-Seq.png)

## For, Loop, Until, While

Express processing loop using FOR:
```
Server:
  HandleFormSubmit:
    validate input
    For each element in input:
      process element
```
See [/assets/for-loop.sysl](/assets/for-loop.sysl) for complete example.

`FOR` keyword is case insensitive. Here is how sysl will render these statements:

![](/assets/for-loop-Seq.png)

You can use `Loop`, `While`, `Until`, `Loop-N` as well (all case-insensitive).



## Multiple Declarations

Sysl allows you to define an application in multiple places. There is no redefinition error in sysl.

```
UserService:
  Login: ...

UserService:
  Register: ...
```

Result will be as-if it was declared like so:

```
UserService:
  Login: ...
  Register: ...
```

## Projects
Most of the changes to your system will be done as part of a well defined `project` or a `software release`.

`TODO: Elaborate`

## Imports


To keep things modular, sysl allows you to import definitions created in other `.sysl` files.

E.g. `server.sysl`
```
Server:
  Login: ...
  Register: ...
```

and you use `import` in `client.sysl`

```
import server

Client:
  Login:
    Server <- Login
```

Above code assumes, server and client files are in the same directory. If they are in different directories, you must have atleast a common root directory and `import /path/from/root`.

All sysl commands accept `--root` argument. Run `sysl -h` or `reljam -h` for more details.



## Generate Diagrams

Once your design is complete, its time to get some output from Sysl. Sysl supports generating diagrams of following types:
  * Sequence Diagram
  * Integration Diagram
  * Data Diagram

Sysl aims to generate code and documentation from only one source of truth i.e. `.sysl` files.

## Sequence Diagrams

You can generate the Sequence Diagram using the following command:

```
sysl sd -o 'call-login-sequence.png' -s 'MobileApp <- Login' call.sysl
```
You can omit the the `.sysl` and sysl will pickup the correct file.
```
sysl sd -o 'call-login-sequence.png' -s 'MobileApp <- Login' call
```

Here is the output that you should see:

![](//assets/call-Seq.png)

See [/assets/call.sysl](/assets/call.sysl) for complete example.

### How sysl generates sequence diagram?

Let's breakdown the `sd` aka `sequence diagram` command:
```
sysl sd -o 'call-login-sequence.png' -s 'MobileApp <- Login' call.sysl
```
  * `-o` specifies the output filename
  * `-s` specifies the start endpoint
  * `call.sysl` the source to start the analysis from

Sysl analyzes the starting endpoint and finds all the `call`s that this endpoint makes to other endpoints (including the ones to other applications). It finds all the transitive dependencies till there are none.

In the above diagram, `DB` is the last app in this flow. Sysl also captures the return data that each endpoint returns to its caller. See below for more details.

#### Format Arguments
The default diagram by default only shows the data type that is returned by an endpoint. You can instruct `sysl` to show the arguments to your endpoint in a sequence diagram.

Command:

`sysl sd -o 'call-login-sequence.png' --epfmt '%(epname) %(args)' -s 'MobileApp <- Login' /assets/call.sysl -v call-login-sequence.png`
See [/assets/args.sysl](/assets/args.sysl) for complete example.

![](/assets/args-Seq.png)

A bit more explanation is required regarding `epname` and `args` keywords that are used in `epfmt` command line argument. See section on [Attributes](#epfmt) below.

## Integration Diagram

`TODO`

See: run `sysl ints -h` for more details.

## Data Diagram

See [Data Models](#data-models) on types of data models and how to render them.


## Type
`!type`
The type keyword is used to define a type. 
In the following example we define a `Post` type made up of multiple attributes.
```
  !type Post:
    userId <: int
    id <: int
    title <: string
    body <: string
```

## Alias
`!alias`
Alias' can be used to simplify a type;

```
  !alias Posts:
    sequence of Post
```
## View
`!view`
Views are sysl's functions; we can use them in the transformation language, see [docs/transformation.html]for more info

## Union
`!union`
Unions are a union type; 
`!union string, int32`
can either be a string, int32, but not both.

## Table
`!table` Add more here

## Wrap
`!wrap` Add more here



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
| info | Build information for sysl executable |
| env | Sysl environment variables value |


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



