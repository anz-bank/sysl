import { IRenderable } from "./common";
import { TypeConstraint, TypeConstraintLength } from "./constraint";

export class Primitive implements IRenderable {
    primitive: TypePrimitive;
    constraints: TypeConstraint[];

    constructor(primitive: TypePrimitive, constraints?: TypeConstraint[]) {
        this.primitive = primitive;
        this.constraints = constraints ?? [];
    }

    private constraintStr(): string {
        const constraint = this.constraints?.length ? this.constraints[0] : null;
        const isNumber = (n?: number) => n != null && !isNaN(n);

        const lengthStr = (length: TypeConstraintLength) => {
            if (isNumber(length.max) && isNumber(length.min)) {
                return `(${length.min}..${length.max})`;
            } else if (isNumber(length.max)) {
                return `(${length.max})`;
            } else if (isNumber(length.min)) {
                return `(${length.min}..)`;
            }
            return "";
        };

        if (constraint) {
            if (isNumber(constraint.precision) && isNumber(constraint.scale)) {
                return `(${constraint.precision}.${constraint.scale})`;
            }
            if (constraint.length) {
                return lengthStr(constraint.length);
            }
            if (constraint.bitWidth) {
                return constraint.bitWidth.toString();
            }
        }

        return "";
    }

    toSysl(): string {
        return `${this.primitive.toLowerCase()}${this.constraintStr()}`;
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
