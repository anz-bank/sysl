import "jest-extended";
import { ElementKind, ElementRef } from "./common";

const appRef = ElementRef.parse("Ns::App");
const typeRef = ElementRef.parse("Ns::App.Type");
const fieldRef = ElementRef.parse("Ns::App.Type.Field");

describe("Creating", () => {
    test.concurrent("via constructor, valid", () => {
        expect(new ElementRef([], "App")).toEqual(expect.objectContaining(
            { namespace: [], appName: "App", typeName: "", fieldName: "", kind: ElementKind.App }));
        expect(new ElementRef(["Ns1"], "App")).toEqual(expect.objectContaining(
            { namespace: ["Ns1"], appName: "App", typeName: "", fieldName: "", kind: ElementKind.App }));
        expect(new ElementRef(["Ns"], "App", "Type")).toEqual(expect.objectContaining(
            { namespace: ["Ns"], appName: "App", typeName: "Type", fieldName: "", kind: ElementKind.Type }));
        expect(new ElementRef(["Ns"], "App", "Type", "Field")).toEqual(expect.objectContaining(
            { namespace: ["Ns"], appName: "App", typeName: "Type", fieldName: "Field", kind: ElementKind.Field }));
        expect(new ElementRef(["Ns ."], "App .", "Type .", "Field .")).toEqual(expect.objectContaining(
            { namespace: ["Ns ."], appName: "App .", typeName: "Type .", fieldName: "Field .", kind: ElementKind.Field }));
        expect(new ElementRef(["int"], "date", "bool", "any")).toEqual(expect.objectContaining(
            { namespace: ["int"], appName: "date", typeName: "bool", fieldName: "any", kind: ElementKind.Field }));
        expect(new ElementRef(["1Ns"], "2App", "3Type", "4Field")).toEqual(expect.objectContaining(
            { namespace: ["1Ns"], appName: "2App", typeName: "3Type", fieldName: "4Field", kind: ElementKind.Field }));
        });
    
    test.concurrent("via constructor, invalid", () => {
        expect(() => new ElementRef([""], "App")).toThrow();
        expect(() => new ElementRef(["Ns", ""], "App")).toThrow();
        expect(() => new ElementRef(["", "Ns"], "App")).toThrow();
        expect(() => new ElementRef([], "")).toThrow();
        expect(() => new ElementRef([], "App", "", "Field")).toThrow();
        expect(() => new ElementRef([], "App", undefined, "Field")).toThrow();
    });

    test.concurrent("via parse, valid", () => {
        expect(ElementRef.parse("App")).toEqual(expect.objectContaining(
            { namespace: [], appName: "App", typeName: "", fieldName: "", kind: ElementKind.App }));
        expect(ElementRef.parse("Ns1 :: App")).toEqual(expect.objectContaining(
            { namespace: ["Ns1"], appName: "App", typeName: "", fieldName: "", kind: ElementKind.App }));
        expect(ElementRef.parse("Ns1 :: App.Type")).toEqual(expect.objectContaining(
            { namespace: ["Ns1"], appName: "App", typeName: "Type", fieldName: "", kind: ElementKind.Type }));
        expect(ElementRef.parse("Ns1 :: App.Type.Field")).toEqual(expect.objectContaining(
            { namespace: ["Ns1"], appName: "App", typeName: "Type", fieldName: "Field", kind: ElementKind.Field }));
        expect(ElementRef.parse("  Ns1:: Ns2  ::App .Type . Field     ")).toEqual(expect.objectContaining(
            { namespace: ["Ns1", "Ns2"], appName: "App", typeName: "Type", fieldName: "Field", kind: ElementKind.Field }));
        expect(ElementRef.parse("Ns%20%2E :: App%20%2E.Type%20%2E.Field%20%2E")).toEqual(expect.objectContaining(
            { namespace: ["Ns ."], appName: "App .", typeName: "Type .", fieldName: "Field .", kind: ElementKind.Field }));
        expect(ElementRef.parse("%69nt :: %64ate.%62ool.%61ny")).toEqual(expect.objectContaining(
            { namespace: ["int"], appName: "date", typeName: "bool", fieldName: "any", kind: ElementKind.Field }));
         expect(ElementRef.parse("%31Ns :: %32App.%33Type.%34Field")).toEqual(expect.objectContaining(
            { namespace: ["1Ns"], appName: "2App", typeName: "3Type", fieldName: "4Field", kind: ElementKind.Field }));
        });

    test.concurrent("via parse, invalid", () => {
        expect(() => ElementRef.parse("")).toThrow();
        expect(() => ElementRef.parse(" ")).toThrow();
        expect(() => ElementRef.parse("::")).toThrow();
        expect(() => ElementRef.parse(":: App")).toThrow();
        expect(() => ElementRef.parse("Ns1 :: :: App")).toThrow();
        expect(() => ElementRef.parse("Ns1 : App")).toThrow();
        expect(() => ElementRef.parse("Ns1 ::")).toThrow();
        expect(() => ElementRef.parse("Ns 1 :: App")).toThrow();
        expect(() => ElementRef.parse("App.Ty pe")).toThrow();
        expect(() => ElementRef.parse("int")).toThrow();
        expect(() => ElementRef.parse("App.int")).toThrow();
        expect(() => ElementRef.parse("int :: App")).toThrow();
        expect(() => ElementRef.parse("any.App")).toThrow();
        expect(() => ElementRef.parse("App.")).toThrow();
        expect(() => ElementRef.parse("App..Field")).toThrow();
        expect(() => ElementRef.parse("App.Type.Field.Something")).toThrow();
        expect(() => ElementRef.parse("App.Type::Field")).toThrow();
        expect(() => ElementRef.parse("App.Type.Fi eld")).toThrow();
        expect(() => ElementRef.parse("App.Type.Fi\neld")).toThrow();
        expect(() => ElementRef.parse("1Ns::App.Type.Field")).toThrow();
        expect(() => ElementRef.parse("Ns::2App.Type.Field")).toThrow();
        expect(() => ElementRef.parse("Ns::App.3Type.Field")).toThrow();
        expect(() => ElementRef.parse("Ns::App.Type.4Field")).toThrow();

        expect(ElementRef.tryParse("")).toBeUndefined();
        expect(ElementRef.tryParse(" ")).toBeUndefined();
        expect(ElementRef.tryParse("::")).toBeUndefined();
        expect(ElementRef.tryParse(":: App")).toBeUndefined();
        expect(ElementRef.tryParse("Ns1 :: :: App")).toBeUndefined();
        expect(ElementRef.tryParse("Ns1 : App")).toBeUndefined();
        expect(ElementRef.tryParse("Ns1 ::")).toBeUndefined();
        expect(ElementRef.tryParse("Ns 1 :: App")).toBeUndefined();
        expect(ElementRef.tryParse("App.Ty pe")).toBeUndefined();
        expect(ElementRef.tryParse("int")).toBeUndefined();
        expect(ElementRef.tryParse("App.int")).toBeUndefined();
        expect(ElementRef.tryParse("int :: App")).toBeUndefined();
        expect(ElementRef.tryParse("any.App")).toBeUndefined();
        expect(ElementRef.tryParse("App.")).toBeUndefined();
        expect(ElementRef.tryParse("App..Field")).toBeUndefined();
        expect(ElementRef.tryParse("App.Type.Field.Something")).toBeUndefined();
        expect(ElementRef.tryParse("App.Type::Field")).toBeUndefined();
        expect(ElementRef.tryParse("App.Type.Fi eld")).toBeUndefined();
        expect(ElementRef.tryParse("App.Type.Fi\neld")).toBeUndefined();
        expect(ElementRef.tryParse("1Ns::App.Type.Field")).toBeUndefined();
        expect(ElementRef.tryParse("Ns::2App.Type.Field")).toBeUndefined();
        expect(ElementRef.tryParse("Ns::App.3Type.Field")).toBeUndefined();
        expect(ElementRef.tryParse("Ns::App.Type.4Field")).toBeUndefined();
    });
});

describe("Sysl rendering", () => {
    test.concurrent("standard", () => {
        expect(new ElementRef([], "App").toSysl()).toEqual("App");
        expect(new ElementRef(["Ns1"], "App").toSysl()).toEqual("Ns1 :: App");
        expect(new ElementRef(["Ns1", "Ns2"], "App").toSysl()).toEqual("Ns1 :: Ns2 :: App");
        expect(new ElementRef([], "App", "Type").toSysl()).toEqual("App.Type");
        expect(new ElementRef([], "App", "Type", "Field").toSysl()).toEqual("App.Type.Field");
        expect(new ElementRef(["Ns ."], "App .", "Type .", "Field .").toSysl()).toEqual(
            "Ns%20%2E :: App%20%2E.Type%20%2E.Field%20%2E"
        );
        expect(new ElementRef(["int"], "date", "bool", "any").toSysl()).toEqual(
            "%69nt :: %64ate.%62ool.%61ny"
        );
        expect(new ElementRef(["1Ns"], "2App", "3Type", "4Field").toSysl()).toEqual(
            "%31Ns :: %32App.%33Type.%34Field"
        );
    });

    test.concurrent("compact", () => {
        expect(new ElementRef([], "App").toSysl(true)).toEqual("App");
        expect(new ElementRef(["Ns1"], "App").toSysl(true)).toEqual("Ns1::App");
        expect(new ElementRef(["Ns1", "Ns2"], "App").toSysl(true)).toEqual("Ns1::Ns2::App");
        expect(new ElementRef(["Ns ."], "App .", "Type .", "Field .").toSysl(true)).toEqual(
            "Ns%20%2E::App%20%2E.Type%20%2E.Field%20%2E"
        );
        expect(new ElementRef(["int", "string"], "date", "bool", "any").toSysl(true)).toEqual(
            "%69nt::%73tring::%64ate.%62ool.%61ny"
        );
        expect(new ElementRef(["1Ns", "2Ns"], "3App", "4Type", "5Field").toSysl(true)).toEqual(
            "%31Ns::%32Ns::%33App.%34Type.%35Field"
        );
    });
});

describe("App parts", () => {
    test.concurrent("to", () => {
        expect(new ElementRef([], "App").toAppParts()).toEqual(["App"]);
        expect(new ElementRef(["Ns1"], "App").toAppParts()).toEqual(["Ns1", "App"]);
        expect(new ElementRef(["Ns1", "Ns2"], "App").toAppParts()).toEqual(["Ns1", "Ns2", "App"]);
        expect(new ElementRef([], "App", "Type").toAppParts()).toEqual(["App"]);
        expect(new ElementRef([], "App", "Type", "Field").toAppParts()).toEqual(["App"]);
        expect(new ElementRef(["Ns ."], "App .", "Type .", "Field .").toAppParts()).toEqual(
            ["Ns .", "App ."]
        );
        expect(new ElementRef(["1Ns", "2Ns"], "3App", "4Type", "5Field").toAppParts()).toEqual(
            ["1Ns", "2Ns", "3App"]
        );
    });

    test.concurrent("to safe", () => {
        expect(new ElementRef([], "App").toAppPartsSafe()).toEqual(["App"]);
        expect(new ElementRef(["Ns1"], "App").toAppPartsSafe()).toEqual(["Ns1", "App"]);
        expect(new ElementRef(["Ns1", "Ns2"], "App").toAppPartsSafe()).toEqual(["Ns1", "Ns2", "App"]);
        expect(new ElementRef([], "App", "Type").toAppPartsSafe()).toEqual(["App"]);
        expect(new ElementRef([], "App", "Type", "Field").toAppPartsSafe()).toEqual(["App"]);
        expect(new ElementRef(["Ns ."], "App .", "Type .", "Field .").toAppPartsSafe()).toEqual(
            ["Ns%20%2E", "App%20%2E"]
        );
        expect(new ElementRef(["1Ns", "2Ns"], "3App", "4Type", "5Field").toAppPartsSafe()).toEqual(
            ["%31Ns", "%32Ns", "%33App"]
        );
    });

    test.concurrent("from", () => {
        expect(ElementRef.fromAppParts(["App"]).toSysl()).toEqual("App");
        expect(ElementRef.fromAppParts(["Ns1", "App"]).toSysl()).toEqual("Ns1 :: App");
        expect(ElementRef.fromAppParts(["Ns1", "Ns2", "App"]).toSysl()).toEqual("Ns1 :: Ns2 :: App");
        expect(ElementRef.fromAppParts(["Ns .", "App ."]).toSysl()).toEqual("Ns%20%2E :: App%20%2E");
        expect(ElementRef.fromAppParts(["1Ns", "2App"]).toSysl()).toEqual("%31Ns :: %32App");
    });

    test.concurrent("from safe", () => {
        expect(ElementRef.fromAppPartsSafe(["App"]).toSysl()).toEqual("App");
        expect(ElementRef.fromAppPartsSafe(["Ns1", "App"]).toSysl()).toEqual("Ns1 :: App");
        expect(ElementRef.fromAppPartsSafe(["Ns1", "Ns2", "App"]).toSysl()).toEqual("Ns1 :: Ns2 :: App");
        expect(ElementRef.fromAppPartsSafe(["Ns%20%2E", "App%20%2E"]).toSysl()).toEqual("Ns%20%2E :: App%20%2E");
        expect(ElementRef.fromAppPartsSafe(["%31Ns", "%32App"]).toSysl()).toEqual("%31Ns :: %32App");
    });
});

describe("Genealogy", () => {
    test.concurrent("parent or self", () => {
        const parent = (r: string) => ElementRef.parse(r).toParentOrSelf().toString();
        expect(parent("App")).toEqual("App");
        expect(parent("Ns1::App")).toEqual("Ns1");
        expect(parent("Ns1::Ns2::App")).toEqual("Ns1::Ns2");
        expect(parent("App.Type")).toEqual("App");
        expect(parent("App.Type.Field")).toEqual("App.Type");
    });

    test.concurrent("parent or throw", () => {
        const parent = (r: string) => ElementRef.parse(r).toParent().toString();
        expect(() => parent("App")).toThrow();
        expect(parent("Ns1::App")).toEqual("Ns1");
        expect(parent("Ns1::Ns2::App")).toEqual("Ns1::Ns2");
        expect(parent("App.Type")).toEqual("App");
        expect(parent("App.Type.Field")).toEqual("App.Type");
    });

    const hasDescent = (descendant: string, ancestor: string) =>
        ElementRef.parse(descendant).isDescendantOf(ElementRef.parse(ancestor));

    test.concurrent("is descendant of, without namespace", () => {
        expect(hasDescent("App", "App")).toBeFalse();
        expect(hasDescent("App", "App.Type")).toBeFalse();
        expect(hasDescent("App", "App.Type.Field")).toBeFalse();

        expect(hasDescent("App.Type", "App")).toBeTrue();
        expect(hasDescent("App.Type", "App1")).toBeFalse();
        expect(hasDescent("App.Type", "App.Type")).toBeFalse();
        expect(hasDescent("App.Type", "App.Type.Field")).toBeFalse();

        expect(hasDescent("App.Type.Field", "App")).toBeTrue();
        expect(hasDescent("App.Type.Field", "App1")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.Type")).toBeTrue();
        expect(hasDescent("App.Type.Field", "App.Type1")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.Type.Field")).toBeFalse();
    });

    test.concurrent("is descendant of, with namespace", () => {
        expect(hasDescent("Ns1", "Ns1")).toBeFalse();
        expect(hasDescent("Ns1", "Ns2")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns2")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns3")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns2::App")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns2::App1")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns2::App.Type")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns2::App.Type1")).toBeFalse();
        expect(hasDescent("Ns1", "Ns1::Ns2::App.Type.Field")).toBeFalse();

        expect(hasDescent("Ns1::Ns2", "Ns1")).toBeTrue();
        expect(hasDescent("Ns1::Ns2", "Ns2")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns2")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns3")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns2::App")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns2::App1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns2::App.Type")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns2::App.Type1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2", "Ns1::Ns2::App.Type.Field")).toBeFalse();

        expect(hasDescent("Ns1::Ns2::App", "Ns1")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App", "Ns2")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns2")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns3")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns2::App")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns2::App1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns2::App.Type")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns2::App.Type1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App", "Ns1::Ns2::App.Type.Field")).toBeFalse();

        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns2")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns2")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns3")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns2::App")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns2::App1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns2::App.Type")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns2::App.Type1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type", "Ns1::Ns2::App.Type.Field")).toBeFalse();

        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns2")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns2")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns3")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns2::App")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns2::App1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns2::App.Type")).toBeTrue();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns2::App.Type1")).toBeFalse();
        expect(hasDescent("Ns1::Ns2::App.Type.Field", "Ns1::Ns2::App.Type.Field")).toBeFalse();        
    });
});

describe("Equality", () => {
    test.concurrent("equals", () => {
        expect(ElementRef.parse("App").equals(ElementRef.parse("App"))).toBeTrue();
        expect(ElementRef.parse("App.Type").equals(ElementRef.parse("App.Type"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").equals(ElementRef.parse("App.Type.Field"))).toBeTrue();
        expect(ElementRef.parse("Ns::App").equals(ElementRef.parse("Ns::App"))).toBeTrue();

        expect(ElementRef.parse("App").equals(ElementRef.parse("App1"))).toBeFalse();
        expect(ElementRef.parse("App").equals(ElementRef.parse("App.Type"))).toBeFalse();
        expect(ElementRef.parse("App.Type").equals(ElementRef.parse("App"))).toBeFalse();
        expect(ElementRef.parse("App.Type").equals(ElementRef.parse("App.Type1"))).toBeFalse();
        expect(ElementRef.parse("App.Type").equals(ElementRef.parse("App1.Type"))).toBeFalse();
        expect(ElementRef.parse("App.Type").equals(ElementRef.parse("App.Type.Field"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").equals(ElementRef.parse("App.Type.Field1"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").equals(ElementRef.parse("App.Type1.Field"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").equals(ElementRef.parse("App1.Type.Field"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").equals(ElementRef.parse("App.Type"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").equals(ElementRef.parse("App"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").equals(ElementRef.parse("Ns::App1"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").equals(ElementRef.parse("Ns1::App"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").equals(ElementRef.parse("Ns::Ns::App"))).toBeFalse();
    });

    test.concurrent("types equal", () => {
        expect(ElementRef.parse("App").typesEqual(ElementRef.parse("App"))).toBeTrue();
        expect(ElementRef.parse("App.Type").typesEqual(ElementRef.parse("App.Type"))).toBeTrue();
        expect(ElementRef.parse("App.Type").typesEqual(ElementRef.parse("App.Type.Field"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").typesEqual(ElementRef.parse("App.Type"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").typesEqual(ElementRef.parse("App.Type.Field"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").typesEqual(ElementRef.parse("App.Type.Field1"))).toBeTrue();

        expect(ElementRef.parse("App").typesEqual(ElementRef.parse("App.Type"))).toBeFalse();
        expect(ElementRef.parse("App.Type").typesEqual(ElementRef.parse("App"))).toBeFalse();
        expect(ElementRef.parse("App.Type").typesEqual(ElementRef.parse("App.Type1"))).toBeFalse();
        expect(ElementRef.parse("App.Type").typesEqual(ElementRef.parse("App1.Type"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").typesEqual(ElementRef.parse("App.Type1.Field"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").typesEqual(ElementRef.parse("App1.Type.Field"))).toBeFalse();
    });

    test.concurrent("apps equal", () => {
        expect(ElementRef.parse("App").appsEqual(ElementRef.parse("App"))).toBeTrue();
        expect(ElementRef.parse("App").appsEqual(ElementRef.parse("App.Type"))).toBeTrue();
        expect(ElementRef.parse("App.Type").appsEqual(ElementRef.parse("App"))).toBeTrue();
        expect(ElementRef.parse("App.Type").appsEqual(ElementRef.parse("App.Type"))).toBeTrue();
        expect(ElementRef.parse("App.Type").appsEqual(ElementRef.parse("App.Type1"))).toBeTrue();
        expect(ElementRef.parse("App.Type").appsEqual(ElementRef.parse("App.Type.Field"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").appsEqual(ElementRef.parse("App.Type.Field"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").appsEqual(ElementRef.parse("App.Type.Field1"))).toBeTrue();
        expect(ElementRef.parse("App.Type.Field").appsEqual(ElementRef.parse("App.Type1.Field"))).toBeTrue();
        expect(ElementRef.parse("Ns::App").appsEqual(ElementRef.parse("Ns::App"))).toBeTrue();

        expect(ElementRef.parse("App").appsEqual(ElementRef.parse("App1"))).toBeFalse();
        expect(ElementRef.parse("App.Type").appsEqual(ElementRef.parse("App1.Type"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").appsEqual(ElementRef.parse("App1.Type.Field"))).toBeFalse();
        expect(ElementRef.parse("App.Type.Field").appsEqual(ElementRef.parse("App.Type"))).toBeTrue();
        expect(ElementRef.parse("Ns::App").appsEqual(ElementRef.parse("App"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").appsEqual(ElementRef.parse("Ns::App1"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").appsEqual(ElementRef.parse("Ns1::App"))).toBeFalse();
        expect(ElementRef.parse("Ns::App").appsEqual(ElementRef.parse("Ns::Ns::App"))).toBeFalse();
    });
});

describe("copying",() => {
    test.concurrent("with", () => {
        expect(appRef.with({ typeName: "Type" })).toEqual(typeRef);
        expect(appRef.with({ typeName: "Type", fieldName: "Field" })).toEqual(fieldRef);
        expect(() => appRef.with({ fieldName: "Field" })).toThrow();

        expect(fieldRef.with({})).toEqual(fieldRef);
        expect(fieldRef.with({ namespace: undefined, appName: undefined, typeName: undefined, fieldName: undefined }))
            .toEqual(fieldRef);
        
        expect(fieldRef.with({ namespace: ["Ns1"] }).toString()).toEqual("Ns1::App.Type.Field");
        expect(fieldRef.with({ appName: "App1" }).toString()).toEqual("Ns::App1.Type.Field");
        expect(fieldRef.with({ typeName: "Type1" }).toString()).toEqual("Ns::App.Type1.Field");
        expect(fieldRef.with({ fieldName: "Field1" }).toString()).toEqual("Ns::App.Type.Field1");
        expect(fieldRef.with({
            namespace: ["Ns1"],
            appName: "App1",
            typeName: "Type1",
            fieldName: "Field1",
         }).toString()).toEqual("Ns1::App1.Type1.Field1");

        expect(fieldRef.with({ namespace: [] }).toString()).toEqual("App.Type.Field");
        expect(fieldRef.with({ fieldName: "" })).toEqual(typeRef);
        expect(fieldRef.with({ typeName: "", fieldName: "" })).toEqual(appRef);

        expect(() => fieldRef.with({ namespace: [""] })).toThrow();
        expect(() => fieldRef.with({ appName: "" })).toThrow();
        expect(() => fieldRef.with({ typeName: "" })).toThrow();
    });

    test.concurrent("truncate", () => {
        expect(appRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(appRef.truncate(ElementKind.Type)).toEqual(appRef);
        expect(appRef.truncate(ElementKind.Field)).toEqual(appRef);
        expect(() => appRef.truncate(ElementKind.Endpoint)).toThrow();
        expect(() => appRef.truncate(ElementKind.Statement)).toThrow();
        expect(() => appRef.truncate(ElementKind.Parameter)).toThrow();

        expect(typeRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(typeRef.truncate(ElementKind.Type)).toEqual(typeRef);
        expect(typeRef.truncate(ElementKind.Field)).toEqual(typeRef);
        expect(() => typeRef.truncate(ElementKind.Endpoint)).toThrow();
        expect(() => typeRef.truncate(ElementKind.Statement)).toThrow();
        expect(() => typeRef.truncate(ElementKind.Parameter)).toThrow();

        expect(fieldRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(fieldRef.truncate(ElementKind.Type)).toEqual(typeRef);
        expect(fieldRef.truncate(ElementKind.Field)).toEqual(fieldRef);
        expect(() => fieldRef.truncate(ElementKind.Endpoint)).toThrow();
        expect(() => fieldRef.truncate(ElementKind.Statement)).toThrow();
        expect(() => fieldRef.truncate(ElementKind.Parameter)).toThrow();    });

    test.concurrent("clone", () => {
        expect(appRef.clone()).toEqual(appRef);
        expect(typeRef.clone()).toEqual(typeRef);
        expect(fieldRef.clone()).toEqual(fieldRef);
    });
});

describe("computed properties", () => {
    test.concurrent("x", () => {
        expect(appRef).toEqual(expect.objectContaining(
            { isApp: true, isType: false, isField: false, numericKind: 1 }));
        expect(typeRef).toEqual(expect.objectContaining(
            { isApp: false, isType: true, isField: false, numericKind: 2 }));
        expect(fieldRef).toEqual(expect.objectContaining(
            { isApp: false, isType: false, isField: true, numericKind: 3 }));
    });
});
