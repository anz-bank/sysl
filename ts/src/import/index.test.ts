import "jest-extended";
import path from "path";
import { flattenStatement, importAndMerge } from ".";
import { ElementRef, Statement } from "../model";

describe("import and merge", () => {
    test("via CLI", async () => {
        const source = testFile("test.proto");
        const existing = testFile("test.proto.sysl");
        const merged = await importAndMerge({
            input: source,
            output: existing,
            format: "protobufDir",
            shallow: true,
        });
        const app = merged.model.getApp(new ElementRef(["Test"], "SearchService"));
        expect(app.endpoints[0].children.flatMap(flattenStatement).map((s: Statement) => s.toString())).toEqual([
            "hello",
            "if world",
            "world",
            "return ok <: Test :: Types.SearchResponse",
            "Foreign <- Endpoint",
        ]);
    }, 10000); // Slow due to arr.ai-based proto import.
});

function testFile(name: string): string {
    return path.relative(process.cwd(), path.join(__dirname, "tests", name));
}
