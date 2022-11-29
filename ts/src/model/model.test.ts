import "reflect-metadata";
import { readFile } from "fs/promises";
import "jest-extended";
import { realign } from "../common/format";
import { allItems } from "../common/iterate";
import { Annotation, Tag } from "./attribute";
import { Model } from "./model";
import "./renderers";
import { Action, Endpoint, Param, Statement } from "./statement";
import { Type } from "./type";
import { Application } from "./application";

const allPath = "../ts/test/all.sysl";
let allModel: Model;
let allSysl: string;

describe("Constructors", () => {
    test("New Model", () => {
        expect(new Model({})).not.toBeNull();
    });

    test("New Application", () => {
        var app = new Application({ name: "Foo" });
        expect(app).toHaveProperty("namespace", []);
        expect(app).toHaveProperty("name", "Foo");

        app = new Application({ namespace: ["Ns1", "Ns2"], name: "Foo" });
        expect(app).toHaveProperty("namespace", ["Ns1", "Ns2"]);
        expect(app).toHaveProperty("name", "Foo");
    });

    test("New Type", () => {
        expect(new Type("Foo")).toHaveProperty("name", "Foo");
    });

    test("New Endpoint", () => {
        expect(new Endpoint({ name: "Foo" })).toHaveProperty("name", "Foo");
    });

    test("New Param", () => {
        expect(new Param("foo", [])).toHaveProperty("name", "foo");
    });

    test("New Statement", () => {
        expect(new Statement({ value: new Action("foo") })).toHaveProperty("value.action", "foo");

        expect(Statement.action("foo")).toHaveProperty("value.action", "foo");
    });

    test("New Annotation", () => {
        expect(new Annotation({ name: "foo", value: "bar" })).toMatchObject({
            name: "foo",
            value: "bar",
        });
    });

    test("New Tag", () => {
        expect(new Tag({ name: "foo" })).toHaveProperty("name", "foo");
    });
});

describe("Serialization", () => {
    describe("Application", () => {
        test("empty", () => {
            expect(new Application({ name: "Foo" }).toSysl()).toEqual(
                realign(`
                    Foo:
                        ...`)
            );
        });
    });

    describe("Annotation", () => {
        test("escaped quotes", () => {
            const anno = new Annotation({ name: "foo", value: `"bar"` });
            expect(anno.toSysl()).toEqual(`foo = "\\"bar\\""`);
        });
    });
});

describe("Parent and Model", () => {
    test("all", async () => {
        const model = await Model.fromFile(allPath);
        expect(allItems(model).every(i => i.model === model)).toEqual(true);
    });
});

describe("Roundtrip", () => {
    // All
    test("AllRoundtrip", async () => {
        allModel = await Model.fromFile(allPath);
        allSysl = (await readFile(allPath)).toString();
        expect(allModel.filterByFile(allPath).toSysl()).toEqual(allSysl);
    });

    const cases = {
        EmptyApp: realign(
            `
            App:
                ...
            `
        ),
        EmptyAppWithSubpackages: realign(
            `
            App :: with :: subpackages:
                ...
            `
        ),
        AppWithTag: realign(
            `
            App [~abstract]:
                ...
            `
        ),
        InlineAnno: {
            input: realign(
                `
                App [name="value"]:
                    ...
                `
            ),
            output: realign(
                `
                App:
                    @name = "value"
                    ...
                `
            ),
        },
        StringAnnoInApp: realign(
            `
            App:
                @name = "value"
            `
        ),
        StringAnnoInType: realign(
            `
            App:
                !type Type:
                    @name = "value"
            `
        ),
        StringAnnoInField: realign(
            `
            App:
                !type Type:
                    Field <: int:
                        @name = "value"
            `
        ),
        StringAnnoEscaped: realign(
            `
            App:
                @name = "a \\"value\\""
                ...
            `
        ),
        MultilineAnno: realign(
            `
            App:
                @name =:
                    | anno
                    |  indented
                    |   across
                    |
                    |    multiple lines
                ...
            `
        ),
        ArrayAnno: realign(
            `
            App:
                @name = ["value1", "value2"]
                ...
            `
        ),
        NestedArrayAnno: realign(
            `
            App:
                @name = [["value1", "value2"], ["value3", "value4"]]
                ...
            `
        ),
        Endpoint: realign(
            `
            App:
                SimpleEp:
                    ...
            `
        ),
        EndpointWithTag: realign(
            `
            App:
                SimpleEp [~ignore]:
                    ...
            `
        ),
        EndpointWithAnno: realign(
            `
            App:
                SimpleEp:
                    @name = "value"
                    ...
            `
        ),
        EndpointWithInlineAnno: {
            input: realign(
                `
                App:
                    SimpleEp [name="value"]:
                        ...
                `
            ),
            output: realign(
                `
                App:
                    SimpleEp:
                        @name = "value"
                        ...
                `
            ),
        },
        EndpointWithUntypedParam: realign(
            `
            App:
                SimpleEp (foo):
                    ...
            `
        ),
        EndpointWithNamedPrimitiveParam: realign(
            `
            App:
                SimpleEp (param <: string):
                    ...
            `
        ),
        EndpointWithUnnamedRefParam: realign(
            `
            App:
                SimpleEp (Types.type):
                    ...
            `
        ),
        EndpointWithNamedRefParam: realign(
            `
            App:
                SimpleEp (param <: Types.type):
                    ...
            `
        ),
        EndpointWithPrimitiveParamWithConstraints: realign(
            `
            App:
                SimpleEp (unlimited <: string(5..), limited <: string(5..10), num <: int(5)):
                    ...
            `
        ),
        EndpointWithCall: realign(
            `
            App:
                SimpleEp:
                    App2 <- Endpoint
            `
        ),
        EndpointWithPrimitiveReturn: realign(
            `
            App:
                SimpleEp:
                    return ok <: string
            `
        ),
        EndpointWithRefReturn: realign(
            `
            App:
                SimpleEp:
                    return ok <: Types.type
            `
        ),
        RestEndpoint: realign(
            `
            RestEndpoint:
                /:
                    GET:
                        ...
            `
        ),
        RestEndpointWithoutNesting: realign(
            `
            RestEndpoint:
                /nested/path:
                    GET:
                        ...
            `
        ),
        RestEndpointWithNesting: {
            input: realign(
                `
            RestEndpoint:
                /nested:
                    /path:
                        GET:
                            ...
            `
            ),
            output: realign(
                `
            RestEndpoint:
                /nested/path:
                    GET:
                        ...
            `
            ),
        },
        RestEndpointWithTypeInPath: realign(
            `
            RestEndpoint:
                /pathwithtype/{native <: int}:
                    GET:
                        ...
            `
        ),
        RestEndpointWithQueryParams: realign(
            `
            RestEndpoint:
                /query:
                    GET?native=string&optional=string?:
                        ...
            `
        ),
        RestEndpointWithRefParam: realign(
            `
            RestEndpoint:
                /param:
                    PATCH (t <: Types.Type [~body]):
                        ...
            `
        ),
        RestEndpointWithPrimitiveParam: realign(
            `
            RestEndpoint:
                /param:
                    POST (native <: string):
                        ...
            `
        ),
        RestEndpointWithConstrainedParams: realign(
            `
            RestEndpoint:
                /param:
                    PUT (unlimited <: string(5..), limited <: string(5..10), num <: int(5)):
                        ...
            `
        ),
        Type: realign(
            `
            App:
                !type Type:
                    @annotation = "annotation"
                    nativeTypeField <: string
                    reference <: RestEndpoint.Type
                    optional <: string?
                    set <: set of string
                    sequence <: sequence of string
                    aliasSequence <: AliasSequence
                    with_anno <: string:
                        @annotation = "this is an annotation"
            `
        ),
        Table: realign(
            `
            App:
                !table Table [~tag]:
                    primaryKey <: string [~pk]
                    nativeTypeField <: string
                    reference <: RestEndpoint.Type
                    optional <: string?
                    set <: set of string
                    sequence <: sequence of string
                    with_anno <: string:
                        @annotation = "this is an annotation"
                    decimal_with_precision <: decimal(5.8)
                    string_max_constraint <: string(5)
                    string_range_constraint <: string(5..10)
            `
        ),
        Enum: realign(
            `
            App:
                !enum Enum [~tag]:
                    ENUM_1: 1
                    ENUM_2: 2
                    ENUM_3: 3
            `
        ),
        Alias: realign(
            `
            App:
                !alias A:
                    int
            `
        ),
        TypeRef: realign(
            `
            Namespace :: App:
                !type Type:
                    shortRef <: Type
                    fullRef <: Namespace::App.Type
            `
        ),
        UnsafeNames: realign(
            `
            %28App%29Name%21:
                !type %28Type%29Name%21:
                    %28Field%29Name%21 <: %28App%29Name%21.%28Type%29Name%21 [~%28Tag%29Name%21]
            `
        ),
        // Lists are not well supported, so we substitute them when serializing to source.
        List: {
            input: realign(
                `
                App:
                    !type Type:
                        list(1..1) <: string
                `
            ),
            output: realign(
                `
                App:
                    !type Type:
                        list <: sequence of string
                `
            ),
        },
        PreservesOrder: realign(`
            App2 [~tag2, ~tag1]:
                @anno2 = "2"
                @anno1 = "1"
                !table Table2 [~tag2, ~tag1]:
                    @anno2 = "2"
                    @anno1 = "1"
                    Field2 <: int [~tag2, ~tag1]:
                        @anno2 = "2"
                        @anno1 = "1"
                    Field1 <: int
            
                !table Table1:
                    ...

            App1:
                ...
            `),
    };

    type SyslCase = { input: string; output: string };
    type TestSysl = SyslCase | string;

    // sysl should be of type TestSysl, but the compiler treats `SyslCase | string` as `string`.
    const inputSysl = (sysl: TestSysl): string => (typeof sysl == "string" ? sysl : sysl.input);
    const expectedSysl = (sysl: TestSysl): string => (typeof sysl == "string" ? sysl : sysl.output);

    test.each(Object.entries(cases))("%s", async (_, sysl: TestSysl) => {
        const model = await Model.fromText(inputSysl(sysl as SyslCase));
        expect(model.toSysl()).toEqual(expectedSysl(sysl as SyslCase));
    });
});
