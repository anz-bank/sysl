import "jest-extended";
import "reflect-metadata";
import * as path from "path";
import { PbDocumentModel } from "./model";
import { readFile } from "fs/promises";

const testDir = path.join(__dirname, "..", "..", "test");
const allPath = path.join(testDir, "all.sysl");
const jsonPath = path.join(testDir, "all.json");
const pbPath = path.join(testDir, "all.pb");

describe("PbDocumentModel", () => {
    test("from text", async () => {
        const source = (await readFile(allPath)).toString("utf8");
        const m = await PbDocumentModel.fromText(source, allPath);
        expect(m.apps.size).toEqual(7);
    });

    test("from JSON string", async () => {
        const source = (await readFile(jsonPath)).toString("utf8");
        const m = PbDocumentModel.fromJson(source);
        expect(m.apps.size).toEqual(7);
    });

    test("from JSON buffer", async () => {
        const source = await readFile(jsonPath);
        const m = PbDocumentModel.fromJson(source);
        expect(m.apps.size).toEqual(7);
    });

    test("from pb", async () => {
        const source = await readFile(pbPath);
        const m = await PbDocumentModel.fromPbOrJson(source);
        expect(m.apps.size).toEqual(7);
    });
});
