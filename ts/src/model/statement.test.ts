import { realign } from "../common/format";
import { Model } from "./model";

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
