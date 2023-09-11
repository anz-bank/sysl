import "reflect-metadata";
import { readFile } from "fs/promises";
import "jest-extended";
import { realign } from "../common/format";
import { allItems } from "../common/iterate";
import { Annotation, AnnoValue, Tag } from "./attribute";
import { Model } from "./model";
import "./renderers";
import { Action, Endpoint, Group, Param, Statement } from "./statement";
import { Type } from "./type";
import { Application } from "./application";
import { Field } from "./field";
import { ElementRef, ILocational } from "./common";
import { CloneContext, ModelFilters } from "./clone";
import { Union } from "./union";
import { Primitive, TypePrimitive } from "./primitive";
import { Alias } from "./alias";
import { CollectionDecorator } from "./decorator";

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

    test.concurrent("New Param", () => {
        expect(new Param("foo", [])).toHaveProperty("name", "foo");
    });

    test.concurrent("New Statement", () => {
        expect(new Statement(new Action("foo"))).toHaveProperty("value.action", "foo");

        expect(Statement.action("foo")).toHaveProperty("value.action", "foo");
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
        test.concurrent("typed with absolute path", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        HomeAddress <: Company::App.Address
                    !type Address:
                        City <: string
            `));
            expect(model.getField("Company::App.Customer.HomeAddress").value).toEqual(model.getType("Company::App.Address").toRef());
        });

        test.concurrent("typed with relative path", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        HomeAddress <: Address
                    !type Address:
                        City <: string
            `));
            expect(model.getField("Company::App.Customer.HomeAddress").value).toEqual(model.getType("Company::App.Address").toRef());
        });

        test.concurrent("typed with collection with absolute path", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        HomeAddress <: sequence of Company::App.Address
                    !type Address:
                        City <: string
            `));
            const fieldType = model.getField("Company::App.Customer.HomeAddress").value as CollectionDecorator;
            expect(fieldType).toEqual(new CollectionDecorator(model.getType("Company::App.Address").toRef(), false));
        });

        test.concurrent("typed with collection with relative path", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        HomeAddress <: sequence of Address
                    !type Address:
                        City <: string
            `));
            const fieldType = model.getField("Company::App.Customer.HomeAddress").value as CollectionDecorator;
            expect(fieldType).toEqual(new CollectionDecorator(model.getType("Company::App.Address").toRef(), false));
        });

        test.concurrent("typed with collection-list with relative path", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        HomeAddress(1..1) <: sequence of Address
                    !type Address:
                        City <: string
            `));
            const fieldType = model.getField("Company::App.Customer.HomeAddress").value as CollectionDecorator;
            expect(fieldType).toEqual(new CollectionDecorator(model.getType("Company::App.Address").toRef(), false));
        });

        test.concurrent("typed with set-list with relative path", async () => {
            const model = await Model.fromText(realign(`
                Company :: App:
                    !table Customer:
                        HomeAddress(1..1) <: set of Address
                    !type Address:
                        City <: string
            `));
            const fieldType = model.getField("Company::App.Customer.HomeAddress").value as CollectionDecorator;
            expect(fieldType).toEqual(new CollectionDecorator(model.getType("Company::App.Address").toRef(), true));
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
                        namespace: ["Ns"],
                        name: "App",
                        metadata: { appTag: undefined, appAnno: "App" },
                        locations: {
                            "0": "ts/main.sysl:5:1:12:33",
                            appTag: "ts/main.sysl:5:12::19",
                            appAnno: "ts/main.sysl:6:16::21"
                        },
                        children: [
                            {
                                name: "Type",
                                metadata: { typeTag: undefined, typeAnno: "Type" },
                                locations: {
                                    "0": "ts/main.sysl:7:5:12:33",
                                    typeTag: "ts/main.sysl:7:17::25",
                                    typeAnno: "ts/main.sysl:8:21::27"
                                },
                                children: [
                                    {
                                        name: "Primitive",
                                        primitive: "string",
                                        constraint: "(5..10)",
                                        locations: { "0": "ts/main.sysl:9:9::35" }
                                    },
                                    {
                                        name: "Array",
                                        primitive: "decimal",
                                        constraint: "(5.8)",
                                        collectionType: "sequence",
                                        locations: { "0": "ts/main.sysl:10:9::42" }
                                    },
                                    {
                                        name: "Reference",
                                        ref: "Other.Type",
                                        metadata: { fieldTag: undefined, fieldAnno: "Field" },
                                        locations: {
                                            "0": "ts/main.sysl:11:9:13:2",
                                            fieldTag: "ts/main.sysl:11:34::43",
                                            fieldAnno: "ts/main.sysl:12:26::33"
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
        const outerSt = ep.statements[0];
        const innerSt = new Statement(new Group("subStatement", []));
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
                    ...
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
                        ...
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
        EndpointWithUnnamedRefParam: realign(`
            App:
                SimpleEp (Types.type):
                    ...
            `),
        EndpointWithNamedRefParam: realign(`
            App:
                SimpleEp (param <: Types.type):
                    ...
            `),
        // TODO: Support relative refs in endpoint params
        EndpointWithLocalNamedRefParam: realign(`
            App:
                SimpleEp (param <: App.type):
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
        // TODO: Support relative refs in endpoint params
        RestEndpointWithComplexParamType: realign(`
            RestEndpoint:
                /param:
                    POST (arg <: RestEndpoint.Customer):
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
        `));

        // Element not in model
        expect(model.findElement("Namespace::MissingApp")).toBeUndefined();
        expect(model.findElement("Namespace::App.MissingType")).toBeUndefined();
        expect(model.findElement("Namespace::App.Type.MissingField")).toBeUndefined();
        expect(model.findApp("Namespace::MissingApp")).toBeUndefined();
        expect(model.findType("Namespace::App.MissingType")).toBeUndefined();
        expect(model.findField("Namespace::App.Type.MissingField")).toBeUndefined();
        expect(() => model.getApp("Namespace::MissingApp")).toThrow();
        expect(() => model.getType("Namespace::App.MissingType")).toThrow();
        expect(() => model.getField("Namespace::App.Type.MissingField")).toThrow();

        // Invalid element reference
        expect(() => model.findElement("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findApp("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findType("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.findField("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.getApp("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.getType("Namespace::App.Type.Field.Crazy")).toThrow();
        expect(() => model.getField("Namespace::App.Type.Field.Crazy")).toThrow();

        // Element type mismatch
        expect(() => model.findApp("Namespace::App.Type")).toThrow();
        expect(() => model.findType("Namespace::App.Type.Field")).toThrow();
        expect(() => model.findField("Namespace::App")).toThrow();
        expect(() => model.getApp("Namespace::App.Type")).toThrow();
        expect(() => model.getType("Namespace::App.Type.Field")).toThrow();
        expect(() => model.getField("Namespace::App")).toThrow();

        // Happy path
        expect(model.findElement("Namespace::App")).toBeInstanceOf(Application);
        expect(model.findApp("Namespace::App")).toBeInstanceOf(Application);
        expect(model.getApp("Namespace::App")).toBeInstanceOf(Application);

        expect(model.findElement("Namespace::App.Type")).toBeInstanceOf(Type);
        expect(model.findType("Namespace::App.Type")).toBeInstanceOf(Type);
        expect(model.getType("Namespace::App.Type")).toBeInstanceOf(Type);

        expect(model.findElement("Namespace::App.Type.Field")).toBeInstanceOf(Field);
        expect(model.findField("Namespace::App.Type.Field")).toBeInstanceOf(Field);
        expect(model.getField("Namespace::App.Type.Field")).toBeInstanceOf(Field);
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

        expect(model.toSysl()).toEqual(clonedModel.toSysl());
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

    test.concurrent("All filter visits", async () => {
        const model = await Model.fromFile(allPath);
        const visits: string[] = [];
        model.clone((context, item) => {
            visits.push(`${"  ".repeat(context.depth - 1)}'${(item as any)?.toString()}': ${item.constructor.name}`);
            return true;
        });

        expect(visits.join("\n")).toEqual(realign(`
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
                '...': Statement
                  '...': Action
              '[REST] /pathwithtype/{native}': Endpoint
                '[~rest]': Tag
                'GET': RestParams
                  'native': Param
                    '[param] <: int': Field
                '...': Statement
                  '...': Action
              '[REST] /query': Endpoint
                '[~rest]': Tag
                'GET': RestParams
                  'native': Param
                    '[param] <: string': Field
                  'optional': Param
                    '[param] <: string?': Field
                '...': Statement
                  '...': Action
              '[REST] /param': Endpoint
                '[~rest]': Tag
                't': Param
                  '[param] <: Types.Type': Field
                    '[~body]': Tag
                'PATCH': RestParams
                '...': Statement
                  '...': Action
              '[REST] /param': Endpoint
                '[~rest]': Tag
                'native': Param
                  '[param] <: string': Field
                'POST': RestParams
                '...': Statement
                  '...': Action
              '[REST] /param': Endpoint
                '[~rest]': Tag
                'unlimited': Param
                  '[param] <: string(5..)': Field
                    'length: 5..': TypeConstraint
                'limited': Param
                  '[param] <: string(5..10)': Field
                    'length: 5..10': TypeConstraint
                'num': Param
                  '[param] <: int(5)': Field
                    'length: ..5': TypeConstraint
                'PUT': RestParams
                '...': Statement
                  '...': Action
              '[REST] /report.csv': Endpoint
                '[~rest]': Tag
                'GET': RestParams
                '...': Statement
                  '...': Action
            'SimpleEndpoint': Application
              '[~tag]': Tag
              '[gRPC] SimpleEp': Endpoint
                '[~SimpleEpTag]': Tag
                '@annotation = ...': Annotation
                '@annotation1 = ...': Annotation
                '@annotation2 = ...': Annotation
                '@annotation3 = ...': Annotation
              '[gRPC] SimpleEpWithParamsRef': Endpoint
                '[~tag]': Tag
                'Types.type': Param
                '...': Statement
                  '...': Action
              '[gRPC] SimpleEpWithTypes': Endpoint
                '[~tag]': Tag
                'native': Param
                  '[param] <: string': Field
                '...': Statement
                  '...': Action
              '[gRPC] SimpleEpWithArray': Endpoint
                '[~tag]': Tag
                'unlimited': Param
                  '[param] <: string(5..)': Field
                    'length: 5..': TypeConstraint
                'limited': Param
                  '[param] <: string(5..10)': Field
                    'length: 5..10': TypeConstraint
                'num': Param
                  '[param] <: int(5)': Field
                    'length: ..5': TypeConstraint
                '...': Statement
                  '...': Action
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
                'aliasSequence <: Types.AliasSequence': Field
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
              '[gRPC] IfStmt': Endpoint
                '[~tag]': Tag
                'if predicate1': Statement
                  'if predicate1': Cond
                    'return ok <: string': Statement
                      'return ok <: string': Return
                'else if predicate2': Statement
                  'else if predicate2': Group
                    'Statements <- IfStmt': Statement
                      'Statements <- IfStmt': Call
                'else': Statement
                  'else': Group
                    '...': Statement
                      '...': Action
              '[gRPC] Loops': Endpoint
                '[~tag]': Tag
                'alt predicate': Statement
                  'alt predicate': Group
                    '...': Statement
                      '...': Action
                '[UNTIL] predicate': Statement
                  '[UNTIL] predicate': Loop
                    '...': Statement
                      '...': Action
                'predicate': Statement
                  'predicate': Foreach
                    '...': Statement
                      '...': Action
                'for predicate': Statement
                  'for predicate': Group
                    '...': Statement
                      '...': Action
                'loop predicate': Statement
                  'loop predicate': Group
                    '...': Statement
                      '...': Action
                '[WHILE] predicate': Statement
                  '[WHILE] predicate': Loop
                    '...': Statement
                      '...': Action
              '[gRPC] Returns': Endpoint
                '[~tag]': Tag
                'return ok <: string': Statement
                  'return ok <: string': Return
                'return ok <: Types.Type': Statement
                  'return ok <: Types.Type': Return
                'return error <: Types.Type': Statement
                  'return error <: Types.Type': Return
              '[gRPC] Calls': Endpoint
                '[~tag]': Tag
                'Statements <- Returns': Statement
                  'Statements <- Returns': Call
                'RestEndpoint <- GET /param': Statement
                  'RestEndpoint <- GET /param': Call
              '[gRPC] OneOfStatements': Endpoint
                '[~tag]': Tag
                '[case1,case number 2,\"case 3\",undefined choices]': Statement
                  '[case1,case number 2,\"case 3\",undefined choices]': Alt
                    'case1': AltChoice
                      'return ok <: string': Statement
                        'return ok <: string': Return
                    'case number 2': AltChoice
                      'return ok <: int': Statement
                        'return ok <: int': Return
                    '\"case 3\"': AltChoice
                      'return ok <: Types.Type': Statement
                        'return ok <: Types.Type': Return
                    'undefined': AltChoice
                      'return error <: string': Statement
                        'return error <: string': Return
              '[gRPC] GroupStatements': Endpoint
                '[~tag]': Tag
                'grouped': Statement
                  'grouped': Group
                    'Statements <- GroupStatements': Statement
                      'Statements <- GroupStatements': Call
              '[gRPC] AnnotatedEndpoint': Endpoint
                '[~tag]': Tag
                '@annotation1 = ...': Annotation
                '@annotation2 = ...': Annotation
                '@annotation3 = ...': Annotation
              '[gRPC] AnnotatedStatements': Endpoint
                'Statements <- Miscellaneous': Statement
                  'Statements <- Miscellaneous': Call
                'return ok <: string [annotation=[\"as\", \"an\", \"array\"]] #Doesn't work, annos/tags/comments are part of the name': Statement
                  'return ok <: string [annotation=[\"as\", \"an\", \"array\"]] #Doesn't work, annos/tags/comments are part of the name': Return
                '\"statement\"': Statement
                  '\"statement\"': Action
              '[gRPC] Miscellaneous': Endpoint
                'SimpleEndpoint -> SimpleEp': Statement
                  'SimpleEndpoint -> SimpleEp': Action
            'Unsafe%2FNamespace::Unsafe%2FApp': Application
              '[~tag]': Tag
              '!type Unsafe%2EType': Type
                '[~tag]': Tag
                'Unsafe.Field <: int': Field
                  '[~tag]': Tag
                  '@description = ...': Annotation
            'ImportedApp': Application`, 2))
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
