import { readFile } from "fs/promises";
import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject, TypedJSON } from "typedjson";
import { joinedAppName } from "../common/format";
import { Location } from "../common/location";
import { sortByLocation } from "../common/sort";
import { ElementRef, Import, Model } from "../model";
import { Application } from "../model/application";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor } from "./serialize";
import { PbEndpoint } from "./statement";
import { PbTypeDef } from "./type";
import { spawnBuffer } from "../common/spawn";

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
        const name = this.name.part.at(-1);
        if (!name) throw new Error("Encountered empty app name.");
        return new Application(new ElementRef(this.name.part.slice(0, -1), name), {
            endpoints: sortByLocation(
                Array.from(this.endpoints ?? new Map<string, PbEndpoint>())
                    .filter(([, e]) => e.name != "...") // Bug in Sysl where ellipsis under app appears as endpoint
                    .map(([, e]) => e.toModel(this.name.part))
            ),
            types: sortByLocation(
                [...(this.types ?? new Map<string, PbTypeDef>())].map(([name, t]) => t.toModel(name, false))
            ),
            locations: this.sourceContexts,
            tags: sortByLocation(getTags(this.attrs)),
            annos: sortByLocation(getAnnos(this.attrs)),
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
    static async fromText(syslText: string, syslFilePath: string, maxImportDepth?: number): Promise<PbDocumentModel> {
        const files = [{ path: syslFilePath, content: syslText }];
        return this.fromPbOrJson(JSON.stringify(files), maxImportDepth);
    }

    /** Deserializes a compiled JSON string into a model. */
    static fromJson(json: string | Buffer): PbDocumentModel {
        const serializer = new TypedJSON(PbDocumentModel, {
            errorHandler: error => {
                throw error;
            },
        });

        try {
            return serializer.parse(json.toString("utf-8"))!;
        } catch (error: any) {
            throw new Error("Parsing of the following document failed. See cause for details.\n" + json.toString(), {
                cause: error,
            });
        }
    }

    /**
     * Returns a model compiled by invoking `sysl pb` and passing it the {@link content} via
     * the standard input.
     *
     * @param content is the content to construct a model from. This can be either the bytes of a
     *        precompiled `.pb` file, or a JSON object of the form
     *        `[{"path": "path/to/file.sysl", "content": "sysl source"}, ...]`.
     *        In the first (simple) case, the precompiled model loaded directly into a TypeScript
     *        object. In the second case, the Sysl source in the {@link content} property is
     *        compiled and then loaded.
     * @param maxImportDepth sets the max import depth for the compiler. This has no effect if
     *        {@link content} is a precompiled `.pb` file.
     */
    static async fromPbOrJson(content: string | Buffer, maxImportDepth?: number): Promise<PbDocumentModel> {
        const syslPath = process.env["SYSL_PATH"] ?? "sysl";
        const out = await spawnBuffer(
            syslPath,
            ["pb", "--mode=json", "--compact", `--max-import-depth=${maxImportDepth ?? 0}`],
            {
                input: content,
            }
        );
        return PbDocumentModel.fromJson(out);
    }

    toModel(): Model {
        return new Model({
            imports: (this.imports ?? []).map(i => i.toModel()),
            locations: this.sourceContexts,
            apps: sortByLocation(Array.from(this.apps).map(([, a]) => a.toModel())),
        });
    }
}
