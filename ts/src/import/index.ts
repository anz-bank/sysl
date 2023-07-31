import fs from "fs/promises";
import path from "path";
import { spawnSysl } from "../common/spawn";
import { Action, Alt, Cond, Foreach, Group, Loop, LoopN, Model, Return, Statement } from "../model";

export type ImportOptions = {
    input: string;
    appName?: string;
    output?: string;
    format?: string;
    importPaths?: string;
    shallow?: boolean;
};

export type ImportResult = {
    /** The content of the imported Sysl spec. */
    model: Model;
    /** The path to the existing Sysl spec that was loaded and merged, if any. */
    existingPath?: string;
};

export async function importAndMerge(opts: ImportOptions): Promise<ImportResult> {
    const newMod = importNew(opts);
    const existingPath = opts.output ?? opts.input.replace(path.extname(opts.input), "") + ".sysl";
    const oldMod = await loadExisting(existingPath);
    if (oldMod) mergeExisting(await newMod, oldMod);

    return {
        model: await newMod,
        existingPath,
    };
}

async function importNew(opts: ImportOptions): Promise<Model> {
    const passthroughArgs = Object.entries(opts)
        // Drop any given `output` so we get the output through stdout.
        // Drop `shallow` since it's a boolean flag
        .filter(([k]) => k != "output" && k != "shallow")
        .map(([k, v]) => `--${k.replace(/([A-Z])/g, "-$1").toLowerCase()}=${v}`);
    if (opts.shallow) passthroughArgs.push("--shallow");
    return Model.fromText((await spawnSysl(["import", ...passthroughArgs])).toString());
}

async function loadExisting(existingPath: string): Promise<Model | undefined> {
    // prettier-ignore
    const exists = await fs.open(existingPath).then(() => true).catch(() => false);
    return exists ? Model.fromFile(existingPath) : undefined;
}

/** Merges aspects of `oldModel` that should be retained when reimported into `newModel`. */
export function mergeExisting(newModel: Model, oldModel: Model): void {
    mergeEndpoints(newModel, oldModel);
    mergeImports(newModel, oldModel);
    mergeHeader(newModel, oldModel);
}

/** Merges endpoint statements from `oldModel` into `newModel`. */
export function mergeEndpoints(newModel: Model, oldModel: Model): void {
    newModel.apps.forEach(app => {
        app.endpoints.forEach(ep => {
            const oldEp = oldModel.findApp(app.toRef())?.endpoints.find(e => e.name == ep.name);
            if (!oldEp) return;
            if (oldEp.statements.length == 1 && isPlaceholder(oldEp.statements[0])) return;

            const oldRets = oldEp.statements.flatMap(flattenStatement).filter(s => s.value instanceof Return);
            const oldPayloads = new Set(oldRets.map(r => (r.value as Return).payload));
            const newRets = ep.statements.flatMap(flattenStatement).filter(s => s.value instanceof Return);
            const addRets = newRets.filter(r => !oldPayloads.has((r.value as Return).payload));
            if (oldRets.length + addRets.length != newRets.length) {
                console.warn(`${app.toRef().toSysl()}.${ep.name} has a return not present in the new model`);
            }
            ep.statements = [...oldEp.statements, ...addRets];
        });
    });
}

function isPlaceholder(s: Statement): boolean {
    return s.value instanceof Action && s.value.action == "...";
}

type ParentStatementType = Cond | Loop | LoopN | Foreach | Group;
const parentStatementClasses = [Cond, Loop, LoopN, Foreach, Group];

export function flattenStatement(s: Statement): Statement[] {
    if (parentStatementClasses.some(Cls => s.value instanceof Cls)) {
        return [s, ...(s.value as ParentStatementType).stmt.flatMap(flattenStatement)];
    }
    if (s.value instanceof Alt) {
        const children = s.value.choices.flatMap(c => c.stmt);
        return [s, ...children.flatMap(flattenStatement)];
    }
    return [s];
}

/** Adds any unique import statements from `oldModel` to `newModel`. */
export function mergeImports(newModel: Model, oldModel: Model): void {
    const newImportPaths = new Set(newModel.imports.map(i => i.filePath));
    newModel.imports.push(...oldModel.imports.filter(i => !newImportPaths.has(i.filePath)));
}

/** Replaces the `DO NOT EDIT` header added by most importers with more precise guidance. */
export function mergeHeader(newModel: Model, _oldModel: Model): void {
    newModel.header = newModel.header?.replace(
        "Code generated by Sysl. DO NOT EDIT.",
        "Code generated by Sysl. CAUTION: Edit only imports and statements."
    );
}
