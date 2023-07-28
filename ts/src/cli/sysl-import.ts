#!/usr/bin/env node

import { program } from "@commander-js/extra-typings";
import fs from "fs/promises";
import path from "path";
import { spawnSysl } from "../common/spawn";
import { mergeExisting } from "../import";
import { Model } from "../model";

type ImportOptions = {
    input: string;
    appName?: string | undefined;
    output?: string | undefined;
    format?: string | undefined;
    importPaths?: string | undefined;
    shallow?: true | undefined;
};

const opts = program
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

importAndMerge(opts);

async function importAndMerge(opts: ImportOptions): Promise<void> {
    const newMod = importNew(opts);
    const existingPath = opts.output ?? opts.input.replace(path.extname(opts.input), "") + ".sysl";
    const oldMod = await loadExisting(existingPath);
    if (oldMod) {
        mergeExisting(await newMod, oldMod);
    }

    const out = (await newMod).toSysl();
    if (opts.output) {
        await fs.writeFile(existingPath, out);
    } else {
        process.stdout.write(out);
    }
}

async function importNew(opts: ImportOptions): Promise<Model> {
    const passthroughArgs = Object.entries(opts)
        .filter(([k]) => k != "output") // Drop any given `output` so we get the output through stdout.
        .map(([k, v]) => `--${k.replace(/([A-Z])/g, "-$1").toLowerCase()}=${v}`);
    return Model.fromText((await spawnSysl(["import", ...passthroughArgs])).toString());
}

async function loadExisting(existingPath: string): Promise<Model | undefined> {
    // prettier-ignore
    const exists = await fs.open(existingPath).then(() => true).catch(() => false);
    return exists ? Model.fromFile(existingPath) : undefined;
}
