import "jest-extended";
import path from "path";
import { importAndMerge } from ".";

describe("import and merge", () => {
    test("via CLI", async () => {
        const source = testFile("test.proto");
        const existing = testFile("test.proto.sysl");
        await importAndMerge({
            input: source,
            output: existing,
            format: "protobufDir",
            shallow: true,
        });
        // Slow due to arr.ai-based proto import.
    }, 10000);
});

function testFile(name: string): string {
    return path.relative(process.cwd(), path.join(__dirname, "tests", name));
}
