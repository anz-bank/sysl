import { readFile } from "fs/promises";
import "reflect-metadata";
import { allItems } from "../common/iterate";
import { Location } from "../common/location";
import { PbDocumentModel } from "../pbModel/model";
import { Application } from "./application";
import {
    IRenderable,
} from "./common";
import { setModelDeep } from "./element";

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

    static async fromFile(syslFilePath: string, maxImportDepth: number = 0): Promise<Model> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath, maxImportDepth);
    }

    static async fromText(
        syslText: string,
        syslFilePath: string = "untitled.sysl",
        maxImportDepth: number = 0
    ): Promise<Model> {
        // TODO: Improve performance by only reading the first part of the file
        const lines = syslText.split(/\r?\n/);
        const until = lines.findIndex(l => !l.startsWith("#"));
        const header = lines.slice(0, until).join("\n");

        const pb = await PbDocumentModel.fromText(syslText, syslFilePath, maxImportDepth);

        let newModel = pb.toModel();
        if (header) {
            newModel.header = header;
        }

        return newModel;
    }

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
