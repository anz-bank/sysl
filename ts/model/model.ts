import { readFile } from "fs/promises";
import { initializeLinq } from "linq-to-typescript";
import "reflect-metadata";
import { Location } from "../location";
import { PbDocumentModel } from "../pbModel/model";
import { indent } from "../format";
import { IDescribable, ILocational, IRenderable } from "./common";
import { Endpoint } from "./statement";
import { Type } from "./type";
import { Annotation, Tag } from "./attribute";
import { renderAnnos, addTags } from "./renderers";

initializeLinq();

export class Import {
    filePath: string;
    locations: Location[];
    appAlias?: string;

    constructor(filePath: string, locations: Location[], appAlias: string) {
        this.filePath = filePath;
        this.locations = locations;
        this.appAlias = appAlias ?? undefined;
    }

    toSysl(): string {
        return `import ${this.filePath}${
            this.appAlias?.any() ? ` as ${this.appAlias}` : ""
        }`;
    }
}

export class Application implements IDescribable, ILocational, IRenderable {
    name: string;
    endpoints: Endpoint[];
    types: Type[];
    locations: Location[];
    tags: Tag[];
    annos: Annotation[];

    constructor(
        name: string,
        endpoints: Endpoint[],
        types: Type[],
        locations: Location[],
        tags: Tag[],
        annos: Annotation[]
    ) {
        this.name = name;
        this.endpoints = endpoints;
        this.types = types;
        this.name = name;
        this.locations = locations;
        this.tags = tags;
        this.annos = annos;
    }

    toSysl(): string {
        let sysl = `${addTags(this.name, this.tags)}:`;
        if (this.annos?.any()) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        if (this.endpoints.any()) {
            sysl += `\n${this.endpoints
                .where(e => !e.isPubsub)
                .select(e => indent(e.toSysl()))
                .toArray()
                .join("\n\n")}`;
        }
        if (this.types.any()) {
            sysl += `\n${this.types
                .select(t => indent(t.toSysl()))
                .toArray()
                .join("\n\n")}`;
        }
        return sysl;
    }
}

export class Model implements IRenderable {
    imports: Import[];
    locations: Location[];
    apps: Application[];
    header: string;

    constructor(
        imports: Import[],
        locations: Location[],
        apps: Application[],
        header: string
    ) {
        this.imports = imports;
        this.locations = locations;
        this.apps = apps;
        this.header = header;
    }

    static async fromFile(syslFilePath: string): Promise<Model> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath);
    }

    static async fromText(
        syslText: string,
        syslFilePath: string = "untitled.sysl"
    ): Promise<Model> {
        // TODO: Improve performance by only reading the first part of the file
        const header = syslText
            .split(/\r?\n/)
            .takeWhile(l => l.startsWith("#"))
            .toArray()
            .join("\n");

        const pb = await PbDocumentModel.fromText(syslText, syslFilePath);

        let newModel = pb.toModel();
        if (header) {
            newModel.header = header;
        }

        return newModel;
    }

    filterByFile(file: string): Model {
        return new Model(
            this.imports
                ? this.imports
                      .where(i => file.includes(i.locations[0]?.file!))
                      .toArray()
                : [],
            this.locations,
            this.apps
                ?.where(a => file.includes(a.locations[0]?.file!))
                .toArray(),
            this.header
        );
    }

    toSysl(): string {
        let sysl = "";

        if (this.header) {
            sysl = this.header + "\n\n";
        }

        if (this.imports?.any()) {
            sysl +=
                this.imports
                    .select(i => i.toSysl())
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
