import "reflect-metadata";
import { Location } from "../location";
import { indent } from "../util";
import { BaseType, ComplexType, SimpleType } from "./common";
import { Element } from "./element";

export type TypeValue = boolean | number | string;

export class TypeConstraintRange {
    min: TypeValue | undefined;
    max: TypeValue | undefined;

    constructor(min: TypeValue | undefined, max: TypeValue | undefined) {
        this.min = min;
        this.max = max;
    }

}

export class TypeConstraintLength {
    min: number | undefined;
    max: number | undefined;

    constructor(min: number | undefined, max: number | undefined) {
        this.min = min;
        this.max = max;
    }
}

/** e.g.: 3 decimal places = {base = 10, index = -3} */
export class TypeConstraintResolution {
    base: number | undefined;
    index: number | undefined;

    constructor(base: number | undefined, index: number | undefined) {
        this.base = base;
        this.index = index;
    }
}

export class TypeConstraint {
    range: TypeConstraintRange | undefined;
    length: TypeConstraintLength | undefined;
    resolution: TypeConstraintResolution | undefined;
    precision: number;
    scale: number;
    bitWidth: number;

    constructor(range: TypeConstraintRange | undefined,
        length: TypeConstraintLength | undefined,
        resolution: TypeConstraintResolution | undefined,
        precision: number,
        scale: number,
        bitWidth: number) {
        this.range = range;
        this.length = length;
        this.resolution = resolution;
        this.precision = precision;
        this.scale = scale;
        this.bitWidth = bitWidth;
    }
}

export class Primitive extends SimpleType {
    primitive: TypePrimitive;
    constraints: TypeConstraint[];

    constructor(primitive: TypePrimitive, constraints: TypeConstraint[]) {
        super();
        this.primitive = primitive;
        this.constraints = constraints ?? [];
    }

    private constraintStr(): string {
        const constraint = this.constraints?.any() ? this.constraints[0] : null;
        const isPresentAndNumber = (n: number | undefined) => n && !isNaN(n);

        const lengthStr = (length: TypeConstraintLength) => {
            if (isPresentAndNumber(length.max) && isPresentAndNumber(length.min)) {
                return `(${length.min}..${length.max})`
            }
            else if (isPresentAndNumber(length.max)) {
                return `(${length.max})`
            }
            else if (isPresentAndNumber(length.min)) {
                return `(${length.min}..)`
            }
            return '';
        }

        if (constraint) {
            if (isPresentAndNumber(constraint.precision) && isPresentAndNumber(constraint.scale)) return `(${constraint.precision}.${constraint.scale})`
            if (constraint.length) return lengthStr(constraint.length);
        }

        return '';
    }

    override toSysl(): string {
        return `${this.primitive.toLowerCase()}${this.constraintStr()}`
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

export class Enum extends SimpleType {
    members: EnumValue[]

    constructor(members: EnumValue[]) {
        super();
        this.members = members ?? [];
    }

    override toSysl(): string {
        return `${indent(this.members.select(e => e.toSysl()).toArray().join("\n"))}`
    }
}

export class EnumValue {
    name: string;
    value: number;

    constructor(name: string, value: number) {
        this.name = name;
        this.value = value;
    }

    toSysl(): string {
        return `${this.name}: ${this.value}`
    }
}

export class Alias extends SimpleType {
    type: Type

    constructor(type: Type) {
        super();
        this.type = type;
    }

    override toSysl(): string {
        return `\n${indent(this.type.toSysl())}`
    }

}

export class Union extends SimpleType {
    types: Type[]

    constructor(types: Type[]) {
        super();
        this.types = types ?? [];
    }

    override toSysl(): string {
        return `${indent(this.types.select(o => o.toSysl()).toArray().join("\n"))}`;
    }
}

// The type of `!type`, `!table` and endpoints.
export class Struct extends SimpleType {
    elements: Element<ComplexType>[];

    constructor(elements: Element<ComplexType>[]) {
        super()
        this.elements = elements ?? [];
    }

    override toSysl(): string {
        return indent(this.elements.select(e => e.toSysl()).toArray().join("\n"));
    }
}

// Used for `sequence of`, `set of` or references to types in other places in the model
export class TypeDecorator<T extends BaseType> extends SimpleType {
    innerType: T;
    kind: DecoratorKind;

    constructor(innerType: T, kind: DecoratorKind) {
        super()
        this.innerType = innerType;
        this.kind = kind;
    }

    override toSysl(): string {
        switch (this.kind) {
            case DecoratorKind.Set:
                return `set of ${this.innerType.toSysl()}`;
            case DecoratorKind.Sequence:
                return `sequence of ${this.innerType.toSysl()}`;
            default:
                return `${this.innerType.toSysl()}`;
        }
    }
}

export enum DecoratorKind {
    Set,
    Sequence,
    Reference
}

export class Reference extends SimpleType {
    name: string;

    constructor(name: string) {
        super();
        this.name = name;
    }

    override toSysl(): string {
        return this.name;
    }
}

export class Type extends ComplexType {
    opt: boolean;
    value: Primitive | Struct | TypeDecorator<SimpleType> | Alias | Enum | Union | Reference;

    constructor(discriminator: string,
        name: string,
        opt: boolean,
        value: Primitive | Struct | TypeDecorator<SimpleType> | Alias | Enum | Union | Reference,
        locations: Location[],
    ) {
        super(discriminator, locations, name);
        this.opt = opt;
        this.value = value;
    }

    override toSysl(): string {
        return `${this.value.toSysl()}${this.opt ? '?' : ''}`;
    }
}
