import { readFile } from "fs/promises";
import { initializeLinq } from "linq-to-typescript";
import { exec } from "promisify-child-process";
import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject, TypedJSON } from "typedjson";
import { Location } from "../location";
import { Element } from "../model/element";
import { Application, Import, Model } from "../model/model";
import { Endpoint } from "../model/statement";
import { Type } from "../model/type";
import { getAnnos, getTags, joinedAppName, serializerFor, sortElements } from "../util";
import { PbAttribute } from "./attribute";
import { PbAppName } from "./appname";
import { PbEndpoint } from "./statement";
import { PbTypeDef } from "./type";

initializeLinq();

@jsonObject
export class PbImport {
    @jsonMember target!: string;
    @jsonMember sourceContext!: Location;
    @jsonMember name?: PbAppName;

    toModel(): Import {
        return new Import(this.target, [this.sourceContext], joinedAppName(this.name?.part ?? []));
    }
}

@jsonObject
export class PbApplication {
    @jsonMember longName?: string;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMember name!: PbAppName;
    @jsonMapMember(String, PbAttribute, serializerFor(PbAttribute))
    attrs?: Map<string, PbAttribute>;
    @jsonMapMember(String, PbEndpoint, serializerFor(PbEndpoint))
    endpoints?: Map<string, PbEndpoint>;
    @jsonMapMember(String, PbTypeDef, serializerFor(PbTypeDef))
    types?: Map<string, PbTypeDef>;

    toModel(): Application {
        return new Application(
            sortElements(this.endpoints?.select(e => new Element<Endpoint>(getTags(e[1].attrs), getAnnos(e[1].attrs), e[1].toModel(this.name.part))).toArray() ?? []),
            sortElements(this.types?.select(t => new Element<Type>(getTags(t[1].attrs), getAnnos(t[1].attrs), t[1].toModel(t[0]))).toArray() ?? []),
            joinedAppName(this.name.part),
            this.sourceContexts,
        );
    }
}

@jsonObject
export class PbDocumentModel {
    @jsonArrayMember(PbImport) imports?: PbImport[];
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMapMember(String, PbApplication, serializerFor(PbApplication))
    apps!: Map<string, PbApplication>;

    static async fromFile(syslFilePath: string): Promise<PbDocumentModel> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath);
    }

    static async fromText(syslText: string, syslFilePath: string = "untitled.sysl"): Promise<PbDocumentModel> {
        const syslPath = process.env["SYSL_PATH"] ?? "sysl";
        const cmd = `${syslPath} pb --mode=json`;
        const proc = exec(cmd);
        try {
            console.debug(`Sysl Path: ${syslPath}\n${(await exec(`${syslPath} info`)).stdout}`)
            proc.stdin?.end(JSON.stringify([{ path: syslFilePath, content: syslText }]));
            const json = (await proc).stdout;
            const serializer = new TypedJSON(PbDocumentModel, { errorHandler: (error) => { throw error; } });
            const model = serializer.parse(json)!;
            return model;
        }
        catch {
            throw new Error("Sysl binary not found. Please visit https://sysl.io/docs/installation for installation instructions.")
        }
    }

    toModel(): Model {
        return new Model(
            this.imports?.select(i => i.toModel()).toArray() ?? [],
            this.sourceContexts,
            sortElements(this.apps.select(a => new Element<Application>(getTags(a[1].attrs), getAnnos(a[1].attrs), a[1].toModel())).toArray()),
            "");
    }
}
