import { Location, Offset } from "./location";

function loc(startLine: number, startCol?: number, endLine?: number, endCol?: number): Location {
    return new Location("test.sysl", new Offset(startLine, startCol), new Offset(endLine, endCol));
}

function locStr(startLineOnly: boolean, startLine: number, startCol?: number, endLine?: number, endCol?: number) {
    return loc(startLine, startCol, endLine, endCol).toString(startLineOnly);
}

test.concurrent("toString", () => {
    expect(locStr(false, 0)).toEqual("test.sysl:1");
    expect(locStr(false, 0, 4)).toEqual("test.sysl:1:5");
    expect(locStr(false, 0, 4, 0, 9)).toEqual("test.sysl:1:5::10");
    expect(locStr(false, 0, 4, 1, 9)).toEqual("test.sysl:1:5:2:10");
    expect(locStr(false, 0, 4, 1, 2)).toEqual("test.sysl:1:5:2:3");

    expect(locStr(true, 0)).toEqual("test.sysl:1");
    expect(locStr(true, 0, 4)).toEqual("test.sysl:1");
    expect(locStr(true, 0, 4, 0, 9)).toEqual("test.sysl:1");
    expect(locStr(true, 0, 4, 1, 9)).toEqual("test.sysl:1");
    expect(locStr(true, 0, 4, 1, 2)).toEqual("test.sysl:1");
});

test.concurrent("parse", () => {
    expect(Location.parse("test.sysl:1")).toEqual(loc(0));
    expect(Location.parse("test.sysl:1:5")).toEqual(loc(0, 4));
    expect(Location.parse("test.sysl:1:5::10")).toEqual(loc(0, 4, 0, 9));
    expect(Location.parse("test.sysl:1:5:2:10")).toEqual(loc(0, 4, 1, 9));
    expect(Location.parse("test.sysl:1:5:2:3")).toEqual(loc(0, 4, 1, 2));

    expect(() => Location.parse(" :1")).toThrow();
    expect(() => Location.parse("test.sysl")).toThrow();
    expect(() => Location.parse("test.sysl:")).toThrow();
    expect(() => Location.parse("test.sysl:x")).toThrow();
    expect(() => Location.parse("test.sysl:1:x")).toThrow();
    expect(() => Location.parse("test.sysl:1:2:3")).toThrow();
    expect(() => Location.parse("test.sysl:1:2:3:4:5")).toThrow();
});
