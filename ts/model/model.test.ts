import { readFile } from "fs/promises";
import { realign } from "../format";
import { Model } from "./model";
import "./renderers";
import "jest-extended";

const allPath = "ts/test/all.sysl";
let allModel: Model;
let allSysl: string;

beforeAll(async () => {
    process.chdir("../");
    allModel = await Model.fromFile(allPath);
    allSysl = (await readFile(allPath)).toString();
});

test("AllRoundtrip", () => {
    expect(allModel.filterByFile(allPath).toSysl()).toEqual(allSysl);
});

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
            | across
            | multiple lines
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
