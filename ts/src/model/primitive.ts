import { IRenderable } from "./common";
import { CloneContext } from "./clone";
import { TypeConstraint, Range } from "./constraint";

export class Primitive implements IRenderable {
    constructor(public primitive: TypePrimitive, public constraint?: TypeConstraint) {
        this.primitive = primitive;
        this.constraint = constraint;
    }

    public constraintStr(): string {
        const isNumber = (n?: number) => n != null && !isNaN(n);

        const lengthStr = (length: Range) => {
            if (isNumber(length.max) && isNumber(length.min)) {
                return `(${length.min}..${length.max})`;
            } else if (isNumber(length.max)) {
                return `(${length.max})`;
            } else if (isNumber(length.min)) {
                return `(${length.min}..)`;
            }
            return "";
        };

        if (this.constraint) {
            if (isNumber(this.constraint.precision) && isNumber(this.constraint.scale)) {
                return `(${this.constraint.precision}.${this.constraint.scale})`;
            }
            if (this.constraint.length) {
                return lengthStr(this.constraint.length);
            }
            if (this.constraint.bitWidth) {
                return this.constraint.bitWidth.toString();
            }
        }

        return "";
    }

    toSysl(): string {
        return `${this.primitive.toLowerCase()}${this.constraintStr()}`;
    }

    toString(): string {
        return this.primitive.toLowerCase();
    }

    clone(context = new CloneContext()): Primitive {
        return new Primitive(this.primitive, context.applyUnder(this.constraint));
    }
}

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
