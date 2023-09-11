import { readFile } from "fs/promises";
import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject, TypedJSON } from "typedjson";
import { joinedAppName } from "../common/format";
import { Location } from "../common/location";
import { sortByLocation } from "../common/sort";
import { Element, ElementRef, Model, ParseParams } from "../model";
import { Application } from "../model/application";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor } from "./serialize";
import { PbEndpoint } from "./statement";
import { PbTypeDef } from "./type";
import { spawnBuffer } from "../common/spawn";
import { Import } from "../model/import";

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
        const appRef = new ElementRef(this.name.part.slice(0, -1), name);
        const types = Array.from(this.endpoints ?? new Map<string, PbEndpoint>())
            .filter(([, e]) => e.name != "..." && !e.isPubsub) // Bug where ellipsis under app appears as endpoint
            .map(([, e]) => e.toModel(this.name.part));
        const endpoints = Array.from(this.types ?? new Map<string, PbTypeDef>(), ([name, t]) =>
            t.toModel(name, false, appRef)
        );
        return new Application(appRef, {
            children: sortByLocation([...types, ...endpoints]),
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

    /**
     * Creates a {@link PbDocumentModel} from a Sysl file on disk.
     *
     * @param syslFilePath The path to the Sysl file.
     * @param params The parameters used to parse Sysl. {@link ParseParams.syslRoot} is ignored.
     * @returns A promise of a {@link PbDocumentModel} that corresponds to the provided Sysl file.
     */
    static async fromFile(syslFilePath: string, params: ParseParams = {}): Promise<PbDocumentModel> {
        const syslText = (await readFile(syslFilePath)).toString();
        return this.fromText(syslText, syslFilePath, params);
    }

    /**
     * Creates a {@link PbDocumentModel} from an in-memory Sysl file.
     *
     * @param syslText The contents of the Sysl file.
     * @param syslFilePath The reported path of the Sysl file. This path doesn't have to correspond to an actual file on
     * the disk, but will be assumed to be one for when reporting locations of sysl elements.
     * @param params The parameters used to parse Sysl. {@link ParseParams.syslRoot} is ignored.
     * @returns A promise of a {@link PbDocumentModel} that corresponds to the provided Sysl text.
     */
    static async fromText(syslText: string, syslFilePath: string, params: ParseParams = {}): Promise<PbDocumentModel> {
        const files = [{ path: syslFilePath, content: syslText }];
        return this.fromPbOrJson(JSON.stringify(files), params);
    }

    /**
     * Creates a {@link PbDocumentModel} from an in-memory PB file in JSON format.
     *
     * @param json The contents of the JSON-formatted PB file.
     * @returns A {@link PbDocumentModel} that corresponds to the provided JSON.
     */
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
     * Creates a {@link PbDocumentModel} from in-memory data, either a pb file or a collection of Sysl files.
     *
     * @param content A string or buffer from which to construct a {@link PbDocumentModel}. This can be either the bytes
     * of a precompiled `.pb` file, or serialised JSON that contains an in-memory file system in the form of an array of
     * `{ path, content }` objects. For example:
     * `[{"path": "model/backend/database.sysl", "content": "MyShop :: Backend :: Database\n  ..."}]`.
     * @param params The parameters used to parse Sysl. {@link ParseParams.syslRoot} is ignored.
     * @returns A promise of a {@link PbDocumentModel} that corresponds to the provided content. This setting has no
     * effect if {@link content} is a precompiled `.pb` file.
     */
    static async fromPbOrJson(
        content: string | Buffer,
        { maxImportDepth, alwaysFetch }: ParseParams = {}
    ): Promise<PbDocumentModel> {
        const syslPath = process.env["SYSL_PATH"] ?? "sysl";
        const args = ["pb", "--mode=json", "--compact"];
        if ((maxImportDepth ?? 0) > 0) args.push(`--max-import-depth=${maxImportDepth}`);
        if (alwaysFetch == false) args.push("--no-forced-fetch");
        const out = await spawnBuffer(syslPath, args, { input: content });
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
