---
title: "Diagrams"
date: 2018-02-28T10:11:24+11:00
description: "Diagram generation from Sysl"
weight: 40
bref: "Diagram generation from Sysl"
topic: "Diagrams"
layout: "single"
toc: true
---

Sysl can generate various different diagrams. Sysl aims to generate code and documentation from only one source of truth i.e. .sysl files.

## Sequence diagrams

For an example refer to [sequence-diagrams](/docs/byexample/sequence-diagrams/)
For details on the command, refer to [sequence](/docs/commands/sd)

### Format Arguments

The default diagram by default only shows the data type that is returned by an endpoint. You can instruct `sysl` to show the arguments to your endpoint in a sequence diagram.

Command:

`sysl sd -o 'call-login-sequence.png' --epfmt '%(epname) %(args)' -s 'MobileApp <- Login' /assets/call.sysl -v call-login-sequence.png`
See [/assets/args.sysl](/assets/args.sysl) for complete example.

![](/assets/args-Seq.png)

A bit more explanation is required regarding `epname` and `args` keywords that are used in `epfmt` command line argument. See section on [Attributes](#epfmt) below.

### Using attributes in appfmt and epfmt

`appfmt` and `epfmt` (app and endpoint format respectively) can be passed to
`sd`, `ints` commands. They control how the application or endpoint name is
rendered as text. There default value is `%(appname)` and `%(epname)`
respectively. These internal attributes are:

    * appname - short name of the application
    * epname - short name of the endpoint
    * eplongname - Long quoted name of the endpoint.
    * controls - controls defined on your endpoint

Complete example:

```
App "Descriptive Long Application name":
  Endpoint-1 "Descriptive Long name for Endpoint 1":
    ...
  Endpoint-2 "Descriptive Long name for Endpoint 2":
    ...
```

Where:

- appname - App
- epname - Endpoint-1 or Endpoint-2
- eplongname - "Descriptive Long name for Endpoint 1" or "Descriptive Long name
  for Endpoint 2"

You can also refer to the attributes that you added by using `[]` or the
Collector syntax.

#### Using user defined attributes in fmt

You can use your attributes in `epfmt` or `appfmt` arguments in the following
ways:

- `%(@attrib_name)` : use `@` to refer to `attrib_name`.
- `%(@attrib_name? yes_stmt | no_stmt)`: use `?` to test for existence of value.
  This is ternary operator, which allows you to to execute `yes_stmt` or
  `no_stmt` depending on the result.
- `%(@attrib_name=='some_value'? yes_stmt | no_stmt)` : compare attrib's value
  to some constant.
- `%(@attrib_name=='some_value'? yes_stmt | @attrib_name=='some_other_value'? | ...)` : nested checks.

Now, `stmt` can be any of the following types:

- plain-text: will be copied as-is to the output
- `<color red>text or %(attrib_name)</color>`: use html like syntax to color the
  output.

Example:

```html
appfmt="%(@status?<color red>%(appname)</color>|%(appname))" epfmt="%(@status?
<color green>%(epname)</color>|%(epname))"
```

See [attribs.sysl](assets/attribs.sysl) for complete example. Notice how
`appfmt` and `epfmt` use `%(@status)`.

![](assets/attribs-Seq.png)


## Integration diagrams

For an example refer to [integration-diagrams](/docs/byexample/integration-diagrams/)
For details on the command, refer to [integrations](/docs/commands/integrations)

## Datamodel diagrams

For an example refer to [data-model-diagrams](/docs/byexample/data-model-diagrams/)
For details on the command, refer to [datamodel](/docs/commands/datamodel)

