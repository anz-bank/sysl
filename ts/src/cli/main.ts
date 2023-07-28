#!/usr/bin/env node -r ts-node/register

// To run this TypeScript source, use:
// $ node -r ts-node/register src/cli/main.ts [command] [options]

import { program } from "@commander-js/extra-typings";
import { spawnSysl } from "../common/spawn";

program
    .executableDir(__dirname)
    .name("sysl")
    .command("import [options]", "import a file into Sysl")
    .on("command:*", async command => process.stdout.write(await spawnSysl(command)))
    .parse();
