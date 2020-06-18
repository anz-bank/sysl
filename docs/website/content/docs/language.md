---
title: "Language Specification"
date: 2018-02-28T10:11:24+11:00
description: "Get an in depth understanding of the Sysl language"
weight: 50
bref: "Get an in depth understanding of the Sysl language"
topic: "Language"
layout: "quickreference"
toc: true
---

## Introduction

Sysl is a system modelling language designed for modelling distributed web applications. Sysl allows you to specify [Application](#applications) behaviour and Data Models that are shared between your applications. Another related concept is of software [Projects](#projects) where you can document what changes happened in each project or a release. Your complete system can be described as forest of trees, where one tree represents one application or data model. Sysl uses indentation to represent parent-child or `has` relationships. E.g. an `application` has `endpoints` or a `table` has `columns`.

To explain these concepts, we will design an application called `MobileApp` which interacts with another application called `Server`.

## Applications

An **application** is an independent entity that provides services via its various **endpoints**.

Here is how an application is defined in sysl.

```
MobileApp:
  ...

Server:
  ...
```

`MobileApp` and `Server` are user-defined Applications that do not have any endpoints yet. We will design this app as we move along.

Notes about sysl syntax:

- `:` and `...` have special meaning. `:` followed by an `indent` is used to create a parent-child relationship.
  - All lines after `:` should be indented. The only exception to this rule is when you want to use the shortcut `...`.
  - The `...` (aka shortcut) means that we don't have enough details yet to describe how this endpoint behaves. Sysl allows you to take an iterative approach in documenting the behaviour. You add more as you know more.

### Application Names

Application names can contain any lowercase or uppercase letter, number, underscore or hyphen.
Application names must begin with either a letter or an underscore.

Valid

```
  Mobile App:
    Login: ...
```

Invalid

```
  Mobile App*:
    Login: ...
```

#### Long Names

Shorter application names are preferred to enable shorter expressions when specifying references to the app. Sysl allows the definition of a LongName using the following syntax

```
  Mobile App "My Awesome Mobile Application":
    Endpoint: ...
```

#### Namespaces

Application names can also be namespaced using two colons. This allows the formation of a hierarchical structure in Sysl, so that applications can be grouped.

```
  Payments :: CreditCard "Credit Card Payment Service":
    Pay: ...
```

Namespaces can be nested arbitrarily deep.
```
  Payments :: CreditCard :: Validate "Credit Card Validator":
    Pay: ...
```

### Multiple Declarations

Sysl allows you to define an application in multiple places. There is no
redefinition error in sysl.

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

- Again, `...` is used to show we don't have enough details yet about each endpoint.
- All endpoints are indented. Use a `tab` or `spaces` to indent.
- These endpoints can also be REST api's. See section on [Rest](#rest) below on how to define rest api endpoints.

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

- `<:` is used to define the arguments to `Login` endpoint.
- `!type` is used to define a new data type `LoginData`.
  - Note the indent to create fields `username` and `password` of type `string`.
  - See [Data Types](#data-types) to see what all the supported data types.
- Data types (like `LoginData`) belong to the app under which it is defined.
- Refer to the newly defined type by its fully qualified name. e.g. `Server.LoginData`.

### Rest

Rest is a very common architectural style for defining web services. Here is how
you can define a web service:

```
Server:
  /path:
    HTTP_METHOD:
      <describe behaviour using statements>
```

where

`HTTP_METHOD` can be one of `GET`, `PUT`, `POST`, `DELETE` and `PATCH`.

Nested-Paths:

You can breakdown long URL's to reduce repetition. See below for complete
example:

```
AccountTransactionApi [package="io.sysl.account.api"]:
    /accounts [interface="Accounts"]:
        /{account_number<:int}:
            GET:
              BankDatabase <- GetAccount(account_number)

            /withdraw:
              POST (Transaction):
                BankDatabase <- WithdrawFunds(account_number)

            /transactions:
              GET ?start_date=string&end_date=string:
                BankDatabase <- QueryTransactions(account_number, start_date, end_date)

    !type Account:
        account_number <: int?
        account_type <: string?
        account_status <: string?
        account_balance <: int?

    !type Transaction:
        transaction_id <: int?
        transaction_type <: string?
        transaction_date_time <: date?
        transaction_amount <: int?
        from_account_number <: Accounts.account_number
        to_account_number <: Accounts.account_number

ATM:
    GetBalance:
        AccountTransactionApi <- GET /accounts/{account_number}
        Return balance
    Withdraw:
        AccountTransactionApi <- POST /accounts/{account_number}/withdraw
        Withdraw funds
```

Here are few things to notice

- `interface` attribute allows you to specify the name of the generated
  interface class. This interface has all the methods of your api.
- Query parameters `start_date` and `end_date` (of type `string`) that you can
  pass to `GET /accounts/{account_number}/transactions` endpoint.
- `ATM <-GetBalance` makes a call to `AccountTransactionApi <- GET /accounts/{account_number}`.

See [api.sysl](/assets/api.sysl) for complete example.

Command to generate code:

Run this in the same directory as `api.sysl`:

```
reljam spring-rest-service api AccountTransactionApi
```

### Parameter Types

Different types of parameters can be defined for an endpoint

#### REST Endpoint Specific

##### Path Parameters

Path parameters can be defined using curly brackets in the path name.

In the example above
```
AccountTransactionApi [package="io.sysl.account.api"]:
    /accounts [interface="Accounts"]:
        /{account_number<:int}:
            GET:
              BankDatabase <- GetAccount(account_number)
```

##### Query Parameters

Query parameters can be defined in the method statement using a `?queryParamName=SomePrimitiveType`
Multiple query parameters can be separated using the & character
When a reference to a type must be made, curly brackets must surround the type reference

`?myQueryParam={QueryParamType}`

```
Server:
    /first:
        GET ?depth=int&limit=int?&offset=int?:
            return ok
    /second :
        GET ?tags={Tags}
          return ok

    !alias Tags:
        sequence of string
```

#### Header and Body Parameters

Header and body parameters can be defined using brackets (foo <: int)
To create a body parameter, we use the pattern ~body e.g `(bodyParam <: int [~body])`
To create a header parameter, we use the pattern ~header e.g `(bodyParam <: int [~header])`
Larameters are separated by commas
```
Server:
    /first:
      GET (filter <: int, offset <: int, tags <: Tags)
    /second:
        GET (foo <: int [~body]) ?depth=int [~bar]:
            ...

    !alias Tags:
        sequence of string
```


## Data-Types

Sysl supports following data types out of the box.

```
decimal
int64
float64
string
bool
date
datetime
```

We can also define our own datatypes using the `!type` keyword within an application.

```
!type response:
  data <: string
  type <: int
```

### Type-Names

Type names follow the same rules as application names:
Type names can contain any lowercase or uppercase letter, number, underscore or hyphen.
Type names must begin with either a letter or an underscore.
Type names MAY contain spaces (this may change in the future)

#### Special Characters

If special characters such as ':' or '.' are needed in a type or endpoint name, this can be expressed in Sysl by using their URL encoded equivalent instead. For an example, refer to [Special Characters](https://sysl.io/docs/byexample/special-characters/). [Url Encoder](https://www.urlencoder.org/)

### Optional Types

We can define optional parameters and fields in sysl using a postfix `?`.

The following example defines an optional token field in the type response

```
!type response:
  data <: string
  type <: int
  token <: string?
```

### Enumerations

You can define an enumeration using the `!enum` syntax to give a name to an
enumeration type and list the enumerated names and their values. For example:

```
Server:
  Login (request <: Server.LoginData):
    return Server.Code

  !enum Code:
    success: 1
    invalid: 2
    tooManyAttempts: 3

  !type LoginData: ...
```

An enumeration is a type and can be referenced in the same way that other types
are referenced else where in a sysl specification.

NOTE: The syntax for enumerations will likely change from `name: value` to
`name = value` in future. Limitations in the current parser prevent the second
form from parsing.


### Return response

An endpoint can return a response to the caller. Everything after `return` keyword till the end-of-line is considered response payload. 

You can return:

- A a primitive Sysl type
  - e.g `return string`
  - Can be one of any [primitive sysl type](#Data-Types)
- A Sysl type - formal type to return to the caller
  - e.g `return LoginData`
- An expression of a Sysl Type
  - e.g `return sequence of string`
- A named return response
  - e.g `return error <: string`
  - e.g `return 200 <: OrderData`
  - e.g `return 200 <: AnotherApp.OrderData`
- An empty response
  - e.g `return ok`

```
MobileApp:
  Name:
    return string
  Login:
    return Server.LoginData
  Search:
    return sequence of string
  Order:
    return UserPreference

  !type UserPreference
    Geography <: string

Server:
  !type LoginData
    userID <: string
```


### Attributes

You can attach more metadata to your application specification. This information
can be used by sysl plugins to extend the default functionality. The attributes
are added inside square brackets `Application [...attributes]`. Sysl attributes
are of two types: `Patterns` and `Key-Value` pairs:

- Patterns

  A pattern is `~` followed by a word that means something to you. E.g.
  `[~tag]`.

  ```
  Application [~rest]:
  ```

  In the above example, `rest` is a pattern.

- Key-Value pair

  As the name suggests, you can associate some data with your application or an
  endpoint.

```
Application [version="1.1"]:
```

The value can be a string `"foo"`, an array of strings `["foo", "bar"]`, array
of array of strings `[["foo"], ["bar"]]`

A complete example:

```
BizApp [version="1.1", clients=["web", "daemon"]]:
```

#### Reserved Attributes

Sysl defines some internal attributes that you can use to customize the look of
your diagrams.

- Changing Application icons in Sequence Diagram

  The default icon for the app is a `circle with an arrow`. You can change this
  icon to:

  - human - `App [~human]:`
  - database - `DataBase [~db]`
  - External App - `IdentityProvider [~external]` - In an enterprise system, you
    might have some external third-party system that your app might interact
    with. Mark an app as `~external` and sysl will place that app at extreme
    right of a sequence diagram.

Complete example that uses the above patterns:

```
User [~human]:
  Check Balance:
    MobileApp <- Login
    MobileApp <- Check Balance

MobileApp [~ui]:
  Login:
    Server <- Login
  Check Balance:
    Server <- Read User Balance

Server:
  Login:
    do input validation
    DB <- Save
    return success or failure

  Read User Balance:
    DB <- Load
    return balance

DB [~db]:
  Save: ...
  Load: ...

Project [seqtitle="Diagram"]:
  Seq:
    User <- Check Balance
```

Here is the result:

![](assets/app_icons-Seq.png)

## Statements

Our `MobileApp` does not have any detail yet on how it behaves. Let's use sysl statements to describe behaviour. Sysl supports following types of statements:

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

See [call.sysl](/assets/call.sysl) for complete example.

Now we have all the ingredients to draw a sequence diagram. Here is one generated by sysl for the above example:

![](/assets/call-Seq.png)

See [Generate Diagrams](../../features/diagrams/) on how to draw sequence and other types of diagrams using sysl.

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

See [if-else.sysl](/assets/if-else.sysl) for complete example.

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

See [for-loop.sysl](/assets/for-loop.sysl) for complete example.

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


### Internal relative file

You have `server.sysl`, `client.sysl` and `deps/deps.sysl`. `server.sysl` and
`client.sysl` files are in the same directory. `deps.sysl` file in the
subdirectory. You can import `server.sysl` and `deps.sysl` files in
`client.sysl` as below:

```sysl
# client.sysl

import server
import deps/deps

Client:
  Login:
    Server <- Login
  Dep:
  	Deps <- Dep
```

```sysl
# server.sysl

Server:
  Login: ...
  Register: ...
```

```sysl
# deps.sysl

Deps:
	Dep: ...
```

### Internal absolute file

If the imported are in the same project but outside of current folder, you must
have at least a common root directory and `import /path/from/root`.

All sysl commands accept `--root` argument. Run `sysl help` for more details.

```sysl
# <root-dir>/servers/server.sysl

Server:
  Login: ...
  Register: ...
```

```sysl
# <root-dir>/clients/client.sysl

import /servers/server

Client:
  Login:
    Server <- Login
```

### External file

Sysl supports importing sysl files via the web using Sysl Modules, which are
based on Go Modules.

The imported sysl repository needs to be initialised with a `go.mod` file(run
`go mod init` under the repository working directory). There's no need to
include go code in the repository. As long as there's a `go.mod` file, the
repository will be treated as a sysl/go module.

Import statements are prefaced with `//` e.g `import //the/external/repo/filepath`

```sysl
# This file is located at github.com/foo/bar/servers/server.sysl

Server:
  Login: ...
  Register: ...
```

```sysl
# This file is located at github.com/your/repo/client.sysl

import //github.com/foo/bar/servers/server

Client:
  Login:
    Server <- Login
```

To disable Sysl Modules, set the environment variable
`SYSL_MODULES=off`.

### Non-sysl file

When you import a sysl file, you can omit the `.sysl` file extension.

To import a non-sysl file like swagger file, you can `import foreign_import_swagger.yaml as com.foo.bar.app ~swagger`.

Valid types include
- ~sysl
- ~swagger
- ~openapi3

## Projects

Most of the changes to your system will be done as part of a well defined
`project` or a `software release`.

`TODO: Elaborate`

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


## Data Models

Sysl supports defining two types of data models, one for your database and other
for your app. You can refer to these models in your app just like any other app.

### Relational Data Model

Relational Data model is the most common way of persisting data in a database.
You can define your data model directly in sysl. `data.sysl`

```
CustomerOrderModel:
  !table Customer:
    customer_id <: int [~pk]

  !table CustomerOrder:
    order_id <: int [~pk]
    customer_id <: Customer.customer_id
```

In the above example:

- `CustomerOrderModel` is a user-defined top-level app that contains definitions
  of various tables or types in your data model.
- Customer.customer_id is a primary key.
- CustomerOrder has a foreign key customer_id which refers to the primary key of
  Customer (i.e. customer_id).

See [data.sysl](/assets/data.sysl) for complete example for Relation model.

![](assets/data-Relational-Model.png)

### Object Model

Define a typical in-memory Object model of an application like so:

```
CustomerOrderModel:
  !type Address:
    line_1 <: string
    city <: string

  !type Customer:
    customer_id <: int
    addresses <: set of Address

  !type CustomerOrder:
    order_id <: int
    customer <: Customer
```

Note:

- `set of Address` - Set is the only collection type

See [data.sysl](/assets/data.sysl) for complete example for Object model.

## Events, publisher and subscriber

Sysl has support for the publisher-subscriber model. In the example below,
`UserService` has `RegisterEvent` that is subscribed by `EmailNotifier` and
`SmsNotifier` applications.

```
UserService:
  <-> RegisterEvent: ...

  Register:
    do registration
    . <- RegisterEvent

EmailNotifier:
    UserService -> RegisterEvent:
        EmailNotifier got the RegisterEvent

SmsNotifier:
    UserService -> RegisterEvent:
        SmsNotifier got RegisterEvent
```

![](assets/events-Seq.png)

## Collector

Project files use the Collector syntax to add additional layer of information on
top of the definitions. Best example of this is where you want to capture how an
API evolved over time.

E.g. `server.sysl`

```
Server:
  Login:
    do input validation
    DB <- Save
    return success or failure

  Read User Balance:
    DB <- Load
    return balance
```

Say, in `project-1.sysl` you modified the call to `DB <- Save`, you can capture
this information like so:

```
Server [status="modified"]:
  .. * <- *:
    Login [status="modified"]
    DB <- Save [status="modified"]
```

`.. * <- *` is the way to tell sysl to take an existing definition of the
application and `merge` the attributes from:

- endpoints
- calls

Then for diagrams related to Project-1, you can show this information in the
diagrams which endpoints are getting reused and which ones are new. This
requires usage of `epfmt`, where you can use the value of the variable `status`
like so `%(@status)`. See [Attributes](#attributes) for more details.

By creating separate file for each project, you will always be able to recreate
necessary documentation at that point in time. This can help you answer
questions like `When or why was this introduced?`.
