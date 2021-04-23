---
slug: sysl-catalog
title: Introducing Sysl Catalog
author: Joshua Carpeggiani
author_title: Sysl Core Team
author_url: https://github.com/joshcarp
author_image_url: https://avatars3.githubusercontent.com/u/32605850?s=400&v=4
tags: []
---

A Markdown/HTML + Diagram generator for Sysl specifications.

## Objective

The objective of sysl-catalog is to create the most seamless experience for engineers to document their API behaviour, as well as creating a standardised way multiple teams can create documentation whilst gaining mutual benefit from already existing documentation.

## Background

Let’s say that a team wants some diagrams to represent how their services interact with other services. First, the team needs to choose what format to use, then the team needs to decide on where the docs are going to be hosted, and how often they should be updated.

<!--truncate-->

Let’s say that this team (Team A) chooses PlantUML to create sequence, data model, and integration diagrams for their service. They choose to generate to a docs directory on their single repository and use some proto plugins to automate their Markdown generation:

- Protoc-gen-doc to generate a digest of what how their protos are structured
- Protoc-gen-uml to generate plantuml diagrams to generate diagrams for the needed diagrams
- Manually written Markdown to describe how their services interact with different services
- Manually written sequence diagrams to describe which dependencies are called
  This works fine; the team has somewhat automated their documentation workflow, with some manual parts.

A couple of months pass and now there’s another team which relies on team A’s service heavily. They are releasing soon and need to create release documentation; so they decide to use the same method that Team A is using.

Now there’s a problem; there are two teams with two separate sets of documentation. Some of it is manual and some of it is automated. This can cause problems for multiple reasons:

- Because there’s no synchronisation between Team A and B’s documentation, a specific change by Team A likely won’t show in Team B’s documentation
- If manually written documentation isn't repeated then the representation of their dependency is limited by a hyperlink to Team A's documentation without fully integrating
- The decoupling of documentation and code means they will likely drift apart over time
  This is what sysl-catalog is trying to solve

sysl-catalog uses the Sysl language as an intermediary between different formats to be able to generate different views of how services work.

## What is Sysl?

Sysl is a “system specification language”; think of it like swagger or protos, but a much higher level, and with the ability to represent not only types, applications and endpoints, but interactions between those applications and endpoints; it plans to define what the code does itself.

## What does sysl-catalog do?

Sysl-catalog is just a static site generator.

Sysl-catalog parses a sysl file (with the .sysl extension) and represents it in a visual form; It can represent endpoints (in sequence diagrams) request/response types or database tables, as well as integration diagrams.

It uses go's `text/template` to do this, and if any addition is needed to be made, custom templates can be used (see templates for examples)
TODO: Fix this link
