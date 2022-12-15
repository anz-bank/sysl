import { readFile } from "fs/promises";
import "reflect-metadata";
import { walk } from "../common/iterate";
import { Location } from "../common/location";
import { PbDocumentModel } from "../pbModel/model";
import { Application } from "./application";
import { ElementKind, ElementRef, IRenderable } from "./common";
import { Element } from "./element";
import path from "path";
import fs from "fs";

export type ImportParams = {
    filePath: string;
    locations: Location[];
    appAlias: string;
};

export class Import {
    filePath: string;
    locations: Location[];
    appAlias?: string;

    constructor({ filePath, locations, appAlias }: ImportParams) {
        this.filePath = filePath;
        this.locations = locations;
        this.appAlias = appAlias ?? undefined;
    }

    toSysl(): string {
        return `import ${this.filePath}${this.appAlias ? ` as ${this.appAlias}` : ""}`;
    }
}

export type ModelParams = {
    header?: string | undefined;
    imports?: Import[];
    apps?: Application[];
    locations?: Location[];
    syslRoot?: string;
};

export class Model implements IRenderable {
    header: string | undefined;
    imports: Import[];
    apps: Application[];
    locations: Location[];
    syslRoot: string;

    constructor({ header, imports, apps, locations, syslRoot: rootPath }: ModelParams = {}) {
        this.header = header;
        this.imports = imports ?? [];
        this.apps = apps ?? [];
        this.locations = locations ?? [];
        this.syslRoot = rootPath ?? ".";
        walk(this, {}); // Attach model and parent
    }

    /**
     * Parses Sysl text file into a {@link Model}.
     * @param cwdPath The file CWD-based path that contains the Sysl text to parse.
     * @param maxImportDepth Optional. The maximum depth to follow import statements, where 0 means unlimited depth,
     * 1 means to remain at the level of the supplied file (do not to follow imports at all), 2 means to delve one
     * deeper (supplied file plus one extra depth), etc. If not specified, defaults to allowing unlimited depth.
     * Limiting the depth may significantly improve performance, especially if it causes remote imports to be skipped.
     * @returns A {@link Model} representing the supplied Sysl file.
     */
    static async fromFile(cwdPath: string, maxImportDepth: number = 0): Promise<Model> {
        const syslText = (await readFile(cwdPath)).toString();
        const syslRoot = this.findRoot(cwdPath, ".sysl") || this.findRoot(cwdPath, ".git") || path.resolve(cwdPath);
        return this.fromText(syslText, cwdPath, syslRoot, maxImportDepth);
    }

    /**
     * Parses Sysl text into a {@link Model}.
     * @param syslText The Sysl text to parse.
     * @param syslFilePath Optional. The file path where the Sysl text came from, used to populate the {@link Location}
     * information. If not specified, "untitled.sysl" is used.
     * @param maxImportDepth Optional. The maximum depth to follow import statements, where 0 means unlimited depth,
     * 1 means to remain at the level of the supplied file (do not to follow imports at all), 2 means to delve one
     * deeper (supplied file plus one extra depth), etc. If not specified, defaults to allowing unlimited depth.
     * Limiting the depth may significantly improve performance, especially if it causes remote imports to be skipped.
     * @returns A {@link Model} representing the supplied Sysl text.
     */
    static async fromText(
        syslText: string,
        syslFilePath: string = "untitled.sysl",
        syslRoot: string = ".",
        maxImportDepth: number = 0
    ): Promise<Model> {
        // TODO: Improve performance by only reading the first part of the file
        const lines = syslText.split(/\r?\n/);
        const until = lines.findIndex(l => !l.startsWith("#"));
        const header = lines.slice(0, until).join("\n");

        const pb = await PbDocumentModel.fromText(syslText, syslFilePath, maxImportDepth);

        let newModel = pb.toModel();
        newModel.syslRoot = syslRoot;
        if (header) {
            newModel.header = header;
        }

        return newModel;
    }

    /**
     * Returns a new model with only the apps that were defined in the specified source file. Non-app elements that
     * came from other source and were merged into apps from the specified source file will not be filtered out.
     * @param cwdPath The CWD-based file by which to filter apps.
     * @returns A clone of the current {@link Model} instance which only includes apps that were specified by the
     * supplied source file.
     */
    filterByFile(cwdPath: string): Model {
        cwdPath = this.convertSyslPath(cwdPath);
        return new Model({
            imports: this.imports.filter(i => cwdPath.includes(i.locations[0]?.file!)),
            locations: this.locations,
            apps: this.apps.filter(a => cwdPath.includes(a.locations[0]?.file!)),
            header: this.header,
        });
    }

    /**
     * Finds an element in the model that is specified by the supplied {@link ElementRef}, or returns `undefined` if
     * no such element was found.
     * @param ref The element reference used to locate the element in the model.
     * @returns An {@link Element} that corresponds to the supplied {@link ElementRef}, or `undefined` if not found.
     */
    findElement(ref: ElementRef): Element | undefined {
        const app = this.apps.find(a => ref.appsEqual(a.toRef()));
        if (!app || ref.kind == ElementKind.App) return app;

        const type = app.types.find(t => ref.typesEqual(t.toRef()));
        if (!type || ref.kind == ElementKind.Type) return type;

        return type.children.find(f => ref.equals(f.toRef()));
    }

    /**
     * Renders this model as a Sysl file. If this model contains elements from multiple source files (e.g. if it had
     * import statements), then all elements from all source files will be rendered as a single file. To only render
     * elements for a specific file, first call {@link filterByFile}.
     * @returns A string containing the Sysl textual representation of this model.
     */
    toSysl(): string {
        let sysl = "";

        if (this.header) {
            sysl = this.header + "\n\n";
        }

        if (this.imports.length) {
            sysl += this.imports.map(i => i.toSysl()).join("\n") + "\n\n";
        }

        sysl += this.apps.map(a => `${a.toSysl()}`).join("\n\n");

        sysl += "\n";

        return sysl;
    }

    /**
     * Converts a CWD-based path to a Sysl path that is relative to the {@link syslRoot}, which {@link Location} uses to
     * specify all paths.
     *
     * @param cwdPath A CWD-based path to convert.
     * @param cwd Optional. The CWD to assume, or if unspecified, uses the real CWD.
     * @returns A path relative to the sysl root
     * @throws {@link Error} Thrown when the specified path is outside the {@link syslRoot}, which is not allowed.
     */
    public convertSyslPath(cwdPath: string, cwd: string = process.cwd()) {
        const resolvedPath = path.resolve(cwd, cwdPath);
        const syslRelativePath = path.relative(this.syslRoot, resolvedPath);
        if (syslRelativePath.startsWith("..")) {
            const pathDescription = resolvedPath == cwdPath ? "" : ` (resolved to '${resolvedPath}')`;
            throw new Error(
                `The provided path '${cwdPath}'${pathDescription} is outside the Sysl root path of ${this.syslRoot} so it cannot be converted to a Sysl path.`
            );
        }
        return syslRelativePath;
    }

    private static findRoot(filePath: string, sentinelName: string): string | undefined {
        var root = path.dirname(path.resolve(filePath));
        while (true) {
            if ([...fs.readdirSync(root)].some(f => f == sentinelName)) return root;
            const parent = path.resolve(path.join(root, ".."));
            if (parent == root) return undefined;
            root = parent;
        }
    }
}
