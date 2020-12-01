---
id: sysl-catalog-cmd
title: Sysl Catalog CLI Reference
sidebar_label: sysl-catalog
---

:::info
This command requires the `SYSL_PLANTUML` environment variable to be set or passed in as a flag. [Follow the PlantUML instructions](../plantuml.md) for more details.
:::

## Summary

`sysl-catalog` is a standalone CLI tool that allows you to generate documentation from your Sysl specifications. It currently supports the following output formats:

- Markdown
- HTML (using [GoldMark](https://github.com/yuin/goldmark))

## Usage

```bash
usage: sysl-catalog [<flags>] <input>
```

Args:

- `<input>` Input sysl file to generate documentation for

## Optional Flags

- `--help` Show context-sensitive help (also try --help-long and --help-man).
- `--plantuml=PLANTUML` PlantUML service to use
- `-p, --port=":6900"` Port to serve on
- `--type="markdown"` Type of output. Supported Formats: (`markdown`|`html`)
- `-o, --output=OUTPUT` Output directory to generate to
- `-v, --verbose` Verbose logs
- `--templates=TEMPLATES` Custom templates to use, separated by a comma
- `--outputFileName=""` Output file name for pages; `{{.Title}}`
- `--serve` Start an HTTP server and preview documentation
- `--noCSS` Disable adding CSS to served HTML
- `--disableLiveReload` Disable live reload
- `--noImages` Disable images creation
- `--embed` Embed images instead of creating SVGs
- `--mermaid` Use Mermaid diagrams where possible (not currently supported)
- `--redoc` Generate ReDoc for specs imported from OpenAPI. Must be run on a Git repo.
- `--imageDest=IMAGEDEST` Optional image directory destination (can be outside output)

## Requirements

To see an example of a sysl file used to generate documentation, refer to [demo.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/demo.sysl)

1. `@package` attribute must be specified:

Currently the package name is not inferred from the application name (`MobileApp`), so this needs to be added (`ApplicationPackage`).

```sysl
MobileApp:
    @package = "ApplicationPackage"
    Login(input <: Server.Request):
        Server <- Authenticate
        return ok <: MegaDatabase.Empty
```

1. Application names might need to be prefixed to parameter types if the type is defined in another application, since defined parameters are under scope of the application it is defined in:

```diff
MobileApp:
    @package = "ApplicationPackage"
+    Login(input <: Server.Request):
-    Login(input <: Request):
        Server <- Authenticate
        return ok <: MegaDatabase.Empty
```

3. Add `~ignore` to applications/projects that are to be ignored in the Markdown creation:

```sysl
ThisAppShouldntShow[~ignore]:
    NotMySystem:
        ...
# Or ignore only specific endpoints
ThisAppShouldShow:
    NotMySystem[~ignore]:
        ...
```

## Example

## Troubleshooting

On macOS, if your `launchctl limit maxfiles` setting is too low (e.g 256) you might see the error message "too many open files" when running make.

You can set the current session limit higher with:

```bash
sudo launchctl limit maxfiles 65536 200000
```

And add the following line to your `.bash_profile` or analogous file:

```bash
ulimit -n 65536 200000
```
