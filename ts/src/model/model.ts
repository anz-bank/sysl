import { readFile } from "fs/promises";
import "reflect-metadata";
import { Location } from "../common/location";
import { PbDocumentModel } from "../pbModel/model";
import { Application } from "./application";
import { ElementID, ElementKind, ElementRef, IRenderable } from "./common";
import { CloneContext, ICloneable, ModelFilter, ModelFilters } from "./clone";
import { Element } from "./element";
import path from "path";
import fs from "fs";
import { Type } from "./type";
import { Field } from "./field";

export type ImportParams = {
    filePath: string;
    locations?: Location[];
    appAlias?: string;
};

export class Import implements ICloneable {
    filePath: string;
    locations: Location[];
    appAlias?: string;

    constructor({ filePath, locations, appAlias }: ImportParams) {
        this.filePath = filePath;
        this.locations = locations ?? [];
        this.appAlias = appAlias ?? undefined;
    }

    toSysl(): string {
        return `import ${this.filePath}${this.appAlias ? ` as ${this.appAlias}` : ""}`;
    }

    toString(): string {
        return this.filePath;
    }

    clone(_context?: CloneContext): Import {
        return new Import({ filePath: this.filePath, appAlias: this.appAlias });
    }
}

export type ModelParams = {
    header?: string | undefined;
    imports?: Import[];
    apps?: Application[];
    locations?: Location[];
    syslRoot?: string;
};

/** Configures how Sysl files should be parsed. */
export type ParseParams = {
    /**
     * Optional, default depends on method used. Specifies the location of the Sysl root which is stored in
     * {@link Model.syslRoot}. This value is not passed on to the Sysl binary, it is only recorded in the {@link Model}.
     */
    syslRoot?: string;
    /**
     * Optional, default 0. Sets the max import depth. A value of 0 means unlimited depth, a value of 1 means ignore all
     * imports, a value of 2 means follow imports one level deep, etc.
     */
    maxImportDepth?: number;
    /**
     * Optional, default true. When set to true, the latest version of remote imports will always be fetched via Git.
     * When set to false, each remote import will only be fetched if that import has no prior cached copy. If there is,
     * the cached copy will be used and no Git pull will be performed for that import. Specifying false when a previous
     * operation already performed the fetch can greatly improve performance by skipping the remote fetching.
     */
    alwaysFetch?: boolean;
};

export class Model implements IRenderable {
    header: string | undefined;
    imports: Import[];
    apps: Application[];
    locations: Location[];
    syslRoot: string;

    constructor({ header, imports, apps, locations, syslRoot }: ModelParams = {}) {
        this.header = header;
        this.imports = imports ?? [];
        this.apps = apps ?? [];
        this.locations = locations ?? [];
        this.syslRoot = syslRoot ?? ".";
        this.attachSubitems();
    }

    /**
     * Parses Sysl text file into a {@link Model}.
     * @param cwdPath The file CWD-based path that contains the Sysl text to parse.
     * @param params The parameters used to parse Sysl. If {@link ParseParams.syslRoot} is not specified, it is
     * autodetected from {@link cwdPath} using the same logic that the Sysl binary uses.
     * @returns A {@link Model} representing the supplied Sysl file.
     */
    static async fromFile(cwdPath: string, params: ParseParams = {}): Promise<Model> {
        const syslText = (await readFile(cwdPath)).toString();
        params.syslRoot =
            params.syslRoot ??
            this.findRoot(cwdPath, ".sysl") ??
            this.findRoot(cwdPath, ".git") ??
            path.resolve(cwdPath);
        return this.fromText(syslText, cwdPath, params);
    }

    /**
     * Parses Sysl text into a {@link Model}.
     * @param syslText The Sysl text to parse.
     * @param syslFilePath Optional. The file path where the Sysl text came from, used to populate the {@link Location}
     * information. If not specified, "untitled.sysl" is used.
     * @param params The parameters used to parse Sysl. If {@link ParseParams.syslRoot} is not specified, it defaults to
     * the current working directory.
     * @returns A {@link Model} representing the supplied Sysl text.
     */
    static async fromText(
        syslText: string,
        syslFilePath: string = "untitled.sysl",
        params: ParseParams = {}
    ): Promise<Model> {
        // TODO: Improve performance by only reading the first part of the file
        const lines = syslText.split(/\r?\n/);
        const until = lines.findIndex(l => !l.startsWith("#"));
        const header = lines.slice(0, until).join("\n");

        const pb = await PbDocumentModel.fromText(syslText, syslFilePath, params);

        let newModel = pb.toModel();
        newModel.syslRoot = params.syslRoot ?? ".";
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
     * Attempts to finds the element specified by the supplied {@link ElementRef}. Returns `undefined` if no such
     * element was found. Currently only supports apps, types and fields.
     * @param ref The reference used to locate the element in the model. Either a string or an instance of `ElementRef`.
     * @returns An {@link Element} that corresponds to the supplied {@link ElementRef}, or `undefined` if not found.
     * @throws Error thrown if the syntax of the supplied string is invalid.
     */
    findElement(ref: ElementID): Element | undefined {
        const parsed = typeof ref == "string" ? ElementRef.parse(ref) : ref;

        const app = this.apps.find(a => parsed.appsEqual(a.toRef()));
        if (!app || parsed.isApp) return app;

        const type = app.types.find(t => parsed.typesEqual(t.toRef()));
        if (!type || parsed.isType) return type;

        return type.children.find(f => parsed.equals(f.toRef()));
    }

    /**
     * Attempts to finds the app specified by the supplied {@link ElementRef}. Returns `undefined` if no such element
     * was found.
     * @param ref The reference used to locate the app in the model. Either a string or an instance of `ElementRef`.
     * @returns An {@link Application} that corresponds to the supplied {@link ElementRef}, or `undefined` if not found.
     * @throws Error thrown if the syntax of the supplied string is invalid.
     * @throws Error thrown if `ref` is not reference to an app.
     */
    findApp(ref: ElementID): Application | undefined {
        const parsed = typeof ref == "string" ? ElementRef.parse(ref) : ref;
        if (!parsed.isApp) throw new Error(`Supplied reference '${ref.toString()} is not an app.`);
        return this.findElement(parsed) as Application;
    }

    /**
     * Attempts to finds the type specified by the supplied {@link ElementRef}. Returns `undefined` if no such type was
     * found.
     * @param ref The reference used to locate the type in the model. Either a string or an instance of `ElementRef`.
     * @returns A {@link Type} that corresponds to the supplied {@link ElementRef}, or `undefined` if not found.
     * @throws Error thrown if the syntax of the supplied string is invalid.
     * @throws Error thrown if `ref` is not reference to a type.
     */
    findType(ref: ElementID): Type | undefined {
        const parsed = typeof ref == "string" ? ElementRef.parse(ref) : ref;
        if (!parsed.isType) throw new Error(`Supplied reference '${ref.toString()} is not a type.`);
        return this.findElement(parsed) as Type;
    }

    /**
     * Attempts to finds the field specified by the supplied {@link ElementRef}. Returns `undefined` if no such field
     * was found.
     * @param ref The reference used to locate the field in the model. Either a string or an instance of `ElementRef`.
     * @returns A {@link Field} that corresponds to the supplied {@link ElementRef}, or `undefined` if not found.
     * @throws Error thrown if the syntax of the supplied string is invalid.
     * @throws Error thrown if `ref` is not reference to a field.
     */
    findField(ref: ElementID): Field | undefined {
        const parsed = typeof ref == "string" ? ElementRef.parse(ref) : ref;
        if (!parsed.isField) throw new Error(`Supplied reference '${ref.toString()} is not a field.`);
        return this.findElement(parsed) as Field;
    }

    /**
     * Retrieves the element specified by the supplied {@link ElementRef}. Throws if the element was not found.
     * @param ref The reference used to locate the element in the model. Either a string or an instance of `ElementRef`.
     * @returns An {@link Element} that corresponds to the supplied {@link ElementRef}.
     * @throws Error thrown if the syntax of the supplied string is invalid.
     * @throws Error thrown if `ref` describes an element not present in the current {@link Model}.
     */
    getElement = (ref: ElementID) => this.getTypedElement<Element>(ref, this.findElement);

    /**
     * Retrieves the app specified by the supplied {@link ElementRef}. Throws if the app was not found.
     * @param ref The reference used to locate the app in the model. Either a string or an instance of `ElementRef`.
     * @returns An {@link Application} that corresponds to the supplied {@link ElementRef}.
     * @throws Error thrown if `ref` is a string with an invalid element reference syntax.
     * @throws Error thrown if `ref` is an {@link ElementRef} with `kind` other than {@link ElementKind.App}.
     * @throws Error thrown if `ref` describes an app not present in the current {@link Model}.
     */
    getApp = (ref: ElementID) => this.getTypedElement<Application>(ref, this.findApp);

    /**
     * Retrieves the type specified by the supplied {@link ElementRef}. Throws if the type was not found.
     * @param ref The reference used to locate the type in the model. Either a string or an instance of `ElementRef`.
     * @returns A {@link Type} that corresponds to the supplied {@link ElementRef}.
     * @throws Error thrown if `ref` is a string with an invalid element reference syntax.
     * @throws Error thrown if `ref` is an {@link ElementRef} with `kind` other than {@link ElementKind.Type}.
     * @throws Error thrown if `ref` describes a type not present in the current {@link Model}.
     */
    getType = (ref: ElementID) => this.getTypedElement<Type>(ref, this.findType);

    /**
     * Retrieves the field specified by the supplied {@link ElementRef}. Throws if the field was not found.
     * @param ref The reference used to locate the field in the model. Either a string or an instance of `ElementRef`.
     * @returns A {@link Field} that corresponds to the supplied {@link ElementRef}.
     * @throws Error thrown if `ref` is a string with an invalid element reference syntax.
     * @throws Error thrown if `ref` is an {@link ElementRef} with `kind` other than {@link ElementKind.Field}.
     * @throws Error thrown if `ref` describes a type not present in the current {@link Model}.
     */
    getField = (ref: ElementID) => this.getTypedElement<Field>(ref, this.findField);

    private getTypedElement<T extends Element>(ref: ElementID, findFunc: (r: ElementID) => T | undefined): T {
        const element = findFunc.bind(this)(ref);
        if (element == undefined) throw new Error(`Requested element '${ref.toString()}' could not be found in model.`);
        return element;
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
     * @throws Error thrown when the specified path is outside the {@link syslRoot}, which is not allowed.
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

    public clone(filter?: ModelFilter, keepLocation?: boolean): Model {
        const model = new Model({ header: this.header, syslRoot: this.syslRoot });
        const context = new CloneContext(model, filter ?? ModelFilters.Default, 0, keepLocation);
        if (context.keepLocation) model.locations = context.recurse(this.locations);
        model.imports = context.recurse(this.imports);
        model.apps = context.recurse(this.apps);
        return model;
    }

    /**
     * Ensures the `.parent` and `.model` properties of descendants are correctly set.
     */
    public attachSubitems() {
        this.apps.forEach(a => {
            a.model = this;
            a.attachSubitems();
        });
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
