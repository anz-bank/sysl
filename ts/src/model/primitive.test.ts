import "jest-extended";
import { Primitive, TypePrimitive } from "./primitive";
import { Range, TypeConstraint } from "./constraint";

test.concurrent("Parsing", async () => {
    expect(Primitive.fromParts("float")).toEqual(new Primitive(TypePrimitive.FLOAT));
    expect(Primitive.fromParts("string")).toEqual(new Primitive(TypePrimitive.STRING));
    expect(Primitive.fromParts("decimal", "(8.5)")).toEqual(
        new Primitive(TypePrimitive.DECIMAL, new TypeConstraint(undefined, 8, 5))
    );

    expect(() => Primitive.fromParts("unknown")).toThrow();
});
