import { readFile } from "fs/promises";
import { exec } from "promisify-child-process";
import "reflect-metadata";
import {
    jsonArrayMember,
    jsonMapMember,
    jsonMember,
    jsonObject,
    TypedJSON,
} from "typedjson";
import { joinedAppName } from "../common/format";
import { Location } from "../common/location";
import { sortLocationalArray } from "../common/sort";
import { Import, Model } from "../model";
import { Application } from "../model/application";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor } from "./serialize";
import { PbEndpoint } from "./statement";
import { PbTypeDef } from "./type";

@jsonObject
export class PbImport {
    @jsonMember target!: string;
    @jsonMember sourceContext!: Location;
    @jsonMember name?: PbAppName;

    toModel(): Import {
        return new Import({
            filePath: this.target,
            locations: [this.sourceContext],
            appAlias: joinedAppName(this.name?.part ?? []),
        });
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
    @jsonMapMember(String, () => PbTypeDef, serializerFor(PbTypeDef))
    types?: Map<string, PbTypeDef>;

    toModel(): Application {
        return new Application({
            name: this.name.part.at(-1),
            namespace: this.name.part.slice(0, -1), 
            endpoints: sortLocationalArray(
                Array.from(this.endpoints ?? new Map()).map(([, e]) =>
                    e.toModel(this.name.part)
                )
            ),
            children: sortLocationalArray(
                Array.from(this.types ?? new Map()).map(([name, t]) => {
                    return t.toModel(name, false);
                })
            ),
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        });
    }
}

@jsonObject
export class PbDocumentModel {
    @jsonArrayMember(PbImport) imports?: PbImport[];
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMapMember(String, PbApplication, serializerFor(PbApplication))
    apps!: Map<string, PbApplication>;

    static async fromFile(syslFilePath: string, maxImportDepth: number): Promise<PbDocumentModel> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath, maxImportDepth);
    }

    /** Compiles and deserializes a Sysl source string into a model. */
    static async fromText(
        syslText: string,
        syslFilePath: string,
        maxImportDepth: number
    ): Promise<PbDocumentModel> {
        const syslPath = process.env["SYSL_PATH"] ?? "sysl";
        const cmd = `${syslPath} pb --mode=json --max-import-depth=${maxImportDepth}`;
        const proc = exec(cmd, { maxBuffer: undefined }); // Do not limit the maximum stdout of the process.
        proc.stdin?.end(
            JSON.stringify([{ path: syslFilePath, content: syslText }])
        );
        return PbDocumentModel.fromJson((await proc).stdout!);
    }

    /** Deserializes a compiled JSON string into a model. */
    static fromJson(json: string | Buffer): PbDocumentModel {
        const serializer = new TypedJSON(PbDocumentModel, {
            errorHandler: error => {
                throw error;
            },
        });
        return serializer.parse(json)!;
    }

    toModel(): Model {
        return new Model({
            imports: (this.imports ?? []).map(i => i.toModel()),
            locations: this.sourceContexts,
            apps: sortLocationalArray(
                Array.from(this.apps).map(([, a]) => a.toModel())
            ),
        });
    }
}
