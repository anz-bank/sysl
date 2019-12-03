---
title: "Diagrams"
date: 2019
weight: 50

---


## Generate Diagrams

Once your design is complete, its time to get some output from Sysl. Sysl supports generating diagrams of following types:
  * Sequence Diagram
  * Integration Diagram
  * Data Diagram

Sysl aims to generate code and documentation from only one source of truth i.e. `.sysl` files.

### Sequence Diagrams
You can generate the Sequence Diagram using the following command:

```bash
sysl sd -o 'call-login-sequence.png' -s 'MobileApp <- Login' call.sysl
```
You can omit the the `.sysl` and sysl will pickup the correct file.
```bash
sysl sd -o 'call-login-sequence.png' -s 'MobileApp <- Login' call
```

Here is the output that you should see:

![](//assets/call-Seq.png)

See [/assets/call.sysl](/assets/call.sysl) for complete example.

#### How sysl generates sequence diagram?
Let's breakdown the `sd` aka `sequence diagram` command:
```bash
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

### Integration Diagram
`TODO`

See: run `sysl ints -h` for more details.

### Data Diagram
See [Data Models](#data-models) on types of data models and how to render them.
