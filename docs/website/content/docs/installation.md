---
title: "Installation"
description: "Sysl can be installed on Windows, MacOS and Linux - follow this guide."
date: 2018-02-27T15:51:27+11:00
weight: 10
draft: false
bref: "Sysl can be installed on Windows, MacOS and Linux - follow this guide"
toc: true
---
Sysl is a CLI (Command Line Interface) that excecutes with the `sysl` command. 

Prerequisites
-----------
- [Go](https://golang.org)
- [PlantUML](https://hub.docker.com/r/plantuml/plantuml-server/) server for diagram generation for use if using the [external service](http://www.plantuml.com/plantuml/) is not appropriate 

Linux/Macos
-----------
1. Clone the repo

`git clone https://github.com/anz-bank/sysl.git`

2. Build the binary

`cd sysl/sysl2/sysl && go build`

(note that sysl2 is the new sysl version written in go)

3. To insall it globally copy it to your bin directory:
`cp sysl /usr/local/bin`

Congrats: Now sysl is installed and the fun can begin!

Windows
-----------
1. Install your Preferred linux distro and then go to section "Linux/Macos"