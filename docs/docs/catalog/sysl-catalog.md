---
id: sysl-catalog
title: Sysl Catalog CLI Reference
sidebar_label: Overview
---

import useBaseUrl from '@docusaurus/useBaseUrl';

`sysl-catalog` is a standalone CLI tool that allows you to generate documentation from your Sysl specifications. It currently supports the following output formats:

- Markdown
- HTML (using [GoldMark](https://github.com/yuin/goldmark))

To install it, [follow the instructions here](sysl-catalog-install.md)

## Example Commands

### Output default Markdown

```bash
sysl-catalog -o=docs/ filename.sysl
```

### Output default HTML

```bash
sysl-catalog -o=docs/ --type=html filename.sysl
```

### Run with custom templates

```bash
sysl-catalog --templates=<filename1.tmpl>,<filename2.tmpl> filename.sysl
```

With this the first template will be executed first, then the second

### Run in server mode

`sysl-catalog` comes with a `serve` mode which will serve on port `:6900` by default

```bash
sysl-catalog --serve <input.sysl>
```

This will start a server and filewatchers to watch the input file and its directories recursively, and any changes will automatically update.

<img alt="server mode" src={useBaseUrl('img/sysl/sysl-catalog-serve.gif')} />

### Generate ReDoc files

```bash
sysl-catalog --redoc filename.sysl
```

This generates a [ReDoc](https://github.com/Redocly/redoc) page that serves the original JSON or YAML OpenAPI spec on GitHub. It currently only supports spec files located in the same repo, and must be run in a Git repo (so that the remote URL can be retrieved using `git`).

### Run in server mode without css/rendered images

```bash
sysl-catalog --serve --noCSS filename.sysl
```

This is useful for rendering raw markdown

<img alt="server mode raw" src={useBaseUrl('img/sysl/sysl-catalog-raw-markdown.png')} />

### Run server with custom template

```bash
sysl-catalog --serve --templates=<filename1.tmpl>,<filename2.tmpl> filename.sysl
```

<img alt="server mode custom" src={useBaseUrl('img/sysl/sysl-catalog-custom-template.png')} />
