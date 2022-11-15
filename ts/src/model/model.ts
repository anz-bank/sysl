import { readFile } from "fs/promises";
import "reflect-metadata";
import { allItems } from "../common/iterate";
import { Location } from "../common/location";
import { PbDocumentModel } from "../pbModel/model";
import { Application } from "./application";
import { ElementKind, ElementRef, IRenderable } from "./common";
import { Element, setModelDeep } from "./element";

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
        return `import ${this.filePath}${
            this.appAlias ? ` as ${this.appAlias}` : ""
        }`;
    }
}

export type ModelParams = {
    header?: string | undefined;
    imports?: Import[];
    apps?: Application[];
    locations?: Location[];
};

export class Model implements IRenderable {
    header: string | undefined;
    imports: Import[];
    apps: Application[];
    locations: Location[];

    constructor({ header, imports, apps, locations }: ModelParams = {}) {
        this.header = header;
        this.imports = imports ?? [];
        this.apps = apps ?? [];
        this.locations = locations ?? [];

        setModelDeep(this, this.apps);
        allItems(this).forEach(i => (i.model = this));
    }

    /**
     * Parses Sysl text file into a {@link Model}.
     * @param syslFilePath The file path that contains the Sysl text to parse.
     * @param maxImportDepth Optional. The maximum depth to follow import statements, where 0 means unlimited depth,
     * 1 means to remain at the level of the supplied file (do not to follow imports at all), 2 means to delve one
     * deeper (supplied file plus one extra depth), etc. If not specified, defaults to allowing unlimited depth.
     * Limiting the depth may significantly improve performance, especially if it causes remote imports to be skipped.
     * @returns A {@link Model} representing the supplied Sysl file.
     */
    static async fromFile(
        syslFilePath: string,
        maxImportDepth: number = 0
    ): Promise<Model> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath, maxImportDepth);
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
        maxImportDepth: number = 0
    ): Promise<Model> {
        // TODO: Improve performance by only reading the first part of the file
        const lines = syslText.split(/\r?\n/);
        const until = lines.findIndex(l => !l.startsWith("#"));
        const header = lines.slice(0, until).join("\n");

        const pb = await PbDocumentModel.fromText(
            syslText,
            syslFilePath,
            maxImportDepth
        );

        let newModel = pb.toModel();
        if (header) {
            newModel.header = header;
        }

        return newModel;
    }

    /**
     * Returns a new model with only the apps that were defined in the specified source file. Non-app elements that
     * came from other source and were merged into apps from the specified source file will not be filtered out.
     * @param file The source file path by which to filter apps.
     * @returns A clone of the current {@link Model} instance which only includes apps that were specified by the
     * supplied source file.
     */
    filterByFile(file: string): Model {
        return new Model({
            imports: this.imports.filter(i =>
                file.includes(i.locations[0]?.file!)
            ),
            locations: this.locations,
            apps: this.apps.filter(a => file.includes(a.locations[0]?.file!)),
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
}
