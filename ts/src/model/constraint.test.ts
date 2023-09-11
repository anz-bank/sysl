import { realign } from "../common/format";
import { TypeConstraint } from "./constraint";
import { Model } from "./model";
import { Primitive } from "./primitive";
import "jest-extended";

//TODO: Remove skip from these tests when bugs are fixed in Sysl binary

test.skip("Decimal with precision and scale", async () => {
    const model = await Model.fromText(realign(`
        App1:
            !type Type1:
                Field1 <: decimal(5.2)
    `));

    const primitive = model.getField("App1.Type1.Field1").value as Primitive;
    expect(primitive).toBeInstanceOf(Primitive);

    const constraint = primitive.constraint!;
    expect(constraint).toBeInstanceOf(TypeConstraint);

    expect(constraint.precision).toBe(5);
    expect(constraint.scale).toBe(2);

    expect(constraint.bitWidth).toBeUndefined();
    expect(constraint.length).toBeUndefined();
});

test.skip("Integer with precision", async () => {
    const model = await Model.fromText(realign(`
        App1:
            !type Type1:
                Field1 <: int(5)
    `));

    const primitive = model.getField("App1.Type1.Field1").value as Primitive;
    expect(primitive).toBeInstanceOf(Primitive);

    const constraint = primitive.constraint!;
    expect(constraint).toBeInstanceOf(TypeConstraint);

    expect(constraint.precision).toBe(5);
    expect(constraint.scale).toBeUndefined();

    expect(constraint.bitWidth).toBeUndefined();
    expect(constraint.length).toBeUndefined();
});
