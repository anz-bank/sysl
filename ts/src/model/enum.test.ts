import { realign } from "../common/format";
import { Enum } from "./enum";
import { Model } from "./model";
import "jest-extended";

test.concurrent("Enum value is number", async () => {
    const model = await Model.fromText(realign(`
        App:
            !enum Enum:
                A: 5
    `));
    const enumVal = (model.getApp("App").children[0] as Enum).children[0].value;
    expect(enumVal).toBeNumber();
    expect(enumVal).toBe(5);
});
