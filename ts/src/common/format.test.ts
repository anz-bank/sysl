import { realign, safeName, unescapeName } from "./format";

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
            safeName("abcdefgihjklmnopqrstuvwxyz-ABCDEFGIHJKLMNOPQRSTUVWXYZ_")
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
            expect(safeName(keyword)).toEqual(keyword + "_");
        });
    });

    test("special characaters underscore", () => {
        expect(safeName(`a/b\\c{d}e f`)).toEqual("a_b_c_d_e_f");
    });

    test("special characaters hex escaped", () => {
        expect(safeName(`hello;,?:@&=+$.!~*'()"world`)).toEqual(
            `hello%3b%2c%3f%3a%40%26%3d%2b%24%2e%21%7e%2a%27%28%29%22world`
        );
    });
});

describe("unescapeName", () => {
    test("no escaping", () => {});

    test("escaped", () => {
        expect(
            unescapeName(
                `hello%3b%2c%3f%3a%40%26%3d%2b%24%2e%21%7e%2a%27%28%29%22world`
            )
        ).toEqual(`hello;,?:@&=+$.!~*'()"world`);
    });
});
