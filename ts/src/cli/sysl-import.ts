#!/usr/bin/env node -r ts-node/register

import { program } from "@commander-js/extra-typings";
import fs from "fs/promises";
import { ImportOptions, importAndMerge } from "../import";

const opts: ImportOptions = program
    .name("import")
    .requiredOption("-i, --input <path>", "path of file to import")
    .option("-a, --app-name <name>", "name of the Sysl app to define in Sysl model.")
    .option("-o, --output <path>", "path of file to write the imported sysl, writes to stdout when not specified")
    .option(
        "-f, --format <name>",
        `format of the input filename. Formats are autodetected where possible, but this can force the use of a particular importer.`
    )
    .option(
        "-I, --import-paths <...paths>",
        "comma-separated list of paths used to resolve imports in the input file. Currently only used for protobuf input."
    )
    .option("-s, --shallow", "does shallow parsing of input and excludes definitions imported by the specification")
    .parse()
    .opts();

importAndMerge(opts).then(async result => {
    if (opts.output) {
        await fs.writeFile(opts.output, result.output);
    } else {
        process.stdout.write(result.output);
    }
});
