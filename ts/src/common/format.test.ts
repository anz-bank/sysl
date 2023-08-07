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
        expect(toSafeName("abcdefgihjklmnopqrstuvwxyz-ABCDEFGIHJKLMNOPQRSTUVWXYZ_")).toEqual(
            "abcdefgihjklmnopqrstuvwxyz-ABCDEFGIHJKLMNOPQRSTUVWXYZ_"
        );
    });

    test("special names", () => {
        const names = [
            { identifier: "int",      safeName: "%69nt"      },
            { identifier: "int32",    safeName: "%69nt32"    },
            { identifier: "int64",    safeName: "%69nt64"    },
            { identifier: "float",    safeName: "%66loat"    },
            { identifier: "float32",  safeName: "%66loat32"  },
            { identifier: "float64",  safeName: "%66loat64"  },
            { identifier: "decimal",  safeName: "%64ecimal"  },
            { identifier: "string",   safeName: "%73tring"   },
            { identifier: "date",     safeName: "%64ate"     },
            { identifier: "datetime", safeName: "%64atetime" },
            { identifier: "bool",     safeName: "%62ool"     },
            { identifier: "bytes",    safeName: "%62ytes"    },
            { identifier: "any",      safeName: "%61ny"      },
            { identifier: "foo15",    safeName: "foo15"      },
            { identifier: "15foo",    safeName: "%315foo"    },
        ];
        names.forEach(name => expect(toSafeName(name.identifier)).toEqual(name.safeName));
        names.forEach(name => expect(fromSafeName(name.safeName)).toEqual(name.identifier));
    });

    test("special characters hex escaped", () => {
        expect(toSafeName(`hello;,?:@&=+$.!~*'()"/\\ world`)).toEqual(
            `hello%3B%2C%3F%3A%40%26%3D%2B%24%2E%21%7E%2A%27%28%29%22%2F%5C%20world`
        );
    });
});

describe("unescapeName", () => {
    test("no escaping", () => {
        expect(fromSafeName(`_hello-world_`)).toEqual(`_hello-world_`);     
    });

    test("escaped lowercase", () => {
        expect(fromSafeName(`hello%3b%2c%3f%3a%40%26%3d%2b%24%2e%21%7e%2a%27%28%29%22%2f%5c%20world`)).toEqual(
            `hello;,?:@&=+$.!~*'()"/\\ world`
        );
    });

    test("escaped uppercase", () => {
        expect(fromSafeName(`hello%3B%2C%3F%3A%40%26%3D%2B%24%2E%21%7E%2A%27%28%29%22%2F%5C%20world`)).toEqual(
            `hello;,?:@&=+$.!~*'()"/\\ world`
        );
    });
});
