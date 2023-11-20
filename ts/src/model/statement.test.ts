import "jest-extended";
import { realign } from "../common/format";
import { Model } from "./model";
import { ElementRef } from "./elementRef";
import { OneOfStatement } from "./statement";

test.concurrent("Zero-statement endpoint, no annos", async () => {
    const sysl = realign(`
        App:
            Endpoint:
                ...
    `);
    const model = await Model.fromText(sysl);

    expect(model.getEndpoint("App.[Endpoint]").children).toBeEmpty();
    expect(model.toSysl()).toBe(sysl);
});

test.concurrent("Zero-statement endpoint, with annos", async () => {
    const sysl = realign(`
        App:
            Endpoint:
                @anno = "value"
    `);
    const model = await Model.fromText(sysl);

    expect(model.getEndpoint("App.[Endpoint]").children).toBeEmpty();
    expect(model.toSysl()).toBe(sysl);
});

test.concurrent("Endpoint param, foreign app", async () => {
    const sysl = realign(`
        App:
            Endpoint (param <: OtherApp.Type):
                ...
    `);
    const model = await Model.fromText(sysl);

    expect(model.getEndpoint("App.[Endpoint]").params[0].value).toEqual(ElementRef.parse("OtherApp.Type"));
    expect(model.toSysl()).toBe(sysl);
});

test.concurrent("Endpoint param, local app shorthand", async () => {
    const sysl = realign(`
        App:
            Endpoint (param <: Type):
                ...
    `);
    const model = await Model.fromText(sysl);

    expect(model.getEndpoint("App.[Endpoint]").params[0].value).toEqual(ElementRef.parse("App.Type"));
    expect(model.toSysl()).toBe(sysl);
});

test.concurrent("Endpoint param, local app fully-qualified", async () => {
    const sysl = realign(`
        App:
            Endpoint (param <: App.Type):
                ...
    `);
    const model = await Model.fromText(sysl);

    expect(model.getEndpoint("App.[Endpoint]").params[0].value).toEqual(ElementRef.parse("App.Type"));
    expect(model.toSysl()).toBe(realign(`
        App:
            Endpoint (param <: Type):
                ...
    `));
});

test.concurrent("REST endpoint special chars", async () => {
    const sysl = realign(`
        App:
            /rest%2Ephp:
                GET:
                    ...
    `);
    const model = await Model.fromText(sysl);

    expect(model.getApp("App").endpoints[0].name).toBe("GET /rest.php");
    expect(model.toSysl()).toBe(sysl);
});

test.concurrent("one of choice with empty name", async () => {
    const sysl = realign(`
        App:
            Endpoint:
                one of:
                    :
                        ...
    `);
    const model = await Model.fromText(sysl);

    expect((model.getStatement("App.[Endpoint].[0,0]") as OneOfStatement).title).toEqual("");
    expect(model.toSysl()).toBe(sysl);
});

// TODO: After parser is fixed to accept quoted group names (like one-of choices), restore test
test.failing("Quoted names for choice and group", async () => {
    const sysl = realign(`
        App:
            Endpoint:
                one of:
                    unquotedChoice:
                        ...
                    "quotedChoice":
                        ...
                unquotedGroup:
                    "quotedGroup":
                        ...
    `);
    const model = await Model.fromText(sysl);

    expect((model.getStatement("App.[Endpoint].[0,0]") as OneOfStatement).title).toEqual("");
    expect(model.toSysl()).toBe(sysl);
});
