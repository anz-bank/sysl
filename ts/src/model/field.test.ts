import { realign } from "../common/format";
import { Model } from "./model";
import { Primitive } from "./primitive";

test.concurrent("Untyped field", async () => {
    const model = await Model.fromText(realign(`
        App:
            !type Type:
                Untyped
    `));
    expect(model.getField("App.Type.Untyped").value).toEqual(Primitive.Any);
    expect(model.toSysl()).toEqual(realign(`
        App:
            !type Type:
                Untyped <: any
    `));
});
