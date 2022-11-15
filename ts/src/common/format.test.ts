import { realign, toSafeName, fromSafeName } from "./format";

test("realign", () => {
    expect(
        realign(`
    App "display name" [~abstract]:
        ...
    `)
    ).toBe(
        `App "display name" [~abstract]:
    ...
`
    );
});

describe("safeName", () => {
    test("safe", () => {
        expect(
            toSafeName("abcdefgihjklmnopqrstuvwxyz-ABCDEFGIHJKLMNOPQRSTUVWXYZ_")
        ).toEqual("abcdefgihjklmnopqrstuvwxyz-ABCDEFGIHJKLMNOPQRSTUVWXYZ_");
    });

    test("keywords", () => {
        [
            "int",
            "int32",
            "int64",
            "float",
            "float32",
            "float64",
            "decimal",
            "string",
            "date",
            "datetime",
            "bool",
            "bytes",
            "any",
        ].forEach(keyword => {
            expect(toSafeName(keyword)).toEqual(keyword + "_");
        });
    });

    test("special characters hex escaped", () => {
        expect(toSafeName(`hello;,?:@&=+$.!~*'()"/\\world`)).toEqual(
            `hello%3B%2C%3F%3A%40%26%3D%2B%24%2E%21%7E%2A%27%28%29%22%2F%5Cworld`
        );
    });
});

describe("unescapeName", () => {
    test("no escaping", () => {});

    test("escaped lowercase", () => {
        expect(
            fromSafeName(
                `hello%3b%2c%3f%3a%40%26%3d%2b%24%2e%21%7e%2a%27%28%29%22%2f%5cworld`
            )
        ).toEqual(`hello;,?:@&=+$.!~*'()"/\\world`);
    });
    test("escaped uppercase", () => {
        expect(
            fromSafeName(
                `hello%3B%2C%3F%3A%40%26%3D%2B%24%2E%21%7E%2A%27%28%29%22%2F%5Cworld`
            )
        ).toEqual(`hello;,?:@&=+$.!~*'()"/\\world`);
    });
});
