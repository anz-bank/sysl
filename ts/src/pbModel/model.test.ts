import "jest-extended";
import "reflect-metadata";
import * as path from "path";
import { PbDocumentModel } from "./model";
import { readFile } from "fs/promises";
import { realign } from "../common/format";

const testDir = path.join(__dirname, "..", "..", "test");
const allPath = path.join(testDir, "all.sysl");
const jsonPath = path.join(testDir, "all.json");
const pbPath = path.join(testDir, "all.pb");

describe("PbDocumentModel", () => {
    test("from text", async () => {
        const source = (await readFile(allPath)).toString("utf8");
        const m = await PbDocumentModel.fromText(source, allPath);
        expect(m.apps.size).toEqual(9);
    });

    test("from JSON string", async () => {
        const source = (await readFile(jsonPath)).toString("utf8");
        const m = PbDocumentModel.fromJson(source);
        expect(m.apps.size).toEqual(9);
    });

    test("from JSON buffer", async () => {
        const source = await readFile(jsonPath);
        const m = PbDocumentModel.fromJson(source);
        expect(m.apps.size).toEqual(9);
    });

    test("from pb", async () => {
        const source = await readFile(pbPath);
        const m = await PbDocumentModel.fromPbOrJson(source);
        expect(m.apps.size).toEqual(9);
    });
});

describe("PbAnnotation", () => {
    test("backslash", async () => {
        // prettier-ignore
        const m = await PbDocumentModel.fromText(
            realign(`
            App [attr=["foo = C:\\\\bar", "\\t", "\\"quote\\""]]:
                @anno = ["foo = C:\\\\bar", "\\t", "\\"quote\\""]
        `),
            "test.sysl"
        );
        const attrs = m.apps.get("App")!.attrs!;
        const attrValues = attrs.get("attr")!.a!.elt;
        const annoValues = attrs.get("anno")!.a!.elt;

        expect(annoValues[0].s).toEqual("foo = C:\\\\bar");
        expect(annoValues[1].s).toEqual("\t");
        expect(annoValues[2].s).toEqual(`"quote"`);

        expect(attrValues[0].s).toEqual("foo = C:\\\\bar");
        expect(attrValues[1].s).toEqual("\t");
        expect(attrValues[2].s).toEqual(`"quote"`);
    });
});
