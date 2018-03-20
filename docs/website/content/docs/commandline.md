  ---
title: "Command line"
description: "Learn to use sysl and reljam command line and its arguments"
date: 2018-02-27T15:55:46+11:00
weight: 20
draft: false
bref: "Sysl and reljam command line arguments"
toc: true
---

Sysl consists of two executables: `sysl` and `reljam`. The **Sy**stem **S**pecification **L**anguage `sysl` is mainly concerned with diagram creation whereas the **Rel**ational **Ja**va **M**odel program `reljam` generates different types of source code output.

Both `sysl` and `reljam` comprise several **sub-commands** for different types of output generation, for example `sysl sd`, `sysl ints`, etc. For a complete list refer to the [output formats documentation](/docs/outputs).

`sysl` and `reljam` have some shared **global options**:

  *  `--help` and `<subcommand> --help` for more help on the commandline
  *  `--version` prints version information kept in lock-step for both executables
  *  `--trace` debug information useful for issue reporting
  *  `--root` <root> directory to which Sysl modules and/or files are relative (default: `.`)
  *  `--out` location or pattern of the output e.g. `--out sequence_diagram.png`

Input
-----
Specify one or more `*.sysl` input files relative to the `--root` directory (default: `.`) in both `sysl` and `reljam`. You can list your desired input either as relative file paths or with the sysl module notation dropping the extension and adding a leading `/`, for example:

    sysl textpb -o out.txt /hello
    sysl textpb -o out.txt hello.sysl

Output
------
The output is specified with the `--out` or `-o` flag. It could be a file name, a directory name or a file name pattern depending on the subcommand used, for example:

* file: `sysl data -o output.png -j Project hello-world.sysl`
* directory: `reljam -o output_directory model hello-world.sysl HelloWorld`
* pattern: `sysl data -o "%(appname)-%(epname).png" -j Project /hello"` generates a png file for each relevant application name and endpoint combination.

Additionally, in `sysl` and `reljam` the output option needs to be used in a slightly different way: `--out` has to specified after the sub-command in `sysl` and before the sub-command in `reljam`:

    sysl data -o output.png -j Project hello-world.sysl
    reljam -o output_directory model hello-world.sysl HelloWorld

Sysl commands
-------------
`sysl` sub-comands are `pb`, `textpb`, `data`, `ints`, `sd`. Find out more about each subcommand with `sysl <subcommand> --help` and in the [output formats documentation](/docs/outputs). You can also find several standalone [examples on GitHub](https://github.com/anz-bank/sysl/tree/master/demo/simple) and generate the following outputs:

  * Data Model diagram:
    - `sysl data -o out.png -j Project /sysl-data`
  * Integration diagrams:
    - `sysl ints -o "out.png" /sysl-ints --project Project`
    - `sysl ints -o "out.png" /sysl-ints --project Project --epa`
  * Sequence diagram:
    - `sysl sd -a Project -o "%(epname)".png /simple-sd`

Reljam commands
---------------
`reljam` sub-comands are `model`, `facade`, `view`, `xsd`, `swagger`, `spring-rest-service`. Here too, you can find out more about each subcommand with `reljam <subcommand> --help` and in the [output formats documentation](/docs/outputs). You can also find several standalone [examples on GitHub](https://github.com/anz-bank/sysl/tree/master/demo/simple) and generate the following outputs:

