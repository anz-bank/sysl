import { realign } from "../common/format";
import { Endpoint } from "./endpoint";
import { Model } from "./model";
import { RestMethod, RestParams } from "./statement";

describe("Rendering", () => {
    test.concurrent("Empty RPC", () => {
        var ep = new Endpoint("Endpoint");
        expect(ep.toSysl()).toEqual(realign(`
            Endpoint:
                ...`));
    });

    test.concurrent("Empty REST", () => {
        var ep = new Endpoint("GET /", { restParams: new RestParams({ method: RestMethod.GET, path: "/" }) });
        expect(ep.toSysl()).toEqual(realign(`
            /:
                GET:
                    ...`));
    });
});

describe("Parsing", () => {
    test.concurrent("Children", async () => {
        const model = await Model.fromText(realign(`
            App:
                Endpoint:
                    Statement
        `));

        expect(model.getEndpoint("App.[Endpoint]").children).toHaveLength(1);
    });
});
