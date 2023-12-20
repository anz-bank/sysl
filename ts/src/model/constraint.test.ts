import { realign } from "../common/format";
import { Range, TypeConstraint } from "./constraint";
import { Model } from "./model";
import { Primitive, TypePrimitive } from "./primitive";
import "jest-extended";

test.concurrent("Constructor", () => {
    // Incompatible constraints
    expect(() => new TypeConstraint(undefined, 10, 2, 64)).toThrow();
    expect(() => new TypeConstraint(new Range(undefined, 10), 10, 2, 64)).toThrow();

    // Scale/precision removes length constraint (TODO: Fix bug in Sysl binary)
    expect(new TypeConstraint(new Range(undefined, 10), 10, 2)).toEqual(new TypeConstraint(undefined, 10, 2));

    // Invalid values
    expect(() => new TypeConstraint(new Range())).toThrow();
    expect(() => new TypeConstraint(new Range(10, 5))).toThrow();
    expect(() => new TypeConstraint(new Range(NaN, 5))).toThrow();
    expect(() => new TypeConstraint(new Range(10, NaN))).toThrow();
    expect(() => new TypeConstraint(undefined, 5)).toThrow();
    expect(() => new TypeConstraint(undefined, undefined, 10)).toThrow();
    expect(() => new TypeConstraint(undefined, 0, 0)).toThrow();
    expect(() => new TypeConstraint(undefined, 5, 10)).toThrow();
    expect(() => new TypeConstraint(undefined, 10, -5)).toThrow();
    expect(() => new TypeConstraint(undefined, -5, -10)).toThrow();
    expect(() => new TypeConstraint(undefined, NaN, 5)).toThrow();
    expect(() => new TypeConstraint(undefined, 10, NaN)).toThrow();
    expect(() => new TypeConstraint(undefined, undefined, undefined, 24 as 32 /* lol */)).toThrow();
});

test.concurrent("Parse", () => {
    expect(TypeConstraint.parse("(8.5)")).toEqual(new TypeConstraint(undefined, 8, 5));
    expect(TypeConstraint.parse("(30)")).toEqual(new TypeConstraint(new Range(undefined, 30)));
    expect(TypeConstraint.parse("(..30)")).toEqual(new TypeConstraint(new Range(undefined, 30)));
    expect(TypeConstraint.parse("(30..)")).toEqual(new TypeConstraint(new Range(30)));
    expect(TypeConstraint.parse("(3..30)")).toEqual(new TypeConstraint(new Range(3, 30)));

    expect(TypeConstraint.parse("64(30)")).toEqual(new TypeConstraint(new Range(undefined, 30), undefined, undefined, 64));
    expect(TypeConstraint.parse("64(..30)")).toEqual(new TypeConstraint(new Range(undefined, 30), undefined, undefined, 64));
    expect(TypeConstraint.parse("64(30..)")).toEqual(new TypeConstraint(new Range(30), undefined, undefined, 64));
    expect(TypeConstraint.parse("64(3..30)")).toEqual(new TypeConstraint(new Range(3, 30), undefined, undefined, 64));

    expect(() => TypeConstraint.parse("(")).toThrow();
    expect(() => TypeConstraint.parse(")")).toThrow();
    expect(() => TypeConstraint.parse("()")).toThrow();
    expect(() => TypeConstraint.parse("(..)")).toThrow();
    expect(() => TypeConstraint.parse("(3..2)")).toThrow();
    expect(() => TypeConstraint.parse("(1..2..3)")).toThrow();
    expect(() => TypeConstraint.parse("(.)")).toThrow();
    expect(() => TypeConstraint.parse("(0.0)")).toThrow();
    expect(() => TypeConstraint.parse("(5.8")).toThrow();
    expect(() => TypeConstraint.parse("(1.2.3)")).toThrow();
    expect(() => TypeConstraint.parse("(1.2.3.4)")).toThrow();
    expect(() => TypeConstraint.parse("((10)")).toThrow();
    expect(() => TypeConstraint.parse("(10))")).toThrow();
    expect(() => TypeConstraint.parse("((10))")).toThrow();

});

test.concurrent("Decimal with precision and scale", async () => {
    const model = await Model.fromText(realign(`
        App1:
            !type Type1:
                Field1 <: decimal(5.2)
    `));

    const primitive = model.getField("App1.Type1.Field1").value as Primitive;
    expect(primitive).toBeInstanceOf(Primitive);
    expect(primitive.primitive).toBe(TypePrimitive.DECIMAL);

    const constraint = primitive.constraint!;
    expect(constraint).toBeInstanceOf(TypeConstraint);

    expect(constraint.precision).toBe(5);
    expect(constraint.scale).toBe(2);

    expect(constraint.bitWidth).toBeUndefined();
    expect(constraint.length).toBeUndefined();
});

test.concurrent("Integer with precision", async () => {
    const model = await Model.fromText(realign(`
        App1:
            !type Type1:
                Field1 <: int64(5..10)
    `));

    const primitive = model.getField("App1.Type1.Field1").value as Primitive;
    expect(primitive).toBeInstanceOf(Primitive);
    expect(primitive.primitive).toBe(TypePrimitive.INT);

    const constraint = primitive.constraint!;
    expect(constraint).toBeInstanceOf(TypeConstraint);

    expect(constraint.precision).toBeUndefined();
    expect(constraint.scale).toBeUndefined();

    expect(constraint.bitWidth).toBe(64);
    expect(constraint.length).toEqual(new Range(5, 10));
});
