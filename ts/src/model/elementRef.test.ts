import "jest-extended";
import { ElementRef } from "./elementRef";
import { ElementKind } from "./elementKind";

const appRef = ElementRef.parse("Ns::App");
const typeRef = ElementRef.parse("Ns::App.Type");
const fieldRef = ElementRef.parse("Ns::App.Type.Field");
const endpointRef = ElementRef.parse("Ns::App.[Endpoint]");
const statementRef = ElementRef.parse("Ns::App.[Endpoint].[0]");

function expectParts(ref: string | ElementRef, ...expected: any[]) {
    const r = ref instanceof ElementRef ? ref : ElementRef.parse(ref);
    const actual = [
        r.namespace,
        r.appName,
        r.typeName || r.endpointName,
        r.endpointName ? r.statementIndices : r.fieldName,
        r.kind];
    expect(actual).toEqual(expected);
}

describe("Constructor", () => {
    test.concurrent("data elements, valid", () => {
        expectParts(new ElementRef([], "App"), [], "App", "", "", ElementKind.App);
        expectParts(new ElementRef(["Ns1"], "App"), ["Ns1"], "App", "", "", ElementKind.App);
        expectParts(new ElementRef(["Ns"], "App", "Type"), ["Ns"], "App", "Type", "", ElementKind.Type);
        expectParts(new ElementRef(["Ns"], "App", "Type", "Field"), ["Ns"], "App", "Type", "Field", ElementKind.Field);
        expectParts(new ElementRef(["Ns ."], "App .", "Type .", "Field ."), ["Ns ."], "App .", "Type .", "Field .", ElementKind.Field);
        expectParts(new ElementRef(["int"], "date", "bool", "any"), ["int"], "date", "bool", "any", ElementKind.Field);
        expectParts(new ElementRef(["1Ns"], "2App", "3Type", "4Field"), ["1Ns"], "2App", "3Type", "4Field", ElementKind.Field);
    });
    
    test.concurrent("data elements, invalid", () => {
        let err = "All namespace parts must be non-empty";
        expect(() => new ElementRef([""], "App")).toThrow(err);
        expect(() => new ElementRef(["Ns", ""], "App")).toThrow(err);
        expect(() => new ElementRef(["", "Ns"], "App")).toThrow(err);

        expect(() => new ElementRef([], "")).toThrow("appName must not be empty.");
        expect(() => new ElementRef([], "App", "", "Field")).toThrow("typeName must be specified");
    });

    test.concurrent("behavior elements, valid", () => {
        expectParts(new ElementRef(["Ns"], "App", "", "", "Endpoint"), ["Ns"], "App", "Endpoint", [], ElementKind.Endpoint);
        expectParts(new ElementRef(["Ns"], "App", "", "", "Endpoint", [0, 3]), ["Ns"], "App", "Endpoint", [0, 3], ElementKind.Statement);
        expectParts(new ElementRef(["Ns ."], "App .", "", "", "Endpoint .", [0, 3]), ["Ns ."], "App .", "Endpoint .", [0, 3], ElementKind.Statement);
        expectParts(new ElementRef(["int"], "date", "", "", "bool", [0, 3]), ["int"], "date", "bool", [0, 3], ElementKind.Statement);
        expectParts(new ElementRef(["1Ns"], "2App", "", "", "3Endpoint", [0, 3]), ["1Ns"], "2App", "3Endpoint", [0, 3], ElementKind.Statement);
    });

    test.concurrent("behavior elements, invalid", () => {
        expect(() => new ElementRef([], "App", "", "", "", [0, 3])).toThrow("endpointName must be specified");
        expect(() => new ElementRef([], "App", "Type", "", "Endpoint")).toThrow("endpointName must not be specified");

        const err = "endpointName or statementIndices must not be specified";
        expect(() => new ElementRef([], "App", "Type", "Field", "Endpoint")).toThrow(err);
        expect(() => new ElementRef([], "App", "Type", "", "Endpoint", [0, 3])).toThrow(err);
        expect(() => new ElementRef([], "App", "Type", "Field", "Endpoint", [0, 3])).toThrow(err);
        expect(() => new ElementRef([], "App", "Type", "", "", [0, 3])).toThrow(err);
        expect(() => new ElementRef([], "App", "Type", "Field", "", [0, 3])).toThrow(err);
    });
});

describe("Parse", () => {
    function expectInvalid(strOrFactory: string, error?: string) {
        expect(() => ElementRef.parse(strOrFactory)).toThrow(error);
        expect(ElementRef.tryParse(strOrFactory)).toBeUndefined();
    }

    test.concurrent("data elements, valid", () => {
        expectParts("App", [], "App", "", "", ElementKind.App);
        expectParts("Ns1 :: App", ["Ns1"], "App", "", "", ElementKind.App);
        expectParts("Ns1 :: App.Type", ["Ns1"], "App", "Type", "", ElementKind.Type);
        expectParts("Ns1 :: App.Type.Field", ["Ns1"], "App", "Type", "Field", ElementKind.Field);
        expectParts("  Ns1:: Ns2  ::App .Type . Field     ", ["Ns1", "Ns2"], "App", "Type", "Field", ElementKind.Field);
        expectParts("Ns%20%2E :: App%20%2E.Type%20%2E.Field%20%2E", ["Ns ."], "App .", "Type .", "Field .", ElementKind.Field);
        expectParts("%69nt :: %64ate.%62ool.%61ny", ["int"], "date", "bool", "any", ElementKind.Field);
        expectParts("%31Ns :: %32App.%33Type.%34Field", ["1Ns"], "2App", "3Type", "4Field", ElementKind.Field);
    });

    test.concurrent("data elements, invalid", () => {
        let err = "Disallowed characters in app";
        expectInvalid("", err);
        expectInvalid(" ", err);
        expectInvalid("::", err);
        expectInvalid(":: App", err);
        expectInvalid("Ns1 :: :: App", err);
        expectInvalid("Ns1 : App", err);
        expectInvalid("Ns1 ::", err);
        expectInvalid("Ns 1 :: App", err);
        expectInvalid("int", err);
        expectInvalid("int :: App", err);
        expectInvalid("any.App", err);
        expectInvalid("1Ns::App.Type.Field", err);
        expectInvalid("Ns::2App.Type.Field", err);

        err = "Disallowed characters in type/field name.";
        expectInvalid("App.Ty pe", err);
        expectInvalid("App.int", err);
        expectInvalid("App.", err);
        expectInvalid("App..Field", err);
        expectInvalid("App.Type::Field", err);
        expectInvalid("App.Type.Fi eld", err);
        expectInvalid("App.Type.Fi\neld", err);
        expectInvalid("Ns::App.3Type.Field", err);
        expectInvalid("Ns::App.Type.4Field", err);

        err = "Too many dots.";
        expectInvalid("App.Type.Field.Something", err);
        expectInvalid("App.Type..Field", err);
    });

    test.concurrent("behavior elements, valid", () => {
        expectParts("Ns1 :: App.[Endpoint]", ["Ns1"], "App", "Endpoint", [], ElementKind.Endpoint);
        expectParts("Ns1 :: App.[Endpoint].[1,0,3]", ["Ns1"], "App", "Endpoint", [1, 0, 3], ElementKind.Statement);

        // Insignificant whitespace
        expectParts("Ns1 :: App.  [Endpoint]  ", ["Ns1"], "App", "Endpoint", [], ElementKind.Endpoint);
        expectParts("Ns1 :: App.[Endpoint]  .[1,0,3]", ["Ns1"], "App", "Endpoint", [1, 0, 3], ElementKind.Statement);
        expectParts("Ns1 :: App.[Endpoint].  [1,0,3]", ["Ns1"], "App", "Endpoint", [1, 0, 3], ElementKind.Statement);
        expectParts("Ns1 :: App . [Endpoint] . [ 1 ,  0 , 3 ]  ", ["Ns1"], "App", "Endpoint", [1, 0, 3], ElementKind.Statement);

        // Significant chars
        expectParts("Ns1 :: App.[Hello World]", ["Ns1"], "App", "Hello World", [], ElementKind.Endpoint);
        expectParts("Ns1 :: App.[/api/customer/{string}]", ["Ns1"], "App", "/api/customer/{string}", [], ElementKind.Endpoint);
        expectParts("Ns1 :: App.[/api/customer[example.com%5d/{string}].[1]", ["Ns1"], "App", "/api/customer[example.com]/{string}", [1], ElementKind.Statement);
        expectParts("Ns1 :: App.[int]", ["Ns1"], "App", "int", [], ElementKind.Endpoint);
    });

    test.concurrent("behavior elements, invalid", () => {
        expectInvalid("App.[Endpoint", "Missing closing bracket");

        let err = "Unexpected characters";
        expectInvalid("App.[Endpoint]x", err);
        expectInvalid("App.x[Endpoint]", err);
        expectInvalid("App.Type.[4]", err);

        err = "Expected open brackets";
        expectInvalid("App.[Endpoint].", err);
        expectInvalid("App.[Endpoint] . ", err);
 
        err = "Invalid statement indices";
        expectInvalid("App.[Endpoint].[,]", err);
        expectInvalid("App.[Endpoint].[a]", err);
        expectInvalid("App.[Endpoint].[1,]", err);
        expectInvalid("App.[Endpoint].[,2]", err);
        expectInvalid("App.[Endpoint].[1,,3]", err);
        expectInvalid("App.[Endpoint].[1,b,3]", err);
        expectInvalid("App.[Endpoint].[1,-2,3]", err);

        err = "No content";
        expectInvalid("App.[Endpoint].[]", err);
        expectInvalid("App.[Endpoint].[  ]", err);
        expectInvalid("App.[].[1,2,3]", err);
        expectInvalid("App.[  ].[1,2,3]", err);
    });
});

describe("String rendering", () => {
    test.concurrent("data elements, standard", () => {
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

    test.concurrent("data elements, compact", () => {
        expect(new ElementRef([], "App").toString()).toEqual("App");
        expect(new ElementRef(["Ns1"], "App").toString()).toEqual("Ns1::App");
        expect(new ElementRef(["Ns1", "Ns2"], "App").toString()).toEqual("Ns1::Ns2::App");
        expect(new ElementRef(["Ns ."], "App .", "Type .", "Field .").toString()).toEqual(
            "Ns%20%2E::App%20%2E.Type%20%2E.Field%20%2E"
        );
        expect(new ElementRef(["int", "string"], "date", "bool", "any").toString()).toEqual(
            "%69nt::%73tring::%64ate.%62ool.%61ny"
        );
        expect(new ElementRef(["1Ns", "2Ns"], "3App", "4Type", "5Field").toString()).toEqual(
            "%31Ns::%32Ns::%33App.%34Type.%35Field"
        );
    });

    test.concurrent("behavior elements, toString", () => {
        expect(new ElementRef(["Ns"], "App", "", "", "Endpoint").toString()).toEqual("Ns::App.[Endpoint]");
        expect(new ElementRef(["Ns"], "App", "", "", "Endpoint", [0, 3]).toString()).toEqual("Ns::App.[Endpoint].[0,3]");
        expect(new ElementRef(["[ Ns ]"], "[ App ]", "", "", "[ Endpoint ]", [0, 3]).toString()).toEqual("%5B%20Ns%20%5D::%5B%20App%20%5D.[[ Endpoint %5D].[0,3]");
        expect(new ElementRef(["int"], "date", "", "", "bool", [0, 3]).toString()).toEqual("%69nt::%64ate.[bool].[0,3]");
        expect(new ElementRef(["1Ns"], "2App", "", "", "3Endpoint", [0, 3]).toString()).toEqual("%31Ns::%32App.[3Endpoint].[0,3]");
    });

    test.concurrent("behavior elements, toSysl", () => {
        expect(() => endpointRef.toSysl()).toThrow();
        expect(() => statementRef.toSysl()).toThrow();
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

    test.concurrent("parent or undefined", () => {
        const parent = (r: string) => ElementRef.parse(r).toParentOrUndefined()?.toString();

        expect(parent("App")).toBeUndefined();
        expect(parent("Ns1::App")).toEqual("Ns1");
        expect(parent("Ns1::Ns2::App")).toEqual("Ns1::Ns2");
        expect(parent("App.Type")).toEqual("App");
        expect(parent("App.Type.Field")).toEqual("App.Type");

        expect(parent("App.[Endpoint]")).toEqual("App");
        expect(parent("App.[Endpoint].[0]")).toEqual("App.[Endpoint]");
        expect(parent("App.[Endpoint].[0,1]")).toEqual("App.[Endpoint].[0]");
        expect(parent("App.[Endpoint].[0,1,2]")).toEqual("App.[Endpoint].[0,1]");
    });

    test.concurrent("parent or throw", () => {
        const parent = (r: string) => ElementRef.parse(r).toParent()?.toString();

        expect(() => parent("App")).toThrow();
        expect(parent("Ns1::App")).toEqual("Ns1");
        expect(parent("Ns1::Ns2::App")).toEqual("Ns1::Ns2");
        expect(parent("App.Type")).toEqual("App");
        expect(parent("App.Type.Field")).toEqual("App.Type");

        expect(parent("App.[Endpoint]")).toEqual("App");
        expect(parent("App.[Endpoint].[0]")).toEqual("App.[Endpoint]");
        expect(parent("App.[Endpoint].[0,1]")).toEqual("App.[Endpoint].[0]");
        expect(parent("App.[Endpoint].[0,1,2]")).toEqual("App.[Endpoint].[0,1]");
    });

    const hasDescent = (descendant: string, ancestor: string) =>
        ElementRef.parse(descendant).isDescendantOf(ElementRef.parse(ancestor));

    test.concurrent("is descendant of, without namespace", () => {
        expect(hasDescent("App", "App")).toBeFalse();
        expect(hasDescent("App", "App.Type")).toBeFalse();
        expect(hasDescent("App", "App.Type.Field")).toBeFalse();
        expect(hasDescent("App", "App.[Endpoint]")).toBeFalse();
        expect(hasDescent("App", "App.[Endpoint].[0]")).toBeFalse();
        expect(hasDescent("App", "App.[Endpoint].[0,1]")).toBeFalse();

        expect(hasDescent("App.Type", "App")).toBeTrue();
        expect(hasDescent("App.Type", "App1")).toBeFalse();
        expect(hasDescent("App.Type", "App.Type")).toBeFalse();
        expect(hasDescent("App.Type", "App.Type.Field")).toBeFalse();
        expect(hasDescent("App.Type", "App.[Endpoint]")).toBeFalse();
        expect(hasDescent("App.Type", "App.[Endpoint].[0]")).toBeFalse();
        expect(hasDescent("App.Type", "App.[Endpoint].[0,1]")).toBeFalse();

        expect(hasDescent("App.Type.Field", "App")).toBeTrue();
        expect(hasDescent("App.Type.Field", "App1")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.Type")).toBeTrue();
        expect(hasDescent("App.Type.Field", "App.Type1")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.Type.Field")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.[Endpoint]")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.[Endpoint].[0]")).toBeFalse();
        expect(hasDescent("App.Type.Field", "App.[Endpoint].[0,1]")).toBeFalse();

        expect(hasDescent("App.[Endpoint]", "App")).toBeTrue();
        expect(hasDescent("App.[Endpoint]", "App1")).toBeFalse();
        expect(hasDescent("App.[Endpoint]", "App.Type")).toBeFalse();
        expect(hasDescent("App.[Endpoint]", "App.Type.Field")).toBeFalse();
        expect(hasDescent("App.[Endpoint]", "App.[Endpoint]")).toBeFalse();
        expect(hasDescent("App.[Endpoint]", "App.[Endpoint].[0]")).toBeFalse();
        expect(hasDescent("App.[Endpoint]", "App.[Endpoint].[0,1]")).toBeFalse();

        expect(hasDescent("App.[Endpoint].[0]", "App")).toBeTrue();
        expect(hasDescent("App.[Endpoint].[0]", "App1")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0]", "App.Type")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0]", "App.Type.Field")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0]", "App.[Endpoint]")).toBeTrue();
        expect(hasDescent("App.[Endpoint].[0]", "App.[Endpoint1]")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0]", "App.[Endpoint].[0]")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0]", "App.[Endpoint].[0,1]")).toBeFalse();

        expect(hasDescent("App.[Endpoint].[0,1]", "App")).toBeTrue();
        expect(hasDescent("App.[Endpoint].[0,1]", "App1")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.Type")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.Type.Field")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.[Endpoint]")).toBeTrue();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.[Endpoint1]")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.[Endpoint].[0]")).toBeTrue();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.[Endpoint].[5]")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0,1]", "App.[Endpoint].[0,1]")).toBeFalse();
        expect(hasDescent("App.[Endpoint].[0,1,2]", "App.[Endpoint].[0,1]")).toBeTrue();
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
        const expectEqual = (first: string, second: string) => expect(ElementRef.parse(first).equals(ElementRef.parse(second)));

        expectEqual("App", "App").toBeTrue();
        expectEqual("App.Type", "App.Type").toBeTrue();
        expectEqual("App.Type.Field", "App.Type.Field").toBeTrue();
        expectEqual("Ns::App", "Ns::App").toBeTrue();
        expectEqual("App.[Endpoint]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0,1]").toBeTrue();

        expectEqual("App", "App1").toBeFalse();
        expectEqual("App", "App.Type").toBeFalse();
        expectEqual("App", "App.Type.Field").toBeFalse();
        expectEqual("App", "App.[Endpoint]").toBeFalse();
        expectEqual("App", "App.[Endpoint].[0]").toBeFalse();
        expectEqual("App", "App.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.Type", "App").toBeFalse();
        expectEqual("App.Type", "App.Type1").toBeFalse();
        expectEqual("App.Type", "App1.Type").toBeFalse();
        expectEqual("App.Type", "App.Type.Field").toBeFalse();
        expectEqual("App.Type", "App.[Type]").toBeFalse();
        expectEqual("App.Type", "App.[Type].[0]").toBeFalse();
        expectEqual("App.Type", "App.[Type].[0,1]").toBeFalse();
        expectEqual("App.Type.Field", "App").toBeFalse();
        expectEqual("App.Type.Field", "App.Type").toBeFalse();
        expectEqual("App.Type.Field", "App.Type.Field1").toBeFalse();
        expectEqual("App.Type.Field", "App.Type1.Field").toBeFalse();
        expectEqual("App.Type.Field", "App1.Type.Field").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type]").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type].[0]").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type].[0,1]").toBeFalse();
        expectEqual("App.[Endpoint]", "App").toBeFalse();
        expectEqual("App.[Endpoint]", "App.Endpoint").toBeFalse();
        expectEqual("App.[Endpoint]", "App.Endpoint.Field").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint1]").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint].[0]").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.Endpoint").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.Endpoint.Field").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[5]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App.Endpoint").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App.Endpoint.Field").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0,5]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0,1,2]").toBeFalse();
        expectEqual("Ns::App", "App").toBeFalse();
        expectEqual("Ns::App", "Ns::App1").toBeFalse();
        expectEqual("Ns::App", "Ns1::App").toBeFalse();
        expectEqual("Ns::App", "Ns::Ns::App").toBeFalse();
    });

    test.concurrent("types equal", () => {
        const expectEqual = (first: string, second: string) => expect(ElementRef.parse(first).typesEqual(ElementRef.parse(second)));

        expectEqual("App", "App").toBeTrue();
        expectEqual("App.Type", "App.Type").toBeTrue();
        expectEqual("App.Type", "App.Type.Field").toBeTrue();
        expectEqual("App.Type.Field", "App.Type").toBeTrue();
        expectEqual("App.Type.Field", "App.Type.Field").toBeTrue();
        expectEqual("App.Type.Field", "App.Type.Field1").toBeTrue();

        expectEqual("App", "App.Type").toBeFalse();
        expectEqual("App", "App.Type.Field").toBeFalse();
        expectEqual("App", "App.[Endpoint]").toBeFalse();
        expectEqual("App", "App.[Endpoint].[0]").toBeFalse();
        expectEqual("App.Type", "App").toBeFalse();
        expectEqual("App.Type", "App.Type1").toBeFalse();
        expectEqual("App.Type", "App1.Type").toBeFalse();
        expectEqual("App.Type", "App.[Type]").toBeFalse();
        expectEqual("App.Type", "App.[Type].[0]").toBeFalse();
        expectEqual("App.Type.Field", "App.Type1.Field").toBeFalse();
        expectEqual("App.Type.Field", "App1.Type.Field").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type]").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type].[0]").toBeFalse();
        expectEqual("App.[Endpoint]", "App").toBeFalse();
        expectEqual("App.[Endpoint]", "App.Endpoint").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint].[0]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0]").toBeFalse();
    });

    test.concurrent("endpoints equal", () => {
        const expectEqual = (first: string, second: string) => expect(ElementRef.parse(first).endpointsEqual(ElementRef.parse(second)));

        expectEqual("App", "App").toBeTrue();
        expectEqual("App.[Endpoint]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[5]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0,1]").toBeTrue();

        expectEqual("App", "App.Type").toBeFalse();
        expectEqual("App", "App.Type.Field").toBeFalse();
        expectEqual("App", "App.[Endpoint]").toBeFalse();
        expectEqual("App", "App.[Endpoint].[0]").toBeFalse();
        expectEqual("App.Type", "App").toBeFalse();
        expectEqual("App.Type", "App.Type").toBeFalse();
        expectEqual("App.Type", "App.Type.Field").toBeFalse();
        expectEqual("App.Type", "App.[Type]").toBeFalse();
        expectEqual("App.Type", "App.[Type].[0]").toBeFalse();
        expectEqual("App.Type.Field", "App.Type").toBeFalse();
        expectEqual("App.Type.Field", "App.Type.Field").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type]").toBeFalse();
        expectEqual("App.Type.Field", "App.[Type].[0]").toBeFalse();
        expectEqual("App.[Endpoint]", "App").toBeFalse();
        expectEqual("App.[Endpoint]", "App.Type").toBeFalse();
        expectEqual("App.[Endpoint]", "App.Type.Field").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint1]").toBeFalse();
        expectEqual("App.[Endpoint]", "App.[Endpoint1].[0]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.Type").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.Type.Field").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint1]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint1].[0]").toBeFalse();
    });

    test.concurrent("apps equal", () => {
        const expectEqual = (first: string, second: string) => expect(ElementRef.parse(first).appsEqual(ElementRef.parse(second)));

        expectEqual("App", "App").toBeTrue();
        expectEqual("App", "App.Type").toBeTrue();
        expectEqual("App", "App.Type.Field").toBeTrue();
        expectEqual("App", "App.[Endpoint]").toBeTrue();
        expectEqual("App", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App", "App.[Endpoint].[0,1]").toBeTrue();
        expectEqual("App.Type", "App").toBeTrue();
        expectEqual("App.Type", "App.Type").toBeTrue();
        expectEqual("App.Type", "App.Type1").toBeTrue();
        expectEqual("App.Type", "App.Type.Field").toBeTrue();
        expectEqual("App.Type", "App.[Endpoint]").toBeTrue();
        expectEqual("App.Type", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.Type", "App.[Endpoint].[0,1]").toBeTrue();
        expectEqual("App.Type.Field", "App").toBeTrue();
        expectEqual("App.Type.Field", "App.Type").toBeTrue();
        expectEqual("App.Type.Field", "App.Type1").toBeTrue();
        expectEqual("App.Type.Field", "App.Type.Field").toBeTrue();
        expectEqual("App.Type.Field", "App.Type.Field1").toBeTrue();
        expectEqual("App.Type.Field", "App.Type1.Field").toBeTrue();
        expectEqual("App.Type.Field", "App.[Endpoint]").toBeTrue();
        expectEqual("App.Type.Field", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.Type.Field", "App.[Endpoint].[0,1]").toBeTrue();
        expectEqual("App.[Endpoint]", "App").toBeTrue();
        expectEqual("App.[Endpoint]", "App.Type").toBeTrue();
        expectEqual("App.[Endpoint]", "App.Type.Field").toBeTrue();
        expectEqual("App.[Endpoint]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint]", "App.[Endpoint1]").toBeTrue();
        expectEqual("App.[Endpoint]", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.[Endpoint]", "App.[Endpoint].[0,1]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.Type").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.Type.Field").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint1]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[5]").toBeTrue();
        expectEqual("App.[Endpoint].[0]", "App.[Endpoint].[0,1]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.Type").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.Type.Field").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint1]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[5]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0,1]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0,5]").toBeTrue();
        expectEqual("App.[Endpoint].[0,1]", "App.[Endpoint].[0,1,2]").toBeTrue();
        expectEqual("Ns::App", "Ns::App").toBeTrue();

        expectEqual("App", "App1").toBeFalse();
        expectEqual("App", "App1.Type").toBeFalse();
        expectEqual("App", "App1.Type.Field").toBeFalse();
        expectEqual("App", "App1.[Endpoint]").toBeFalse();
        expectEqual("App", "App1.[Endpoint].[0]").toBeFalse();
        expectEqual("App", "App1.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.Type", "App1").toBeFalse();
        expectEqual("App.Type", "App1.Type").toBeFalse();
        expectEqual("App.Type", "App1.Type.Field").toBeFalse();
        expectEqual("App.Type", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.Type", "App1.[Endpoint].[0]").toBeFalse();
        expectEqual("App.Type", "App1.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.Type.Field", "App1").toBeFalse();
        expectEqual("App.Type.Field", "App1.Type").toBeFalse();
        expectEqual("App.Type.Field", "App1.Type.Field").toBeFalse();
        expectEqual("App.Type.Field", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.Type.Field", "App1.[Endpoint].[0]").toBeFalse();
        expectEqual("App.Type.Field", "App1.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.[Endpoint]", "App1").toBeFalse();
        expectEqual("App.[Endpoint]", "App1.Type").toBeFalse();
        expectEqual("App.[Endpoint]", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint]", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint]", "App1.[Endpoint].[0]").toBeFalse();
        expectEqual("App.[Endpoint]", "App1.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App1").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App1.Type").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App1.[Endpoint].[0]").toBeFalse();
        expectEqual("App.[Endpoint].[0]", "App1.[Endpoint].[0,1]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App1").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App1.Type").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App1.[Endpoint]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App1.[Endpoint].[0]").toBeFalse();
        expectEqual("App.[Endpoint].[0,1]", "App1.[Endpoint].[0,1]").toBeFalse();
        expectEqual("Ns::App", "App").toBeFalse();
        expectEqual("Ns::App", "Ns::App1").toBeFalse();
        expectEqual("Ns::App", "Ns1::App").toBeFalse();
        expectEqual("Ns::App", "Ns::Ns::App").toBeFalse();
    });
});

describe("Copying",() => {
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

        expect(appRef.with({ endpointName: "Endpoint" })).toEqual(endpointRef);
        expect(appRef.with({ endpointName: "Endpoint", statementIndices: [0] })).toEqual(statementRef);
        expect(endpointRef.with({ statementIndices: [0] })).toEqual(statementRef);

        expect(() => appRef.with({ statementIndices: [0] })).toThrow();
        expect(() => typeRef.with({ statementIndices: [0] })).toThrow();
        expect(() => typeRef.with({ endpointName: "Endpoint" })).toThrow();
    });

    test.concurrent("truncate", () => {
        expect(appRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(appRef.truncate(ElementKind.Type)).toEqual(appRef);
        expect(appRef.truncate(ElementKind.Field)).toEqual(appRef);
        expect(appRef.truncate(ElementKind.Endpoint)).toEqual(appRef);
        expect(appRef.truncate(ElementKind.Statement)).toEqual(appRef);

        expect(typeRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(typeRef.truncate(ElementKind.Type)).toEqual(typeRef);
        expect(typeRef.truncate(ElementKind.Field)).toEqual(typeRef);
        expect(() => typeRef.truncate(ElementKind.Endpoint)).toThrow();
        expect(() => typeRef.truncate(ElementKind.Statement)).toThrow();

        expect(fieldRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(fieldRef.truncate(ElementKind.Type)).toEqual(typeRef);
        expect(fieldRef.truncate(ElementKind.Field)).toEqual(fieldRef);
        expect(() => fieldRef.truncate(ElementKind.Endpoint)).toThrow();
        expect(() => fieldRef.truncate(ElementKind.Statement)).toThrow();

        expect(endpointRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(() => endpointRef.truncate(ElementKind.Type)).toThrow();
        expect(() => endpointRef.truncate(ElementKind.Field)).toThrow();
        expect(endpointRef.truncate(ElementKind.Endpoint)).toEqual(endpointRef);
        expect(endpointRef.truncate(ElementKind.Statement)).toEqual(endpointRef);

        expect(statementRef.truncate(ElementKind.App)).toEqual(appRef);
        expect(() => statementRef.truncate(ElementKind.Type)).toThrow();
        expect(() => statementRef.truncate(ElementKind.Field)).toThrow();
        expect(statementRef.truncate(ElementKind.Endpoint)).toEqual(endpointRef);
        expect(statementRef.truncate(ElementKind.Statement)).toEqual(statementRef);
    });

    test.concurrent("push/pop app", () => {
        const rootApp = ElementRef.parse("App");
        expect(rootApp.pushApp("App2").toString()).toEqual("App::App2");
        expect(appRef.pushApp("App2").toString()).toEqual("Ns::App::App2");
        expect(typeRef.pushApp("App2").toString()).toEqual("Ns::App::App2.Type");
        expect(fieldRef.pushApp("App2").toString()).toEqual("Ns::App::App2.Type.Field");

        expect(() => rootApp.popApp().toString()).toThrow();
        expect(appRef.popApp().toString()).toEqual("Ns");
        expect(typeRef.popApp().toString()).toEqual("Ns.Type");
        expect(fieldRef.popApp().toString()).toEqual("Ns.Type.Field");
        expect(ElementRef.parse("Ns1::Ns2::App").popApp().toString()).toEqual("Ns1::Ns2");
    });

    test.concurrent("stringify", () => {
        const before = {
            num: 42,
            str: "Foo",
            ref: appRef,
            arr: [5, typeRef, "Bar", [fieldRef]],
            obj: { ref: appRef }
        };
        const after = {
            num: 42,
            str: "Foo",
            ref: "Ns::App",
            arr: [5, "Ns::App.Type", "Bar", ["Ns::App.Type.Field"]],
            obj: { ref: "Ns::App" }
        };
        expect(ElementRef.stringSubstitute(before)).toEqual(after);
        expect(before.ref instanceof ElementRef && before.arr[1] instanceof ElementRef).toBeTrue(); // No mutation
    });
});

describe("Computed properties", () => {
    test.concurrent("flags and kind", () => {
        expect(appRef)      .toEqual(expect.objectContaining({ isApp: true,  isType: false, isField: false, isEndpoint: false, isStatement: false, isData: true,  isBehavior: true,  kind: ElementKind.App,       numericKind: 1 }));
        expect(typeRef)     .toEqual(expect.objectContaining({ isApp: false, isType: true,  isField: false, isEndpoint: false, isStatement: false, isData: true,  isBehavior: false, kind: ElementKind.Type,      numericKind: 2 }));
        expect(endpointRef) .toEqual(expect.objectContaining({ isApp: false, isType: false, isField: false, isEndpoint: true,  isStatement: false, isData: false, isBehavior: true,  kind: ElementKind.Endpoint,  numericKind: 2 }));
        expect(fieldRef)    .toEqual(expect.objectContaining({ isApp: false, isType: false, isField: true,  isEndpoint: false, isStatement: false, isData: true,  isBehavior: false, kind: ElementKind.Field,     numericKind: 3 }));
        expect(statementRef).toEqual(expect.objectContaining({ isApp: false, isType: false, isField: false, isEndpoint: false, isStatement: true,  isData: false, isBehavior: true,  kind: ElementKind.Statement, numericKind: 3 }));
    });
});
