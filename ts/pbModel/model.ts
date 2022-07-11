import { readFile } from "fs/promises";
import { initializeLinq } from "linq-to-typescript";
import { exec } from "promisify-child-process";
import "reflect-metadata";
import {
    jsonArrayMember,
    jsonMapMember,
    jsonMember,
    jsonObject,
    TypedJSON,
} from "typedjson";
import { Location } from "../location";
import { Application, Import, Model } from "../model/model";
import { getAnnos, getTags, sortLocationalArray } from "../sort";
import { PbAttribute } from "./attribute";
import { PbAppName } from "./appname";
import { PbEndpoint } from "./statement";
import { PbTypeDef } from "./type";
import { serializerFor } from "./serialize";
import { joinedAppName } from "../format";

initializeLinq();

@jsonObject
export class PbImport {
    @jsonMember target!: string;
    @jsonMember sourceContext!: Location;
    @jsonMember name?: PbAppName;

    toModel(): Import {
        return new Import(
            this.target,
            [this.sourceContext],
            joinedAppName(this.name?.part ?? [])
        );
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
            joinedAppName(this.name.part),
            sortLocationalArray(
                this.endpoints
                    ?.select(e => e[1].toModel(this.name.part))
                    .toArray() ?? []
            ),
            sortLocationalArray(
                this.types?.select(t => t[1].toModel(t[0], false)).toArray() ??
                    []
            ),
            this.sourceContexts,
            getTags(this.attrs),
            getAnnos(this.attrs)
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

    static async fromText(
        syslText: string,
        syslFilePath: string = "untitled.sysl"
    ): Promise<PbDocumentModel> {
        const syslPath = process.env["SYSL_PATH"] ?? "sysl";
        const cmd = `${syslPath} pb --mode=json`;
        const proc = exec(cmd);
        console.debug(
            `Sysl Path: ${syslPath}\n${(await exec(`${syslPath} info`)).stdout}`
        );
        proc.stdin?.end(
            JSON.stringify([{ path: syslFilePath, content: syslText }])
        );
        const json = (await proc).stdout;
        const serializer = new TypedJSON(PbDocumentModel, {
            errorHandler: error => {
                throw error;
            },
        });
        const model = serializer.parse(json)!;
        return model;
    }

    toModel(): Model {
        return new Model(
            this.imports?.select(i => i.toModel()).toArray() ?? [],
            this.sourceContexts,
            sortLocationalArray(
                this.apps.select(a => a[1].toModel()).toArray()
            ),
            ""
        );
    }
}
