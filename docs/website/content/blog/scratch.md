---
title: "Scratching an itch"
date: 2018-02-15T15:55:46+11:00
draft: false
author: Marcelo Cantos
---
Often we have heard that software products were started when somebody scratched their own itch. And so it was with Marcelo and Sysl, here his story:

Sysl was born as a reaction to the state of architecture as I saw it when I joined ANZ (Australia and New Zealand Banking Group Limited): a morass of documents, spreadsheets and Visio diagrams depicting different parts of the architecture elephant with wildly different representations across platforms and no unified formalism to allow the architecture of the system to be understood and analysed as a whole. Moreover, once complete such designs fossilised rapidly, since there was no formal tie-in between them and the code that was derived from them. Generating code from these designs was an entirely manual process. As a coder by trade (and passion), I saw that none of the extant GUI-focused tooling was doing anything to remedy this, and I even suspected that they were part of the problem, so I decided to experiment with a text-based formalism to describe a system I was designing at the time.

I initially implemented a simple syntax for expressing application endpoints as pseudo-code. I defined the pseudo-code to be quite free-flowing and informal, but with just enough formality to properly describe integrations between applications and behavioural sequences. Coupled with a schema definition language, this gave Sysl enough expressive power to automatically generate all the architectural views required in a solution design, including integration views, sequence diagrams and data model views. (Side note: I don’t believe static diagrams are a viable medium to depict software architecture, and I have always had plans to generate dynamic navigable views from Sysl. A Sysl browser, essentially.)

A second core goal of Sysl was to eliminate the tedious manual translation of software designs into code. I wrote a set of code generators that ingested designs expressed in Sysl and automatically spat out code in various forms. The initial use was to generate data model code for a complex product origination system at ANZ. The project was well into development when I arrived and I found that their data modelling efforts were under considerable strain. As the project progressed, necessary model changes were arising quite frequently, but since the model was expressed as a spreadsheet, a Visio diagram, Java code, XSLT, and other forms, and because most of these forms weren’t tied together in an automated way, changes were taking weeks to absorb and were therefore progressing in deeply overlapping waves of change across multiple assets. I did an initial analysis of the various assets and found that almost none of them exhibited perfect agreement. In fact, since the assets were subject to quite different versioning regimes, it was difficult to even pin down what was comparable to what.

The first action I took was to help the data modellers establish a single point of truth. This was a spreadsheet containing a table of all entities and fields in the data model. I got them to add some extra columns such as _primary key_ and _foreign key_, and thus capture integrity rules that previously existed only in the form of arrows on the Viso diagram. I wrote a translation script that converted the spreadsheet into Sysl’s data modelling sublanguage, and then wrote a code generator, which produced Java code that implementing the Sysl model. This code implemented a relational data model in memory, complete with integrity constraints, some simple querying capabilities, and a very handy set of navigational methods that allowed programmer to follow foreign-key relationships from entity to entity.

The effect of this was manifold:

1. The gap between the various representations instantly and forever disappeared, since all changes now flowed automatically from a single source of truth.
2. Developers realised without being told that changes had to flow from upstream into their code.
3. Developers gained a powerful set of data processing capability out of the box, enabling to remove large slabs of hand-written logic from the code base.
4. The time to introduce data model changes went from 2–4 weeks down to 2–4 hours, sometimes even down to minutes.
5. Project leads were much more willing to accepting changes to the data model as the work progressed.

In addition to Java models, I also added JSON/XML reading and writing logic, XSDs, REST API service stubs, and so on, so that none of the code embodying Sysl data models had to be written by hand.

While Sysl has already saved ANZ millions of dollars avoiding wasted time and effort, I don’t feel that it has come into its own just yet. The full promise of Sysl is far from realised, and we continue to work towards a vision of Sysl being the go-to point for expressing and discovering how all systems at ANZ function and connect, and for ensuring that what we design and what we build are (and remain) one and the same.
