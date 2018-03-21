---
title: "Quick start"
description: "Create your first diagrams from sysl demo files in minutes."
date: 2018-02-27T15:51:21+11:00
weight: 1
draft: false
bref: "Create your first diagram from a Sysl specification"
toc: false
---

Sysl consists of two executables: `sysl` and `reljam`.</br> The **Sy**stem **S**pecification **L**anguage `sysl` is mainly concerned with diagram creation whereas the **Rel**ational **Ja**va **M**odel program `reljam` generates different types of source code output.

Install Sysl
------------
Windows users can run standalone Sysl executables whereas users of other operating systems need to work with either Python or Docker.

**Windows users**, download the `sysl-bundle-windows.zip`, containing ``sysl.exe`` and ``reljam.exe``, from the  [Sysl release page](https://github.com/anz-bank/sysl/releases).

**Python users**, install Sysl with [Python 2.7](https://www.python.org/downloads/) and `pip install sysl`.
If your specific environment causes problems you might find our [environment guide](/docs/environment) helpful.

**Docker users**, pull the Sysl image with `docker pull anzbank/sysl` and tag it with
`docker tag anzbank/sysl sysl`.

For more details on installation, please refer to the [installation documentation](/docs/installation). To learn more about the command line options and try running `sysl --help` or read the [command line documentation](/docs/commandline).

Generate text-based output
--------------------------
Download and save the [petshop.sysl](https://raw.githubusercontent.com/anz-bank/sysl/master/demo/petshop/petshop.sysl) from the Sysl repository.
In the commands below `/petshop` refers to this `petshop.sysl` file in the current working directory.


Windows exe and Python users create **Text Protobuf** and **Java** ouptput with:

	sysl textpb -o out/petshop.txt /petshop
	reljam model /petshop PetShopModel

Docker users run:

	docker run sysl sysl textpb -o out/petshop.txt /petshop
	docker run sysl reljam model /petshop PetShopModel

Generate a sequence diagram
---------------------------
Download and save the [Bank demo files](https://github.com/anz-bank/sysl/tree/master/demo/bank) from the Sysl repository:

* [bank.sysl](https://raw.githubusercontent.com/anz-bank/sysl/master/demo/bank/bank.sysl)
* [project.sysl](https://raw.githubusercontent.com/anz-bank/sysl/master/demo/bank/project.sysl)

Sysl uses [PlantUML](http://plantuml.com/) for its diagram creation. You therefore need to set the `SYSL_PLANTUML` environment variable or pass a valid PlantUML service URL via `-p`

	sysl sd -p "http://www.plantuml.com/plantuml" -o "out/sd-%(epname).png" /project -a "Bank :: Sequences" -v

You can also [install PlantUML](http://plantuml.com/starting) locally and update the environment variable accordingly, e.g. `SYSL_PLANTUML=http://localhost:8080/plantuml`.

The generated output should look like this:

![Sequence diagram](/img/sysl/bank-sd.svg)
