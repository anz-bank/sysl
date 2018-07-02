  ---
title: "Command line"
description: "Learn to use the sysl and reljam command line tools."
date: 2018-02-27T15:55:46+11:00
weight: 20
draft: false
bref: "Sysl and reljam command line arguments"
toc: true
---

Sysl consists of two executables: `sysl` and `reljam`.</br> The **Sy**stem **S**pecification **L**anguage `sysl` is mainly concerned with diagram creation whereas the **Rel**ational **Ja**va **M**odel program `reljam` generates different types of source code output.

Both `sysl` and `reljam` comprise several **sub-commands** for different types of output generation, for example `sysl sd`, `reljam model`. Refer to the [output formats documentation](/docs/outputs) for a complete list.

`sysl` and `reljam` have some shared **global options**:

  *  `--help` and `<sub-command> --help` for more help on the commandline
  *  `--version` prints version information
  *  `--trace` debug information useful for issue reporting
  *  `--root DIR` root directory for Sysl modules and/or files, default `.`
  *  `--out OUT` output directory, file name or file pattern

Input
-----
Specify one or more `*.sysl` input files relative to the `--root` directory. The default root directory is `.`, the current working directory.

```
optional arguments:
  -h, --help            show this help message and exit
  --no-validations, --nv
                        suppress validations
  --root ROOT, -r ROOT  sysl root directory for input files (default: .)
  --version, -v         show version number (semver.org standard)
  --trace, -t
```

Output
------
The output is specified with the `--out` or `-o` flag. It could be a file name, a directory name or a file name pattern depending on the sub-command used. Here are three different example:

* **File:** <br/>`sysl data -o output.png -j Project hello-world.sysl`
* **Directory:** <br/>`reljam -o output_directory model hello-world.sysl HelloWorld`
* **Pattern:** <br/>`sysl data -o "%(appname)-%(epname).png" -j Project /hello"` This command generates a png file for each relevant application name `%(appname)` and endpoint `%(epname)` combination.

There is a subtle difference in the usage of the output option for `sysl` as opposed to `reljam`: `--out` has to specified after the sub-command in `sysl` and before the sub-command in `reljam`:

    sysl data -o output.png -j Project hello-world.sysl
    reljam -o output_directory model hello-world.sysl HelloWorld

Sysl commands
-------------
`sysl` sub-comands are `pb`, `textpb`, `data`, `ints`, `sd`. Find out more about each sub-command with `sysl <sub-command> --help` and in the [output formats documentation](/docs/outputs). You can also find several standalone [examples on GitHub](https://github.com/anz-bank/sysl/tree/master/demo/simple) and generate the following outputs:

  * Data Model diagram:
    - `sysl data -o out.png -j Project /sysl-data`
  * Integration diagrams:
    - `sysl ints -o "out.png" /sysl-ints --project Project`
    - `sysl ints -o "out.png" /sysl-ints --project Project --epa`
  * Sequence diagram:
    - `sysl sd -a Project -o "%(epname)".png /simple-sd`

Reljam commands
---------------
`reljam` sub-comands are `model`, `facade`, `view`, `xsd`, `swagger`, `spring-rest-service`. Here too, you can find out more about each sub-command with `reljam <sub-command> --help` and in the [output formats documentation](/docs/outputs). You can also find several standalone [examples on GitHub](https://github.com/anz-bank/sysl/tree/master/demo/simple) and a few selected examples below:

  * Java Model source code:
    - `reljam model reljam-model.sysl HelloWorld`
  * XSD source code:
    - `reljam xsd reljam-xsd.sysl`

