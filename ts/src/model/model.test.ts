import "jest-extended";
import { readFile } from "fs/promises";
import { realign } from "../common/format";
import { Application, AppName, Model } from "./model";
import "./renderers";
import { Primitive, Type, TypePrimitive } from "./type";
import { Action, Endpoint, Param, Statement } from "./statement";
import { Annotation, Tag } from "./attribute";

const allPath = "../ts/test/all.sysl";
let allModel: Model;
let allSysl: string;

describe("Constructors", () => {
    test("New Model", () => {
        expect(new Model({})).not.toBeNull();
    });

    test("New Application", () => {
        const name = AppName.fromString("Foo");
        expect(new Application({ name })).toHaveProperty("name", name);
    });

    test("New Type", () => {
        expect(
            new Type({
                discriminator: "!type",
                name: "Foo",
                opt: true,
                value: new Primitive(TypePrimitive.INT),
            })
        ).toHaveProperty("name", "Foo");
    });

    test("New Endpoint", () => {
        expect(new Endpoint({ name: "Foo" })).toHaveProperty("name", "Foo");
    });

    test("New Param", () => {
        expect(new Param("foo")).toHaveProperty("name", "foo");
    });

    test("New Statement", () => {
        expect(new Statement({ value: new Action("foo") })).toHaveProperty(
            "value.action",
            "foo"
        );

        expect(Statement.action("foo")).toHaveProperty("value.action", "foo");
    });

    test("New Annotation", () => {
        expect(new Annotation({ name: "foo", value: "bar" })).toMatchObject({
            name: "foo",
            value: "bar",
        });
    });

    test("New Tag", () => {
        expect(new Tag({ value: "foo" })).toHaveProperty("value", "foo");
    });
});

describe("Serialization", () => {
    describe("Application", () => {
        test("empty", () => {
            expect(
                new Application({ name: AppName.fromString("Foo") }).toSysl()
            ).toEqual(
                realign(`
                    Foo:
                        ...`)
            );
        });
    });
});

describe("Roundtrip", () => {
    // All
    test("AllRoundtrip", async () => {
        allModel = await Model.fromFile(allPath);
        allSysl = (await readFile(allPath)).toString();
        expect(allModel.filterByFile(allPath).toSysl()).toEqual(allSysl);
    });

    // Applications
    test("EmptyApp", async () => {
        const sysl = realign(`
    App:
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EmptyAppWithSubpackages", async () => {
        const sysl = realign(`
    App :: with :: subpackages:
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("AppWithTag", async () => {
        const sysl = realign(`
    App [~abstract]:
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("InlineAnno", async () => {
        const sysl = realign(`
    App [name="value"]:
        ...
    `);
        const output = realign(`
    App:
        @name = "value"
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(output);
    });

    test("StringAnno", async () => {
        const sysl = realign(`
    App:
        @name = "value"
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("MultilineAnno", async () => {
        const sysl = realign(`
    App:
        @name =:
            | anno
            |  indented
            |   across
            |
            |    multiple lines
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("ArrayAnno", async () => {
        const sysl = realign(`
    App:
        @name = ["value1", "value2"]
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("NestedArrayAnno", async () => {
        const sysl = realign(`
    App:
        @name = [["value1", "value2"], ["value3", "value4"]]
        ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    // GRPC Endpoints
    test("Endpoint", async () => {
        const sysl = realign(`
    App:
        SimpleEp:
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithTag", async () => {
        const sysl = realign(`
    App:
        SimpleEp [~ignore]:
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithAnno", async () => {
        const sysl = realign(`
    App:
        SimpleEp:
            @name = "value"
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithInlineAnno", async () => {
        const sysl = realign(`
    App:
        SimpleEp [name="value"]:
            ...
    `);
        const output = realign(`
    App:
        SimpleEp:
            @name = "value"
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(output);
    });

    test("EndpointWithPrimitiveParam", async () => {
        const sysl = realign(`
    App:
        SimpleEp (param <: string):
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithRefParam", async () => {
        const sysl = realign(`
    App:
        SimpleEp (Types.type):
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithPrimitiveParamWithConstraints", async () => {
        const sysl = realign(`
    App:
        SimpleEp (unlimited <: string(5..), limited <: string(5..10), num <: int(5)):
            ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithCall", async () => {
        const sysl = realign(`
    App:
        SimpleEp:
            App2 <- Endpoint
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithPrimitiveReturn", async () => {
        const sysl = realign(`
    App:
        SimpleEp:
            return ok <: string
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("EndpointWithRefReturn", async () => {
        const sysl = realign(`
    App:
        SimpleEp:
            return ok <: Types.type
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    // REST Endpoints
    test("RestEndpoint", async () => {
        const sysl = realign(`
    RestEndpoint:
        /:
            GET:
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("RestEndpointWithoutNesting", async () => {
        const sysl = realign(`
    RestEndpoint:
        /nested/path:
            GET:
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("RestEndpointWithNesting", async () => {
        const sysl = realign(`
    RestEndpoint:
        /nested:
            /path:
                GET:
                    ...
    `);
        // Currently nested endpoints are not supported so this is the expected output
        const expected = realign(`
    RestEndpoint:
        /nested/path:
            GET:
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(expected);
    });

    test("RestEndpointWithTypeInPath", async () => {
        const sysl = realign(`
    RestEndpoint:
        /pathwithtype/{native <: int}:
            GET:
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("RestEndpointWithQueryParams", async () => {
        const sysl = realign(`
    RestEndpoint:
        /query:
            GET?native=string&optional=string?:
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("RestEndpointWithRefParam", async () => {
        const sysl = realign(`
    RestEndpoint:
        /param:
            PATCH (t <: Types.Type [~body]):
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("RestEndpointWithPrimitiveParam", async () => {
        const sysl = realign(`
    RestEndpoint:
        /param:
            POST (native <: string):
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("RestEndpointWithConstrainedParams", async () => {
        const sysl = realign(`
    RestEndpoint:
        /param:
            PUT (unlimited <: string(5..), limited <: string(5..10), num <: int(5)):
                ...
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    // Types
    test("Type", async () => {
        const sysl = realign(`
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
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("Table", async () => {
        const sysl = realign(`
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
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("Enum", async () => {
        const sysl = realign(`
    App:
        !enum Enum [~tag]:
            ENUM_1: 1
            ENUM_2: 2
            ENUM_3: 3
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("TypeRef", async () => {
        const sysl = realign(`
    Namespace :: App:
        !type Type:
            shortRef <: Type
            fullRef <: Namespace :: App.Type
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });

    test("UnsafeNames", async () => {
        const sysl = realign(`
    %28App%29Name%21:
        !type %28Type%29Name%21:
            %28Field%29Name%21 <: %28App%29Name%21.%28Type%29Name%21 [~%28Tag%29Name%21]
    `);
        const model = await Model.fromText(sysl);
        expect(model.toSysl()).toEqual(sysl);
    });
});
