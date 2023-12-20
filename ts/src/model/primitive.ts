import { IRenderable } from "./common";
import { TypeConstraint } from "./constraint";

export enum TypePrimitive {
    ANY = "ANY",
    BOOL = "BOOL",
    INT = "INT",
    FLOAT = "FLOAT",
    DECIMAL = "DECIMAL",
    /** STRING - Unicode string (Python 2 unicode, Python 3 str, and SQL nvarchar). */
    STRING = "STRING",
    /** BYTES - Octet sequence, like Python 3 bytes and SQL varbinary. */
    BYTES = "BYTES",
    DATE = "DATE",
    DATETIME = "DATETIME",
    UUID = "UUID",
}

export class Primitive implements IRenderable {
    public static readonly Any = new Primitive(TypePrimitive.ANY);

    constructor(public readonly primitive: TypePrimitive, public readonly constraint?: TypeConstraint) {
        // TODO: Check constraint matches TypePrimitive
        this.primitive = primitive;
        this.constraint = constraint;
    }

    toSysl(): string {
        return `${this.primitive.toLowerCase()}${this.constraint?.toString() ?? ""}`;
    }

    toString(): string {
        return this.primitive.toLowerCase();
    }

    private static readonly names = Object.values(TypePrimitive).map(v =>
        v.toString().toLowerCase()
    ) as readonly string[];

    static fromParts(name: string, constraintStr?: string): Primitive {
        if (!Primitive.names.includes(name)) throw new Error(`Unknown primitive: ${name}`);
        let constraint = constraintStr ? TypeConstraint.parse(constraintStr) : undefined;
        return new Primitive(TypePrimitive[name.toUpperCase() as keyof typeof TypePrimitive], constraint);
    }
}
