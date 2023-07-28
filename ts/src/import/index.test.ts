import { writeFile } from "fs/promises";
import "jest-extended";
import path from "path";
import { flattenStatement, mergeExisting } from ".";
import { spawnSysl } from "../common/spawn";
import { Model } from "../model";

describe("import", () => {
    describe("merge", () => {
        test("from files", async () => {
            const source = testFile("test.proto");
            const existing = testFile("test.proto.sysl");
            const oldMod = Model.fromFile(existing);
            const newMod = syslImport(source).then(s => Model.fromText(s));

            mergeExisting(await newMod, await oldMod);

            const ep = (await newMod).findApp("Test::SearchService")!.endpoints[0];
            const stmts = ep.statements.flatMap(flattenStatement);
            expect(stmts).toHaveLength(5);
            expect(stmts[3]).toMatchObject({ value: { payload: "ok <: Test :: Types.SearchResponse" } });
            expect((await newMod).header).toInclude("CAUTION");

            await writeFile(existing, (await newMod).toSysl());
        }, 10000);
    });
});

function testFile(name: string): string {
    return path.relative(process.cwd(), path.join(__dirname, "tests", name));
}

async function syslImport(input: string): Promise<string> {
    return (await spawnSysl(["import", "-i", input, "-f", "protobuf", "--root", __dirname])).toString();
}
