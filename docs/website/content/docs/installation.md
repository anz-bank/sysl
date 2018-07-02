---
title: "Installation"
description: "Sysl can be installed on Windows, MacOS and Linux - follow this guide."
date: 2018-02-27T15:51:27+11:00
weight: 10
draft: false
bref: "Sysl can be installed on Windows, MacOS and Linux - follow this guide"
toc: true
---
Sysl is a CLI (Command Line Interface), split between two executables: `sysl` and `reljam` (Relational Java Model generator).
Windows users can download and run standalone Sysl executables whereas users of other operating systems need to work with either Python or Docker. `sysl --version` will display the currently installed version of Sysl.

Windows Exe
-----------
Windows users can download the `sysl-bundle-windows.zip`, containing the command line tools `sysl.exe` and `reljam.exe`, from the [Sysl release page](https://github.com/anz-bank/sysl/releases>).

Users on other operating systems need to work with Python or Docker.

Python
------
Install [Python 2.7](https://www.python.org/downloads/).
If your specific environment causes problems you might find our [guide](/docs/environment) helpful.

Install Sysl with

	> pip install sysl

Now you can execute Sysl on the command line with

	> sysl   textpb -o out/petshop.txt /demo/petshop/petshop
	> reljam model /demo/petshop/petshop PetShopModel

See `sysl --help` and `reljam --help` for more options.

Docker
------
Install [Docker](https://docs.docker.com/install/) and pull the Docker Image with

	> docker pull anzbank/sysl

Consider tagging the docker image to make commands shorter

	> docker tag anzbank/sysl sysl

Try the following commands

	> docker run sysl
	> docker run sysl sysl -h
	> docker run sysl reljam -h

See [Sysl Image on Docker Hub](https://hub.docker.com/r/anzbank/sysl/) for more details.
