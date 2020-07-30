# Sysl arr.ai (WIP)

This package contains experimental [arr.ai](https://github.com/arr-ai/arrai) code that operates on Sysl models. It seeks to:

* Demonstrate how to use arr.ai to work with Sysl models.
* Replicate `sysl`'s various commands in a more concise, functional style.
* Introduce new capabilities such as relational modelling and data flow analysis.

## Usage

To keep things simple and consistent, all scripts currently work with the same `model.sysl` spec. This spec will evolve over time to contain the full range of Sysl features and patterns to stress test the scripts.

To run the arr.ai scripts, install arr.ai and execute `arrai run [file]`. For example:

```bash
$ arrai run integration_diagram.arrai
```

Each script has an accompanying Markdown file describing the logic and output.

To run all the scripts at once, run `make`. The Makefile commands also generate SVGs for the scripts that produce diagrams.

<!-- TODO(ladeo): Generate these Markdown files from the arr.ai sources. -->

## Development

Many of the scripts will start by importing `model.sysl` (assuming it's in the current working directory) like so:

```arrai
let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('sysl.pb'));
let module = //encoding.proto.decode(sysl, 'Module' , //os.file('model.pb'));
...
```

If `model.sysl` is changed, `model.pb` must be regenerated like so:

```bash
$ sysl pb --mode=pb -o model.pb model.sysl
```

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

This will make VS Code automatically run the script every time you save it.
