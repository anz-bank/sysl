# Sysl arr.ai (WIP)

This package contains experimental [arr.ai](https://github.com/arr-ai/arrai) code that operates on Sysl models. It seeks to:

* Demonstrate how to use arr.ai to work with Sysl models.
* Replicate `sysl`'s various commands in a more concise, functional style.
* Introduce new capabilities such as relational modelling and data flow analysis.

Each script has an accompanying Markdown file describing the logic and output.

## Usage

To keep things simple and consistent, all scripts currently work with the same `model.sysl` spec. This spec will evolve over time to contain the full range of Sysl features and patterns to stress test the scripts.

To run all the scripts at once, run:

```bash
make
```

To run a particular script, run `make` for it's output target. For example, to generate the SVG output of the `data_diagram.arrai` script, run:

```bash
make data_diagram.svg
```

For diagram targets (`svg` and `png`), `make` will use the arr.ai script to produce the PlantUML output, then attempt to render it with a PlantUML server and save the result.

You can of course run the scripts directly as well. For example:

```bash
arrai run data_diagram.arrai > out/data_diagram.puml
```

<!-- TODO(ladeo): Generate these Markdown files from the arr.ai sources. -->

## Development

### Loading the model

Many of the scripts will start by importing the default `model.sysl` model like so:

```arrai
let model = //{./load_model};
...
```

This actually loads the proto-encoded `out/model.pb` file, so if `model.sysl` is changed, `out/model.pb` must be regenerated like so:

```bash
sysl protobuf --mode=pb model.sysl > out/model.pb
```

If an extra command line arg is provided to `arrai run`, the `load_model` script treats it as a path to a `.pb`-encoded Sysl model to use instead.

### Environment

To develop additional arr.ai scripts, you're encouraged to use Visual Studio Code with the `arrai` extension.

To speed up the edit/refresh cycle, install the [Save and Run](https://github.com/wk-j/vscode-save-and-run) extension and configure it with:

```json
"saveAndRun": {
    "commands": [
        {
            "match": "\\.arrai$",
            "cmd": "echo; arrai run ${file}",
            "useShortcut": false,
            "silent": false
        }
    ]
}
```

This will make VS Code automatically run the script every time you save it (after printing a newline for clarity).

For more comprehensive testing, you can also run `make` on every save of every file in this directory.
