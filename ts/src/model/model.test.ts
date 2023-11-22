import "jest-extended";
import { readFile } from "fs/promises";
import { realign } from "../common/format";
import { allItems } from "../common/iterate";
import { Annotation, AnnoValue, Tag } from "./attribute";
import { Model } from "./model";
import { ActionStatement, GroupStatement, ParentStatement, Statement } from "./statement";
import { Endpoint } from "./endpoint";
import { Type } from "./type";
import { Application } from "./application";
import { Field, FieldValue } from "./field";
import { ILocational } from "./common";
import { ElementRef } from "./elementRef";
import { CloneContext, ICloneable, ModelFilters } from "./clone";
import { Union } from "./union";
import { Primitive, TypePrimitive } from "./primitive";
import { Alias } from "./alias";
import { CollectionDecorator } from "./decorator";
import { Element } from "./element";
import { Enum } from "./enum";

const allPath = "../ts/test/all.sysl";
let allModel: Model;
let allSysl: string;

describe("Constructors", () => {
    test.concurrent("New Model", () => {
        expect(new Model({})).not.toBeNull();
    });

    test.concurrent("New Application", () => {
        var app = new Application("Foo");
        expect(app).toHaveProperty("namespace", []);
        expect(app).toHaveProperty("name", "Foo");

        app = new Application("Ns1::Ns2::Foo");
        expect(app).toHaveProperty("namespace", ["Ns1", "Ns2"]);
        expect(app).toHaveProperty("name", "Foo");

        app = new Application("Foo", { namespace: ["Ns1", "Ns2"] });
        expect(app).toHaveProperty("namespace", ["Ns1", "Ns2"]);
        expect(app).toHaveProperty("name", "Foo");

        expect(() => new Application(ElementRef.parse("Foo"), { namespace: ["Ns1", "Ns2"] })).toThrow();
        expect(() => new Application(ElementRef.parse("Ns1::Ns2::Foo.Type"))).toThrow();
        expect(() => new Application("Ns1::Ns2::Foo", { namespace: ["Ns1", "Ns2"] })).toThrow();
        expect(() => new Application("Ns1::Ns2::Foo.Type")).toThrow();
        expect(() => new Application("Ns1::Ns2::Foo.Type.Field.Crazy")).toThrow();
    });

    test.concurrent("New Type", () => {
        expect(new Type("Foo")).toHaveProperty("name", "Foo");
    });

    test.concurrent("New Endpoint", () => {
        expect(new Endpoint("Foo")).toHaveProperty("name", "Foo");
    });

    test.concurrent("New Statement", () => {
        expect(new ActionStatement("foo")).toHaveProperty("action", "foo");

        expect(new ActionStatement("foo")).toHaveProperty("action", "foo");
    });

    test.concurrent("New Annotation", () => {
        expect(new Annotation("foo", "bar")).toMatchObject({ name: "foo", value: "bar" });
        expect(new Annotation("foo", "C:\\\\bar")).toMatchObject({ name: "foo", value: "C:\\\\bar" });
    });

    test.concurrent("New Tag", () => {
        expect(new Tag("foo")).toHaveProperty("name", "foo");
    });
});

describe("Sysl rendering", () => {
    describe("Application", () => {
        test.concurrent("empty", () => {
            expect(new Application("Foo").toSysl()).toEqual(realign(`
                    Foo:
                        ...`) //TODO: Fix missing newline at end of file
            );
        });

        test.concurrent("reverse", async () => {
            const model = await Model.fromText(realign(`
                Foo:
                    ...
            `));

            expect(model.apps).toHaveLength(1);
            expect(model.apps[0].name).toEqual("Foo");
            expect(model.apps[0].namespace).toBeEmpty();
            expect(model.apps[0].children).toBeEmpty();
            expect(model.apps[0].endpoints).toBeEmpty();
            expect(model.apps[0].annos).toBeEmpty();
            expect(model.apps[0].tags).toBeEmpty();
        });
    });

    describe("Field", () => {
        test.concurrent("Reference types", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        Abs <: Company::App.Address
                        Rel <: Address
                        
                        SeqAbs <: sequence of Company::App.Address
                        SeqRel <: sequence of Address
                        SetAbs <: set of Company::App.Address
                        SetRel <: set of Address
                        
                        ListAbs(1..1) <: Company::App.Address
                        ListRel(1..1) <: Address

                        ListSeqAbs(1..1) <: sequence of Company::App.Address
                        ListSeqRel(1..1) <: sequence of Address
                        ListSetAbs(1..1) <: set of Company::App.Address
                        ListSetRel(1..1) <: set of Address
                    !type Address:
                        City <: string
            `));

            const addressRef = model.getType("Company::App.Address").toRef();
            const sequenceRef = new CollectionDecorator(addressRef, false);
            const setRef = new CollectionDecorator(addressRef, true);
            
            model.getType("Company::App.Customer").children.forEach(f => {
                let expected: FieldValue = addressRef;
                if (f.name.includes("Set")) expected = setRef;
                else if (f.name.includes("Seq") || f.name.includes("List")) expected = sequenceRef;
                expect(f.value).toEqual(expected);
            });
        });

        test.concurrent("bigger example", async () => {
            const model = await Model.fromText(realign(`
                App:
                    !type ResponseData:
                        group <: string
                        list(1..1) <: sequence of DataWrapper
                    !type DataWrapper:
                        data <: Data?
                    !type Data:
                        merchants <: sequence of Customer?
                        info <: string
                        responseMsg <: string
                    !type Customer:
                        name <: string
            `));
            const fieldType = model.getField("App.ResponseData.list").value as CollectionDecorator;
            expect(fieldType).toEqual(new CollectionDecorator(model.getType("App.DataWrapper").toRef(), false));
        });

        test.concurrent("escaped backslash", () => {
            const anno = new Annotation("proto_options", ["key = Foo\\\\Bar"]);
            expect(anno.toSysl()).toEqual(`proto_options = ["key = Foo\\\\Bar"]`);
        });
    });

    test.concurrent("Special value handling", () => {
        expect(new Annotation("foo", `"bar"`).toSysl()).toEqual(`foo = "\\"bar\\""`);
        expect(new Annotation("foo", "bar\\baz").toSysl()).toEqual(`foo = "bar\\baz"`);
        expect(new Application(new ElementRef([], "App .")).toSysl()).toEqual("App%20%2E:\n    ...");
        expect(new Application(new ElementRef([], "int")).toSysl()).toEqual("%69nt:\n    ...");
        expect(new Application(new ElementRef([], "1App")).toSysl()).toEqual("%31App:\n    ...");
    });
});

describe("DTO", () => {
    test.concurrent("Data elements", async () => {
        const model = await Model.fromText(realign(`
            # Header

            import other.sysl

            Ns :: App [~appTag]:
                @appAnno = "App"
                !type Type [~typeTag]:
                    @typeAnno = "Type"
                    Primitive <: string(5..10)
                    Array <: sequence of decimal(5.8)
                    Reference <: Other.Type [~fieldTag]:
                        @fieldAnno = "Field"
        `), "main.sysl", { maxImportDepth: 1 });

        expect(model.toDto()).toMatchObject(
            {
                header: "# Header",
                imports: [ { filePath: "other.sysl", locations: ["ts/main.sysl:3:1::18"] } ],
                apps: [
                    {
                        kind: "Application",
                        namespace: ["Ns"],
                        name: "App",
                        metadata: { appTag: undefined, appAnno: "App" },
                        locations: {
                            0: "ts/main.sysl:5:1:12:33",
                            appTag: "ts/main.sysl:5:12::19",
                            appAnno: "ts/main.sysl:6:16::21",
                        },
                        children: [
                            {
                                kind: "Type",
                                name: "Type",
                                metadata: { typeTag: undefined, typeAnno: "Type" },
                                locations: {
                                    0: "ts/main.sysl:7:5:12:33",
                                    typeTag: "ts/main.sysl:7:17::25",
                                    typeAnno: "ts/main.sysl:8:21::27",
                                },
                                children: [
                                    {
                                        kind: "Field",
                                        name: "Primitive",
                                        primitive: "string",
                                        constraint: "(5..10)",
                                        locations: { 0: "ts/main.sysl:9:9::35" },
                                    },
                                    {
                                        kind: "Field",
                                        name: "Array",
                                        primitive: "decimal",
                                        constraint: "(5.8)",
                                        collectionType: "sequence",
                                        locations: { 0: "ts/main.sysl:10:9::42" },
                                    },
                                    {
                                        kind: "Field",
                                        name: "Reference",
                                        ref: "Other.Type",
                                        metadata: { fieldTag: undefined, fieldAnno: "Field" },
                                        locations: {
                                            0: "ts/main.sysl:11:9:13:2",
                                            fieldTag: "ts/main.sysl:11:34::43",
                                            fieldAnno: "ts/main.sysl:12:26::33",
                                        }
                                    },
                                ],
                            }
                        ],
                    }
                ],
            }
        );
    });

    test.concurrent("Rest endpoint", async () => {
        const model = await Model.fromText(realign(`
            Ns :: App:
                /customer/{id <: int}:
                    GET(auth_token <: string(32..64) [~paramTag, paramAnno=[["1"],"2"]]) ?includeOrders=bool? [~restTag]:
                        @restAnno = "GET Customer"
                        if includeOrders == true:
                            return data <: CustomerAndOrders
                        else:
                            throw NotImplementedError
        `), "main.sysl", { maxImportDepth: 1 });

        expect(model.toDto()).toMatchObject(
            {
                apps: [
                    {
                        children: [
                            {
                                name: "GET /customer/{id}",
                                params: [
                                    {
                                        kind: "Field",
                                        name: "auth_token",
                                        optional: false,
                                        primitive: "string",
                                        constraint: "(32..64)",
                                        metadata: {
                                            paramTag: undefined,
                                            paramAnno: [["1"],"2"],
                                        },
                                        locations: {
                                            "paramTag": "ts/main.sysl:3:43::52",
                                            "paramAnno": "ts/main.sysl:3:64::75",
                                        },
                                    }
                                ],
                                restParams: {
                                    method: "GET",
                                    path: "/customer/{id}",
                                    queryParams: [
                                        {
                                            kind: "Field",
                                            name: "includeOrders",
                                            optional: true,
                                            primitive: "bool",
                                            locations: { 0: "ts/main.sysl:3:79::98" },
                                        }
                                    ],
                                    urlParams: [
                                        {
                                            kind: "Field",
                                            name: "id",
                                            optional: false,
                                            primitive: "int",
                                            locations: { 0: "ts/main.sysl:2:15::26" },
                                        }
                                    ],
                                },
                                metadata: { restTag: undefined, restAnno: "GET Customer" },
                                locations: {
                                    0: "ts/main.sysl:3:9:9:2",
                                    restTag: "ts/main.sysl:3:100::108",
                                    restAnno: "ts/main.sysl:4:25::39", // TODO: https://github.com/anzx/sysl/issues/888
                                },
                                children: [
                                    {
                                        kind: "CondStatement",
                                        prefix: "if",
                                        title: "includeOrders == true",
                                        children: [
                                            {
                                                kind: "ReturnStatement",
                                                payload: "data <: CustomerAndOrders",
                                                locations: { 0: "ts/main.sysl:6:17::49" },
                                            }
                                        ],
                                        locations: { 0: "ts/main.sysl:5:13:7:14" },
                                    },
                                    {
                                        kind: "GroupStatement",
                                        prefix: "",
                                        title: "else",
                                        children: [
                                            {
                                                kind: "ActionStatement",
                                                action: "throw NotImplementedError",
                                                locations: { 0: "ts/main.sysl:8:17::42" },
                                            }
                                        ],
                                        locations: { 0: "ts/main.sysl:7:13:8:42" },
                                    },
                                ],
                            }
                        ],
                    }
                ],
            }
        );
    });

    test.only("RPC endpoint", async () => {
        const model = await Model.fromText(realign(`
            Ns :: App:
                GetCustomer(id <: int, auth_token <: string(32..64) [~paramTag, paramAnno=[["1"],"2"]], includeOrders <: bool?) [~rpcTag]:
                    @rpcAnno = "GetCustomer"
                    if includeOrders == true:
                        return data <: CustomerAndOrders
                    else:
                        throw NotImplementedError
        `), "main.sysl", { maxImportDepth: 1 });

        expect(model.toDto()).toMatchObject(
            {
                apps: [
                    {
                        children: [
                            {
                                name: "GetCustomer",
                                params: [
                                    {
                                        kind: "Field",
                                        name: "id",
                                        optional: false,
                                        primitive: "int",
                                        locations: { }, // TODO: https://github.com/anzx/sysl/issues/891
                                    },
                                    {
                                        kind: "Field",
                                        name: "auth_token",
                                        optional: false,
                                        primitive: "string",
                                        constraint: "(32..64)",
                                        metadata: {
                                            paramTag: undefined,
                                            paramAnno: [["1"],"2"],
                                        },
                                        locations: {
                                            "paramAnno": "ts/main.sysl:2:69::90",
                                            "paramTag": "ts/main.sysl:2:58::67",
                                        },
                                    },
                                    {
                                        kind: "Field",
                                        name: "includeOrders",
                                        optional: true,
                                        primitive: "bool",
                                        locations: { },
                                    }

                                ],
                                metadata: { rpcTag: undefined, rpcAnno: "GetCustomer" },
                                locations: {
                                    "0": "ts/main.sysl:2:5:7:38",
                                    "rpcAnno": "ts/main.sysl:3:9::33",
                                    "rpcTag": "ts/main.sysl:2:118::125",
                                },
                                children: [
                                    {
                                        kind: "CondStatement",
                                        prefix: "if",
                                        title: "includeOrders == true",
                                        children: [
                                            {
                                                kind: "ReturnStatement",
                                                payload: "data <: CustomerAndOrders",
                                                locations: { 0: "ts/main.sysl:5:13::45" },
                                            }
                                        ],
                                        locations: { 0: "ts/main.sysl:4:9:6:10" },
                                    },
                                    {
                                        kind: "GroupStatement",
                                        prefix: "",
                                        title: "else",
                                        children: [
                                            {
                                                kind: "ActionStatement",
                                                action: "throw NotImplementedError",
                                                locations: { 0: "ts/main.sysl:7:13::38" },
                                            }
                                        ],
                                        locations: { 0: "ts/main.sysl:6:9:7:38" },
                                    },
                                ],
                            }
                        ],
                    }
                ],
            }
        );
    });

    test.concurrent("All", async () => (await Model.fromFile(allPath)).toDto());
});

describe("Parent and Model", () => {
    test.concurrent("all", async () => {
        const model = await Model.fromFile(allPath);
        expect(allItems(model).every(i => i.model === model)).toEqual(true);
    });

    test.concurrent("move endpoint, attach sub-statement", async () => {
        const model = await Model.fromText(realign(`
            App1:
                Endpoint:
                    statement:
                        ...
        `));
        model.apps.push(new Application("App2"));
        const app1 = model.getApp("App1");
        const app2 = model.getApp("App2");
        const ep = app1.children.pop() as Endpoint;
        const outerSt = ep.children[0] as GroupStatement;
        const innerSt = new GroupStatement("subStatement");
        outerSt.children.push(innerSt);
        app2.children.push(ep);
        model.attachSubitems();

        expect(app1.endpoints).toBeEmpty();
        expect(app2.model === model).toBeTrue();
        expect(app2.parent).toBeUndefined();
        expect(ep.model === model).toBeTrue();
        expect(ep.parent === app2).toBeTrue();
        expect(outerSt.model === model).toBeTrue();
        expect(outerSt.parent === ep).toBeTrue();
        expect(innerSt.model === model).toBeTrue();
        expect(innerSt.parent === outerSt).toBeTrue();
    });
});

describe("Roundtrip", () => {
    test.concurrent("AllRoundtrip", async () => {
        allModel = await Model.fromFile(allPath);
        allSysl = (await readFile(allPath)).toString();
        expect(allModel.filterByFile(allPath).toSysl()).toEqual(allSysl);
    });

    const cases = {
        EmptyApp: realign(`
            App:
                ...
        `),
        EmptyAppWithSubpackages: realign(`
            App :: with :: subpackages:
                ...
        `),
        AppWithTag: realign(`
            App [~abstract]:
                ...
            `),
        InlineAnno: {
            input: realign(`
                App [name="value"]:
                    ...
                `),
            output: realign(`
                App:
                    @name = "value"
                `),
        },
        StringAnnoInApp: realign(
            `
            App:
                @name = "value"
            `
        ),
        StringAnnoInType: realign(`
            App:
                !type Type:
                    @name = "value"
            `),
        StringAnnoInField: realign(`
            App:
                !type Type:
                    Field <: int:
                        @name = "value"
            `),
        StringAnnoEscaped: realign(`
            App:
                @name = "hello;,?:@&=+$.!~*'()\\"/\\\\ world"
            `),
        MultilineAnno: realign(`
            App:
                @name =:
                    | anno
                    |  indented
                    |   across
                    |
                    |    multiple lines
                    |     with special chars: hello;,?:@&=+$.!~*'()"/\\ world
            `),
        ArrayAnno: realign(`
            App:
                @name = ["value1", "value2", "value3 = C:\\\\value\\\\4"]
            `),
        NestedArrayAnno: realign(`
            App:
                @name = [["value1", "value2"], ["value3", "value4"]]
            `),
        Endpoint: realign(`
            App:
                SimpleEp:
                    ...
            `),
        EndpointWithTag: realign(`
            App:
                SimpleEp [~ignore]:
                    ...
            `),
        EndpointWithAnno: realign(`
            App:
                SimpleEp:
                    @name = "value"
            `),
        EndpointWithInlineAnno: {
            input: realign(`
                App:
                    SimpleEp [name="value"]:
                        ...
                `),
            output: realign(`
                App:
                    SimpleEp:
                        @name = "value"
                `),
        },
        EndpointWithUntypedParam: realign(`
            App:
                SimpleEp (foo):
                    ...
            `),
        EndpointWithNamedPrimitiveParam: realign(`
            App:
                SimpleEp (param <: string):
                    ...
            `),
        EndpointWithNamedRefParam: realign(`
            App:
                SimpleEp (param <: Types.type):
                    ...
            `),
        EndpointWithLocalNamedRefParam: realign(`
            App:
                SimpleEp (param <: type):
                    ...
            `),
        EndpointWithPrimitiveParamWithConstraints: realign(`
            App:
                SimpleEp (unlimited <: string(5..), limited <: string(5..10), num <: int(5)):
                    ...
            `),
        EndpointWithCall: realign(`
            App:
                SimpleEp:
                    App2 <- Endpoint
            `),
        EndpointWithPrimitiveReturn: realign(`
            App:
                SimpleEp:
                    return ok <: string
            `),
        EndpointWithRefReturn: realign(`
            App:
                SimpleEp:
                    return ok <: Types.type
            `),
        RestEndpoint: realign(`
            RestEndpoint:
                /:
                    GET:
                        ...
            `),
        RestEndpointWithoutNesting: realign(`
            RestEndpoint:
                /nested/path:
                    GET:
                        ...
            `),
        RestEndpointWithNesting: {
            input: realign(`
                RestEndpoint:
                    /nested:
                        /path:
                            GET:
                                ...
            `),
            output: realign(`
                RestEndpoint:
                    /nested/path:
                        GET:
                            ...
                `),
        },
        RestEndpointWithTypeInPath: realign(`
            RestEndpoint:
                /pathwithtype/{native <: int}:
                    GET:
                        ...
            `),
        RestEndpointWithQueryParams: realign(`
            RestEndpoint:
                /query:
                    GET?native=string&optional=string?:
                        ...
            `),
        RestEndpointWithRefParam: realign(`
            RestEndpoint:
                /param:
                    PATCH (t <: Types.Type [~body]):
                        ...
            `),
        RestEndpointWithPrimitiveParam: realign(`
            RestEndpoint:
                /param:
                    POST (native <: string):
                        ...
            `),
        RestEndpointWithConstrainedParams: realign(`
            RestEndpoint:
                /param:
                    PUT (unlimited <: string(5..), limited <: string(5..10), num <: int(5)):
                        ...
            `),
        RestEndpointWithComplexParamType: realign(`
            RestEndpoint:
                /param:
                    POST (arg <: Customer):
                        ...
            `),
        Type: realign(`
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
            `),
        Table: realign(`
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
            `),
        Enum: realign(`
            App:
                !enum Enum [~tag]:
                    ENUM_1: 1
                    ENUM_2: 2
                    ENUM_3: 3
            `),
        Alias: realign(`
            App:
                !alias A:
                    int
            `),
        TypeRef: {
            input: realign(`
                Namespace :: App:
                    !type Type:
                        shortRef <: Type
                        fullRef <: Namespace::App.Type
                        extRef <: Namespace::External.Type
                
                Namespace :: External:
                    !type Type:
                        ...
            `),
            output: realign(`
                Namespace :: App:
                    !type Type:
                        shortRef <: Type
                        fullRef <: Type
                        extRef <: Namespace::External.Type
                
                Namespace :: External:
                    !type Type:
                        ...
                `),
        },
        UnsafeNames: realign(`
            %28App%29Name%21:
                !type %28Type%29Name%21:
                    %28Field%29Name%21 <: %28Namespace%21::%28App%29Name2%21.%28Type%29Name%21 [~%28Tag%29Name%21]

            %28Namespace%21 :: %28App%29Name2%21:
                !type %28Type%29Name%21:
                    ...
            
            %69nt :: %73tring :: %64ate:
                !type %62ool:
                    %61ny <: any

            %31Ns :: %32Ns :: %33App:
                !type %34Type:
                    %35Field <: int
            `),
        // Lists are not well supported, so we substitute them when serializing to source.
        List: {
            input: realign(`
                App:
                    !type Type:
                        list(1..1) <: string
                `),
            output: realign(`
                App:
                    !type Type:
                        list <: sequence of string
                `),
        },
        ListSequence: {
            input: realign(`
                App:
                    !type Type:
                        list(1..1) <: sequence of string
                `),
            output: realign(`
                App:
                    !type Type:
                        list <: sequence of string
                `),
        },
        ListSet: {
            input: realign(`
                App:
                    !type Type:
                        list(1..1) <: set of string
                `),
            output: realign(`
                App:
                    !type Type:
                        list <: set of string
                `),
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
        // TODO: This is how the pb is generated from the Sysl binary, investigate further
        InlineComplexType: {
            input: realign(`
                App:
                    !type Type:
                        Field <:
                            Subfield <: int
                `),
            output: realign(`
                App:
                    !type Type:
                        Field <: Field

                    !type Type%2EField:
                        Subfield <: int
                `)
        },
    };

    type SyslCase = { input: string; output: string };
    type TestSysl = SyslCase | string;

    // sysl should be of type TestSysl, but the compiler treats `SyslCase | string` as `string`.
    const inputSysl = (sysl: TestSysl): string => (typeof sysl == "string" ? sysl : sysl.input);
    const expectedSysl = (sysl: TestSysl): string => (typeof sysl == "string" ? sysl : sysl.output);

    test.concurrent.each(Object.entries(cases))("%s", async (_, sysl: TestSysl) => {
        const model = await Model.fromText(inputSysl(sysl as SyslCase));
        expect(model.toSysl()).toEqual(expectedSysl(sysl as SyslCase));
    });
});

describe("General methods", () => {
    test.concurrent("toSyslPath", () => {
        const model = new Model({ syslRoot: "/usr/MyShop" });
        // @ts-ignore
        expect(model.convertSyslPath("/usr/MyShop/backend.sysl")).toEqual("backend.sysl");
        // @ts-ignore
        expect(model.convertSyslPath("/usr/MyShop/schema/backend.sysl")).toEqual("schema/backend.sysl");
        // @ts-ignore
        expect(model.convertSyslPath("schema/backend.sysl", "/usr/MyShop")).toEqual("schema/backend.sysl");
        // @ts-ignore
        expect(model.convertSyslPath("backend.sysl", "/usr/MyShop/schema")).toEqual("schema/backend.sysl");
        // @ts-ignore
        expect(model.convertSyslPath("../index.sysl", "/usr/MyShop/schema")).toEqual("index.sysl");
        // @ts-ignore
        expect(() => model.convertSyslPath("../../outside.sysl", "/usr/MyShop/schema")).toThrowError(
            "is outside the Sysl root path"
        );
        // @ts-ignore
        expect(() => model.convertSyslPath("outside.sysl", "/usr/")).toThrowError("is outside the Sysl root path");
    });

    test.concurrent("find/get methods", async () => {
        const model = await Model.fromText(realign(`
            Namespace :: App:
                !type Type:
                    Field <: int
                Endpoint:
                    Statement0:
                        Statement0_0
                        Statement0_1:
                            Statement0_1_0
                        Statement0_2

        `));

        // Element not in model
        expect(model.findElement("Namespace::MissingApp")).toBeUndefined();
        expect(model.findElement("Namespace::App.MissingType")).toBeUndefined();
        expect(model.findElement("Namespace::App.[MissingEndpoint]")).toBeUndefined();
        expect(model.findElement("Namespace::App.Type.MissingField")).toBeUndefined();
        expect(model.findElement("Namespace::App.[Endpoint].[1]")).toBeUndefined();
        expect(model.findElement("Namespace::App.[Endpoint].[0,3]")).toBeUndefined();
        expect(model.findElement("Namespace::App.[Endpoint].[0,0,0]")).toBeUndefined();
        expect(model.findApp("Namespace::MissingApp")).toBeUndefined();
        expect(model.findType("Namespace::App.MissingType")).toBeUndefined();
        expect(model.findEndpoint("Namespace::App.[MissingEndpoint]")).toBeUndefined();
        expect(model.findField("Namespace::App.Type.MissingField")).toBeUndefined();
        expect(model.findStatement("Namespace::App.[Endpoint].[1]")).toBeUndefined();
        expect(model.findStatement("Namespace::App.[Endpoint].[0,3]")).toBeUndefined();
        expect(model.findStatement("Namespace::App.[Endpoint].[0,0,0]")).toBeUndefined();
        expect(() => model.getApp("Namespace::MissingApp")).toThrow();
        expect(() => model.getType("Namespace::App.MissingType")).toThrow();
        expect(() => model.getEndpoint("Namespace::App.[MissingEndpoint]")).toThrow();
        expect(() => model.getField("Namespace::App.Type.MissingField")).toThrow();
        expect(() => model.getStatement("Namespace::App.[Endpoint].[1]")).toThrow();
        expect(() => model.getStatement("Namespace::App.[Endpoint].[0,3]")).toThrow();
        expect(() => model.getStatement("Namespace::App.[Endpoint].[0,0,0]")).toThrow();

        // Invalid element reference
        expect(() => model.findElement("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findApp("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findType("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findEndpoint("Namespace::App.[Endpoint].[]")).toThrow();
        expect(() => model.findField("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findStatement("Namespace::App.[Endpoint].[]")).toThrow();
        expect(() => model.getApp("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.getType("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.getEndpoint("Namespace::App.[Endpoint].[]")).toThrow();
        expect(() => model.getField("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.getStatement("Namespace::App.[Endpoint].[]")).toThrow();

        // Element type mismatch
        expect(() => model.findApp("Namespace::App.Type")).toThrow();
        expect(() => model.findType("Namespace::App.Type.Field")).toThrow();
        expect(() => model.findEndpoint("Namespace::App.[Endpoint].[0]")).toThrow();
        expect(() => model.findField("Namespace::App")).toThrow();
        expect(() => model.findStatement("Namespace::App.[Endpoint]")).toThrow();
        expect(() => model.getApp("Namespace::App.Type")).toThrow();
        expect(() => model.getType("Namespace::App.Type.Field")).toThrow();
        expect(() => model.getEndpoint("Namespace::App.[Endpoint].[0]")).toThrow();
        expect(() => model.getField("Namespace::App")).toThrow();
        expect(() => model.getStatement("Namespace::App.[Endpoint]")).toThrow();

        // Happy path
        expect(model.findElement("Namespace::App")).toBeInstanceOf(Application);
        expect(model.findApp("Namespace::App")).toBeInstanceOf(Application);
        expect(model.getApp("Namespace::App")).toBeInstanceOf(Application);

        expect(model.findElement("Namespace::App.Type")).toBeInstanceOf(Type);
        expect(model.findType("Namespace::App.Type")).toBeInstanceOf(Type);
        expect(model.getType("Namespace::App.Type")).toBeInstanceOf(Type);

        expect(model.findElement("Namespace::App.[Endpoint]")).toBeInstanceOf(Endpoint);
        expect(model.findEndpoint("Namespace::App.[Endpoint]")).toBeInstanceOf(Endpoint);
        expect(model.getEndpoint("Namespace::App.[Endpoint]")).toBeInstanceOf(Endpoint);

        expect(model.findElement("Namespace::App.Type.Field")).toBeInstanceOf(Field);
        expect(model.findField("Namespace::App.Type.Field")).toBeInstanceOf(Field);
        expect(model.getField("Namespace::App.Type.Field")).toBeInstanceOf(Field);

        expect(model.findElement("Namespace::App.[Endpoint].[0]")).toBeInstanceOf(Statement);
        expect(model.findStatement("Namespace::App.[Endpoint].[0]")).toBeInstanceOf(Statement);
        expect(model.getStatement("Namespace::App.[Endpoint].[0]")).toBeInstanceOf(Statement);

        const statement = model.apps[0].endpoints[0].children[0] as ParentStatement;
        expect(model.getStatement("Namespace::App.[Endpoint].[0]")).toBe(statement);
        expect(model.getStatement("Namespace::App.[Endpoint].[0,1]")).toBe(statement.children[1]);
        expect(model.getStatement("Namespace::App.[Endpoint].[0,1,0]"))
            .toBe((statement.children[1] as ParentStatement).children[0]);
    });
});

describe("Cloning", () => {
    test.concurrent("Annotation", async () => {
        const model = await Model.fromText(realign(`
            App1:
                @anno1 = "value1"
                @annoArr1 = ["v1", ["v2"]]

            App2:
                ...
        `));

        const app1 = model.getApp("App1");

        const clonedAnno = app1.annos[0].clone();
        clonedAnno.name = "anno2";                           // Modify name and value to ensure originals don't change
        clonedAnno.value = "value2";                         //

        const clonedAnnoArr = app1.annos[1].clone();
        clonedAnnoArr.name = "annoArr2";                     //
        (clonedAnnoArr.value as AnnoValue[])[0] = "v3";      // Modify name and value to ensure originals don't change
        (clonedAnnoArr.value as AnnoValue[][])[1][0] = "v4"; //

        model.getApp("App2").annos.push(clonedAnno, clonedAnnoArr);

        expect(model.toSysl()).toEqual(realign(`
            App1:
                @anno1 = "value1"
                @annoArr1 = ["v1", ["v2"]]

            App2:
                @anno2 = "value2"
                @annoArr2 = ["v3", ["v4"]]
        `));
    });

    test.concurrent("Tag", async () => {
        const model = await Model.fromText(realign(`
            App1 [~tag1]:
                ...

            App2:
                ...
        `));

        const clonedTag = model.getApp("App1")!.tags[0].clone();
        clonedTag.name = "tag2";                            // Modify name to ensure originals don't change

        model.getApp("App2").tags.push(clonedTag);

        expect(model.toSysl()).toEqual(realign(`
            App1 [~tag1]:
                ...

            App2 [~tag2]:
                ...
        `));
    });

    test.concurrent("Field", async () => {
        const model = await Model.fromText(realign(`
            App1:
                !type Type1:
                    Field1 <: int [~tag1]:
                        @anno = "value"

            App2:
                !type Type2:
                    ...
        `));

        const clonedField = model.getField("App1.Type1.Field1").clone();
        clonedField.name = "Field2";                        //
        clonedField.tags[0].name = "tag2";                  //
        clonedField.annos[0].name = "anno2";                // Modify name to ensure originals don't change
        clonedField.annos[0].value = "value2";              //

        model.getType("App2.Type2").children.push(clonedField);

        expect(model.toSysl()).toEqual(realign(`
            App1:
                !type Type1:
                    Field1 <: int [~tag1]:
                        @anno = "value"

            App2:
                !type Type2:
                    Field2 <: int [~tag2]:
                        @anno2 = "value2"
        `));
    });

    test.concurrent("Field, filtered", async () => {
        const model = await Model.fromText(realign(`
            App1:
                !type Type1:
                    Field1 <: int [~tag1]:
                        @anno = "value"

            App2:
                !type Type2:
                    ...
        `));

        const clonedField = model
            .getField("App1.Type1.Field1")
            .clone(new CloneContext(model, ModelFilters.ExcludeAnnosAndTags));
        clonedField.name = "Field2";                        // Modify name to ensure originals don't change
        model.getType("App2.Type2").children.push(clonedField);

        expect(model.toSysl()).toEqual(realign(`
            App1:
                !type Type1:
                    Field1 <: int [~tag1]:
                        @anno = "value"

            App2:
                !type Type2:
                    Field2 <: int
        `));
    });

    test.concurrent("Type", async () => {
        const model = await Model.fromText(realign(`
            App1:
                !table Type1 [~tag1]:
                    @anno = "value1"
                    Field1 <: int

            App2:
                ...
        `));

        const clonedType = model.getType("App1.Type1").clone();
        clonedType.name = "Type2";                         //
        clonedType.tags[0].name = "tag2";                  //
        clonedType.annos[0].name = "anno2";                // Modify name to ensure originals don't change
        clonedType.annos[0].value = "value2";              //
        clonedType.children[0].name = "Field2"             //

        model.getApp("App2").children.push(clonedType);
        model.attachSubitems();

        expect(model.toSysl()).toEqual(realign(`
            App1:
                !table Type1 [~tag1]:
                    @anno = "value1"
                    Field1 <: int

            App2:
                !table Type2 [~tag2]:
                    @anno2 = "value2"
                    Field2 <: int
        `));
    });

    test.concurrent("App", async () => {
        const model = await Model.fromText(realign(`
            App1 [~tag1]:
                @anno = "value1"
                !table Type1:
                    Field1 <: int
        `));

        const clonedApp = model.getApp("App1").clone();
        clonedApp.name = "App2";                           //
        clonedApp.tags[0].name = "tag2";                   //
        clonedApp.annos[0].name = "anno2";                 // Modify name to ensure originals don't change
        clonedApp.annos[0].value = "value2";               //
        clonedApp.children[0].name = "Type2"               //

        model.apps.push(clonedApp);

        expect(model.toSysl()).toEqual(realign(`
            App1 [~tag1]:
                @anno = "value1"
                !table Type1:
                    Field1 <: int

            App2 [~tag2]:
                @anno2 = "value2"
                !table Type2:
                    Field1 <: int
        `));
    });

    test.concurrent("Union", async () => {
        const model = await Model.fromText(realign(`
            App1:
                !union Union1 [~tag1]:
                    @anno = "value1"
                    int
                    string
                    sequence of decimal(5.8)
                    RestEndpoint.Type
            App2:
                ...
        `));

        const clonedUnion = model.getApp("App1").children.find(c => c.name == "Union1")!.clone() as Union;
        clonedUnion.name = "Union2";                       // Modify name to ensure originals don't change
        clonedUnion.getTag("tag1").name = "tag2";
        clonedUnion.getAnno("anno").value = "value2";
        clonedUnion.members = clonedUnion.members
            .filter(m => !(m instanceof Primitive && m.primitive == TypePrimitive.STRING));

        model.getApp("App2").children.push(clonedUnion);

        expect(model.toSysl()).toEqual(realign(`
            App1:
                !union Union1 [~tag1]:
                    @anno = "value1"
                    int
                    string
                    sequence of decimal(5.8)
                    RestEndpoint.Type

            App2:
                !union Union2 [~tag2]:
                    @anno = "value2"
                    int
                    sequence of decimal(5.8)
                    RestEndpoint.Type
        `));
    });

    test.concurrent("Alias", async () => {
        const model = await Model.fromText(realign(`
            App1:
                !alias Alias1 [~tag1]:
                    @anno = "value1"
                    int
            App2:
                ...
        `));

        const clonedAlias = model.getApp("App1").children.find(c => c.name == "Alias1")!.clone() as Alias;
        clonedAlias.name = "Alias2";                       // Modify name to ensure originals don't change
        clonedAlias.getTag("tag1").name = "tag2";
        clonedAlias.getAnno("anno").value = "value2";
        clonedAlias.value = new Primitive(TypePrimitive.STRING);

        model.getApp("App2").children.push(clonedAlias);

        expect(model.toSysl()).toEqual(realign(`
            App1:
                !alias Alias1 [~tag1]:
                    @anno = "value1"
                    int

            App2:
                !alias Alias2 [~tag2]:
                    @anno = "value2"
                    string
        `));
    });

    test.concurrent("Special chars", async () => {
        const sysl = realign(`
            example%2Ecom [~tag1]:
                @anno = "value1"
                !table www%2Eexample%2Ecom:
                    subdomain%2Eexample%2Ecom <: int
        `);
        const model = await Model.fromText(sysl);
        model.apps[0] = model.apps[0].clone();
        expect(model.toSysl()).toEqual(sysl);
    });

    test.concurrent("All", async () => {
        const model = await Model.fromFile(allPath);
        const clonedModel = model.clone();

        expect(clonedModel.toSysl()).toEqual(model.toSysl());
    });

    test.concurrent("All filtered to current file", async () => {
        const model = await Model.fromFile(allPath);
        const clonedModel = model.clone(ModelFilters.OnlyFromFile(model.convertSyslPath(allPath)));
        expect(clonedModel.findApp("ImportedApp")).toBeUndefined();

        const allSysl = (await readFile(allPath)).toString();
        expect(clonedModel.toSysl()).toEqual(allSysl);
        expect(clonedModel.toSysl()).not.toEqual(model.toSysl());
    });

    test.concurrent("Preserve location when keepLocation=true", async () => {
        const model = await Model.fromFile(allPath);
        const clonedModel = model.clone(undefined, true);

        expect(clonedModel.getApp("Types").locations)
            .toEqual(model.getApp("Types").locations);

        expect(clonedModel.getType("Types.Type").locations)
            .toEqual(model.getType("Types.Type").locations);

        expect(clonedModel.getField("Types.Type.with_anno").locations)
            .toEqual(model.getField("Types.Type.with_anno").locations);

        expect(clonedModel.getField("Types.Type.with_anno").getAnno("annotation").locations)
            .toEqual(model.getField("Types.Type.with_anno").getAnno("annotation").locations);

        expect(clonedModel.getField("Types.Type.with_anno").getTag("tag").locations)
            .toEqual(model.getField("Types.Type.with_anno").getTag("tag").locations);

        expect(clonedModel.getApp("ImportedApp").locations)
            .toEqual(model.getApp("ImportedApp").locations);


        const clonedApp = clonedModel.getApp("Types").children;
        const app = model.getApp("Types").children;

        expect((clonedApp.find(c => c.name == "Enum") as ILocational).locations)
            .toEqual((app.find(c => c.name == "Enum") as ILocational).locations);
        expect((clonedApp.find(c => c.name == "Union") as ILocational).locations)
            .toEqual((app.find(c => c.name == "Union") as ILocational).locations);
        expect((clonedApp.find(c => c.name == "Alias") as ILocational).locations)
            .toEqual((app.find(c => c.name == "Alias") as ILocational).locations);
    });

    test.concurrent("removes location when keepLocation=false", async () => {
        const model = await Model.fromFile(allPath);
        const clonedModel = model.clone();

        expect(clonedModel.getApp("Types").locations).toBeEmpty();
        expect(clonedModel.getType("Types.Type").locations).toBeEmpty();
        expect(clonedModel.getField("Types.Type.with_anno").locations).toBeEmpty();
        expect(clonedModel.getField("Types.Type.with_anno").getAnno("annotation").locations).toBeEmpty();
        expect(clonedModel.getField("Types.Type.with_anno").getTag("tag").locations).toBeEmpty();
        expect(clonedModel.getApp("ImportedApp").locations).toBeEmpty();

        const clonedApp = clonedModel.getApp("Types").children;
        expect((clonedApp.find(c => c.name == "Enum") as ILocational).locations).toBeEmpty();
        expect((clonedApp.find(c => c.name == "Union") as ILocational).locations).toBeEmpty();
        expect((clonedApp.find(c => c.name == "Alias") as ILocational).locations).toBeEmpty();
    });

    function renderAs(model: Model, renderer: (item: ICloneable | ILocational) => string | undefined): string {
        const visits: string[] = [];
        model.clone((context, item) => {
            const render = renderer(item);
            if (render) visits.push(`${"  ".repeat(context.depth - 1)}${render}`);
            return !!render;
        });
        return visits.join("\n");
    }

    test.concurrent("All filter visits", async () => {
        expect(renderAs(await Model.fromFile(allPath), i => `'${i}': ${i.constructor.name}`)).toEqual(realign(`
            'imported.sysl': Import
            'App': Application
              '[~abstract]': Tag
            'AppWithAnnotation': Application
              '[~tag]': Tag
              '@annotation = ...': Annotation
              '@annotation1 = ...': Annotation
              '@annotation2 = ...': Annotation
              '@annotation3 = ...': Annotation
            'App::with::subpackages': Application
              '[~tag]': Tag
            'RestEndpoint': Application
              '[~tag]': Tag
              '[REST] /': Endpoint
                '[~rest]': Tag
                'GET': RestParams
              '[REST] /pathwithtype/{native}': Endpoint
                '[~rest]': Tag
                'GET': RestParams
                  'native <: int': Field
                'action': ActionStatement
              '[REST] /query': Endpoint
                '[~rest]': Tag
                'GET': RestParams
                  'native <: string': Field
                  'optional <: string?': Field
                'action': ActionStatement
              '[REST] /param': Endpoint
                '[~rest]': Tag
                't <: Types.Type': Field
                  '[~body]': Tag
                'PATCH': RestParams
                'action': ActionStatement
              '[REST] /param': Endpoint
                '[~rest]': Tag
                'native <: string': Field
                'POST': RestParams
                'action': ActionStatement
              '[REST] /param': Endpoint
                '[~rest]': Tag
                'unlimited <: string(5..)': Field
                  'length: 5..': TypeConstraint
                'limited <: string(5..10)': Field
                  'length: 5..10': TypeConstraint
                'num <: int(5)': Field
                  'length: ..5': TypeConstraint
                'PUT': RestParams
                'action': ActionStatement
              '[REST] /report.csv': Endpoint
                '[~rest]': Tag
                'GET': RestParams
                'action': ActionStatement
            'SimpleEndpoint': Application
              '[~tag]': Tag
              '[RPC] SimpleEp': Endpoint
                '[~SimpleEpTag]': Tag
                '@annotation = ...': Annotation
                '@annotation1 = ...': Annotation
                '@annotation2 = ...': Annotation
                '@annotation3 = ...': Annotation
              '[RPC] SimpleEpWithParam': Endpoint
                '[~tag]': Tag
                'untypedParam <: any': Field
              '[RPC] SimpleEpWithTypes': Endpoint
                '[~tag]': Tag
                'native <: string': Field
                'action': ActionStatement
              '[RPC] SimpleEpWithArray': Endpoint
                '[~tag]': Tag
                'unlimited <: string(5..)': Field
                  'length: 5..': TypeConstraint
                'limited <: string(5..10)': Field
                  'length: 5..10': TypeConstraint
                'num <: int(5)': Field
                  'length: ..5': TypeConstraint
                'action': ActionStatement
            'Types': Application
              '!type Type': Type
                '[~tag]': Tag
                '@annotation = ...': Annotation
                'nativeTypeField <: string': Field
                  '[~tag]': Tag
                'reference <: RestEndpoint.Type': Field
                  '[~tag]': Tag
                'optional <: string?': Field
                  '[~tag]': Tag
                'set <: set of string': Field
                  '[~tag]': Tag
                'sequence <: sequence of string': Field
                  '[~tag]': Tag
                'aliasSequence <: AliasSequence': Field
                  '[~tag]': Tag
                'with_anno <: string': Field
                  '[~tag]': Tag
                  '@annotation = ...': Annotation
              '!table Table': Type
                '[~tag]': Tag
                'primaryKey <: string': Field
                  '[~pk]': Tag
                'nativeTypeField <: string': Field
                'reference <: RestEndpoint.Type': Field
                'optional <: string?': Field
                'set <: set of string': Field
                'sequence <: sequence of string': Field
                'with_anno <: string': Field
                  '@annotation = ...': Annotation
                'decimal_with_precision <: decimal(5.8)': Field
                  'length: ..5, precision: 5, scale: 8': TypeConstraint
                'string_max_constraint <: string(5)': Field
                  'length: ..5': TypeConstraint
                'string_range_constraint <: string(5..10)': Field
                  'length: 5..10': TypeConstraint
                'int_with_bitwidth <: int64': Field
                  'bitWidth: 64': TypeConstraint
                'float_with_bitwidth <: float64': Field
                  'bitWidth: 64': TypeConstraint
              '!enum Enum': Enum
                '[~tag]': Tag
                'ENUM_1: 1': EnumValue
                'ENUM_2: 2': EnumValue
                'ENUM_3: 3': EnumValue
              '!union Union': Union
                '[~tag]': Tag
                'int': Primitive
                'string': Primitive
                'sequence of decimal(5.8)': CollectionDecorator
                  'length: ..5, precision: 5, scale: 8': TypeConstraint
                'RestEndpoint.Type': ElementRef
              '!union EmptyUnion': Union
                '[~tag]': Tag
              '!Alias Alias': Alias
                '[~tag]': Tag
                '@annotation1 = ...': Annotation
                '@annotation2 = ...': Annotation
                '@annotation3 = ...': Annotation
              '!Alias AliasSequence': Alias
                '[~tag]': Tag
                '@annotation1 = ...': Annotation
                '@annotation2 = ...': Annotation
                '@annotation3 = ...': Annotation
              '!Alias AliasRef': Alias
                '[~tag]': Tag
              '!Alias AliasForeignRef': Alias
                '[~tag]': Tag
              '!Alias AliasForeignRefSet': Alias
                '[~tag]': Tag
            'Statements': Application
              '[~tag]': Tag
              '[RPC] IfStmt': Endpoint
                '[~tag]': Tag
                'if predicate1': CondStatement
                  'return ok <: string': ReturnStatement
                'else if predicate2': GroupStatement
                  'Statements <- IfStmt': CallStatement
                'else': GroupStatement
              '[RPC] Loops': Endpoint
                '[~tag]': Tag
                'until predicate': LoopStatement
                  'action': ActionStatement
                'for each predicate': ForEachStatement
                'while predicate': LoopStatement
              '[RPC] Returns': Endpoint
                '[~tag]': Tag
                'return ok <: string': ReturnStatement
                'return ok <: Types.Type': ReturnStatement
                'return error <: Types.Type': ReturnStatement
              '[RPC] Calls': Endpoint
                '[~tag]': Tag
                'Statements <- Returns': CallStatement
                'RestEndpoint <- GET /param': CallStatement
              '[RPC] OneOfStatements': Endpoint
                '[~tag]': Tag
                'one of': OneOfStatement
                  'case1': GroupStatement
                    'return ok <: string': ReturnStatement
                  'case number 2': GroupStatement
                    'return ok <: int': ReturnStatement
                  '\"case 3\"': GroupStatement
                    'return ok <: Types.Type': ReturnStatement
                  '': GroupStatement
                    'return error <: string': ReturnStatement
              '[RPC] GroupStatements': Endpoint
                '[~tag]': Tag
                'grouped': GroupStatement
                  'Statements <- GroupStatements': CallStatement
              '[RPC] AnnotatedEndpoint': Endpoint
                '[~tag]': Tag
                '@annotation1 = ...': Annotation
                '@annotation2 = ...': Annotation
                '@annotation3 = ...': Annotation
              '[RPC] AnnotatedStatements': Endpoint
                'Statements <- Miscellaneous': CallStatement
                'return ok <: string [annotation=[\"as\", \"an\", \"array\"]] #Doesn't work, annos/tags/comments are part of the name': ReturnStatement
                '\"statement\"': ActionStatement
              '[RPC] Miscellaneous': Endpoint
                'SimpleEndpoint -> SimpleEp': ActionStatement
            'Unsafe%2FNamespace::Unsafe%2FApp': Application
              '[~tag]': Tag
              '!type Unsafe%2EType': Type
                '[~tag]': Tag
                'Unsafe.Field <: int': Field
                  '[~tag]': Tag
                  '@description = ...': Annotation
            'ImportedApp': Application`, 2))
    });

    
    const hasRef = (item: any): item is Element => {
        return (item instanceof Element) &&
            // TODO: Remove when all elements have toRef();
            !(
                item instanceof Enum ||
                item instanceof Union ||
                item instanceof Alias ||
                (item instanceof Field && (!item.parent || item.parent instanceof Endpoint)) // Endpoint param
            );
    }

    test.concurrent("All .toRef()s", async () => {

        expect(renderAs(await Model.fromFile(allPath), i => hasRef(i) ? i.toRef().toString() : undefined))
            .toEqual(realign(`
            App
            AppWithAnnotation
            App::with::subpackages
            RestEndpoint
              RestEndpoint.[GET /]
              RestEndpoint.[GET /pathwithtype/{native}]
                RestEndpoint.[GET /pathwithtype/{native}].[0]
              RestEndpoint.[GET /query]
                RestEndpoint.[GET /query].[0]
              RestEndpoint.[PATCH /param]
                RestEndpoint.[PATCH /param].[0]
              RestEndpoint.[POST /param]
                RestEndpoint.[POST /param].[0]
              RestEndpoint.[PUT /param]
                RestEndpoint.[PUT /param].[0]
              RestEndpoint.[GET /report.csv]
                RestEndpoint.[GET /report.csv].[0]
            SimpleEndpoint
              SimpleEndpoint.[SimpleEp]
              SimpleEndpoint.[SimpleEpWithParam]
              SimpleEndpoint.[SimpleEpWithTypes]
                SimpleEndpoint.[SimpleEpWithTypes].[0]
              SimpleEndpoint.[SimpleEpWithArray]
                SimpleEndpoint.[SimpleEpWithArray].[0]
            Types
              Types.Type
                Types.Type.nativeTypeField
                Types.Type.reference
                Types.Type.optional
                Types.Type.set
                Types.Type.sequence
                Types.Type.aliasSequence
                Types.Type.with_anno
              Types.Table
                Types.Table.primaryKey
                Types.Table.nativeTypeField
                Types.Table.reference
                Types.Table.optional
                Types.Table.set
                Types.Table.sequence
                Types.Table.with_anno
                Types.Table.decimal_with_precision
                Types.Table.string_max_constraint
                Types.Table.string_range_constraint
                Types.Table.int_with_bitwidth
                Types.Table.float_with_bitwidth
            Statements
              Statements.[IfStmt]
                Statements.[IfStmt].[0]
                  Statements.[IfStmt].[0,0]
                Statements.[IfStmt].[1]
                  Statements.[IfStmt].[1,0]
                Statements.[IfStmt].[2]
              Statements.[Loops]
                Statements.[Loops].[0]
                  Statements.[Loops].[0,0]
                Statements.[Loops].[1]
                Statements.[Loops].[2]
              Statements.[Returns]
                Statements.[Returns].[0]
                Statements.[Returns].[1]
                Statements.[Returns].[2]
              Statements.[Calls]
                Statements.[Calls].[0]
                Statements.[Calls].[1]
              Statements.[OneOfStatements]
                Statements.[OneOfStatements].[0]
                  Statements.[OneOfStatements].[0,0]
                    Statements.[OneOfStatements].[0,0,0]
                  Statements.[OneOfStatements].[0,1]
                    Statements.[OneOfStatements].[0,1,0]
                  Statements.[OneOfStatements].[0,2]
                    Statements.[OneOfStatements].[0,2,0]
                  Statements.[OneOfStatements].[0,3]
                    Statements.[OneOfStatements].[0,3,0]
              Statements.[GroupStatements]
                Statements.[GroupStatements].[0]
                  Statements.[GroupStatements].[0,0]
              Statements.[AnnotatedEndpoint]
              Statements.[AnnotatedStatements]
                Statements.[AnnotatedStatements].[0]
                Statements.[AnnotatedStatements].[1]
                Statements.[AnnotatedStatements].[2]
              Statements.[Miscellaneous]
                Statements.[Miscellaneous].[0]
            Unsafe%2FNamespace::Unsafe%2FApp
              Unsafe%2FNamespace::Unsafe%2FApp.Unsafe%2EType
                Unsafe%2FNamespace::Unsafe%2FApp.Unsafe%2EType.Unsafe%2EField
            ImportedApp`, 2))
    });

    test.concurrent("All .toRef() used in get", async () => {
        const model = await Model.fromFile(allPath);
        model.clone((_context, item) => {
            if (hasRef(item)) model.getElement(item.toRef().toString());
            return true;
        });
    });
});

describe("Imports", () => {
    jest.setTimeout(20000);
    test.concurrent("Remote import without fetch", async () => {
        const sysl = "import //github.com/anz-bank/sysl/ts/test/imported.sysl";
        await Model.fromText(sysl);  // warm up

        const t = process.hrtime.bigint();
        await Model.fromText(sysl);
        const withFetchMs = Number((process.hrtime.bigint() - t) / 1000000n);

        const t2 = process.hrtime.bigint();
        await Model.fromText(sysl, undefined, { alwaysFetch: false });
        const withoutFetchMs = Number((process.hrtime.bigint() - t2) / 1000000n);
        
        expect(withFetchMs / 2).toBeGreaterThan(withoutFetchMs);
    });
});

describe("PubSub", () => {
    test.concurrent("Simple", async () => {
        const model = await Model.fromText(realign(`
        From:
            <-> Event: ...
            Endpoint:
                . <- Event
        To:
            From -> Event:
                receive
        `));

        expect(model.getApp("From").children).toMatchObject([
            { name: "Event", isPubsub: true },
            { name: "Endpoint", isPubsub: false },
        ]);

        expect(model.getApp("To").children).toMatchObject([ { name: "From -> Event", isPubsub: false } ]);
    });
});
