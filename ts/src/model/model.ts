import { readFile } from "fs/promises";
import "reflect-metadata";
import { indent, joinedAppName, safeName } from "../common/format";
import { Location } from "../common/location";
import { PbDocumentModel } from "../pbModel/model";
import { Annotation, Tag } from "./attribute";
import { IElement, IElementParams, IRenderable, setParentDeep } from "./common";
import { addTags, renderAnnos } from "./renderers";
import { Endpoint } from "./statement";
import { Type } from "./type";

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

export class AppName {
    parts: string[];

    constructor(parts: string[]) {
        this.parts = parts;
    }

    static fromString(name: string): AppName {
        return new AppName(name.split(/\s*::\s*/));
    }

    toSysl() {
        return joinedAppName(this.parts.map(safeName));
    }
}

export type ApplicationParams = IElementParams & {
    name: AppName;
    endpoints?: Endpoint[];
    types?: Type[];
};

export class Application implements IElement {
    name: AppName;
    endpoints: Endpoint[];
    types: Type[];
    locations: Location[];
    annos: Annotation[];
    tags: Tag[];
    parent: undefined;
    model?: Model;

    constructor({
        name,
        endpoints,
        types,
        locations,
        annos,
        tags,
        model,
    }: ApplicationParams) {
        this.name = name;
        this.endpoints = endpoints ?? [];
        this.types = types ?? [];
        this.locations = locations ?? [];
        this.annos = annos ?? [];
        this.tags = tags ?? [];
        this.model = model;

        setParentDeep(this, this.endpoints, this.types, this.annos, this.tags);
    }

    toSysl(): string {
        let sysl = `${addTags(this.name.toSysl(), this.tags)}:`;
        if (this.annos.length) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        if (this.endpoints.length) {
            sysl += `\n${this.endpoints
                .filter(e => !e.isPubsub)
                .map(e => indent(e.toSysl()))
                .join("\n\n")}`;
        }
        if (this.types.length) {
            sysl += `\n${this.types.map(t => indent(t.toSysl())).join("\n\n")}`;
        }
        if (
            !this.annos.length &&
            !this.endpoints.length &&
            !this.types.length
        ) {
            sysl += `\n${indent("...")}`;
        }
        return sysl;
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
        const lines = syslText.split(/\r?\n/);
        const until = lines.findIndex(l => !l.startsWith("#"));
        const header = lines.slice(0, until).join("\n");

        const pb = await PbDocumentModel.fromText(syslText, syslFilePath);

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
