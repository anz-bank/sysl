import { Location, Offset } from "./location";

test.concurrent("toString", () => {
    function loc(startLineOnly: boolean, startLine: number, startCol?: number, endLine?: number, endCol?: number) {
        return new Location("test.sysl", new Offset(startLine, startCol), new Offset(endLine, endCol)).toString(
            startLineOnly
        );
    }

    expect(loc(false, 0)).toBe("test.sysl:1");
    expect(loc(false, 0, 4)).toBe("test.sysl:1:5");
    expect(loc(false, 0, 4, 0, 9)).toBe("test.sysl:1:5::10");
    expect(loc(false, 0, 4, 1, 9)).toBe("test.sysl:1:5:2:10");
    expect(loc(false, 0, 4, 1, 2)).toBe("test.sysl:1:5:2:3");

    expect(loc(true, 0)).toBe("test.sysl:1");
    expect(loc(true, 0, 4)).toBe("test.sysl:1");
    expect(loc(true, 0, 4, 0, 9)).toBe("test.sysl:1");
    expect(loc(true, 0, 4, 1, 9)).toBe("test.sysl:1");
    expect(loc(true, 0, 4, 1, 2)).toBe("test.sysl:1");
});
