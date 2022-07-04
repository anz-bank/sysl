import { readFile } from "fs/promises";
import { initializeLinq } from "linq-to-typescript";
import "reflect-metadata";
import { Location } from "../location";
import { PbDocumentModel } from "../pbModel/model";
import { indent } from "../util";
import { ComplexType } from "./common";
import { Element } from "./element";
import { Endpoint } from "./statement";
import { Type } from "./type";

initializeLinq();

export class Import {
    filePath: string;
    locations: Location[]
    appAlias?: string

    constructor(filePath: string, locations: Location[], appAlias: string) {
        this.filePath = filePath
        this.locations = locations;
        this.appAlias = appAlias ?? undefined;
    }

    toSysl(): string {
        return `import ${this.filePath}${this.appAlias?.any() ? ` as ${this.appAlias}` : ''}`;
    }
}

export class Application extends ComplexType {
    endpoints: Element<Endpoint>[];
    types: Element<Type>[];

    constructor(endpoints: Element<Endpoint>[],
        types: Element<Type>[],
        name: string,
        locations: Location[],
    ) {
        super("", locations, name);
        this.endpoints = endpoints;
        this.types = types;
        this.name = name;
        this.locations = locations;
    }

    override toSysl(): string {
        let sysl = ``;
        if (this.endpoints.any()) {
            sysl += `${this.endpoints.where(e => !e.content.isPubsub).select(e => indent(e.content.toSysl())).toArray().join("\n\n")}`;
        }
        if (this.types.any()) {
            sysl += `${this.types.select(t => indent(t.toSysl())).toArray().join("\n\n")}`;
        }
        return sysl;
    }
}

export class Model {
    imports: Import[];
    locations: Location[];
    apps: Element<Application>[];
    header: string;

    constructor(imports: Import[], locations: Location[], apps: Element<Application>[], header: string) {
        this.imports = imports;
        this.locations = locations;
        this.apps = apps;
        this.header = header;
    }

    static async fromFile(syslFilePath: string): Promise<Model> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath);
    }

    static async fromText(syslText: string, syslFilePath: string = "untitled.sysl"): Promise<Model> {
        // TODO: Improve performance by only reading the first part of the file
        const header = syslText.split(/\r?\n/).takeWhile(l => l.startsWith("#")).toArray().join("\n");

        const pb = await PbDocumentModel.fromText(syslText, syslFilePath)

        let newModel = pb.toModel();
        if (header) {
            newModel.header = header;
        }

        return newModel;
    }

    filterByFile(file: string): Model {
        return new Model(
            this.imports ? this.imports.where(i => file.includes(i.locations[0]?.file!)).toArray() : [],
            this.locations,
            this.apps?.where(a => file.includes(a.content.locations[0]?.file!)).toArray(),
            this.header
        );
    }

    toSysl(): string {
        let sysl = "";

        if (this.header) {
            sysl = this.header + "\n\n";
        }

        if (this.imports?.any()) {
            sysl += this.imports.select(i => i.toSysl())
                .toArray()
                .join("\n") + "\n\n";
        }

        sysl += this.apps
            .select(a => `${a.toSysl()}`)
            .toArray()
            .join("\n\n");

        sysl += "\n";

        return sysl;
    }
}
