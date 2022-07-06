import "reflect-metadata";
import { indent } from "../format";
import { Location } from "../location";
import { Annotation, Tag } from "./attribute";
import { IDescribable, ILocational, IRenderable } from "./common";
import { renderAnnos, addTags } from "./renderers";

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

    constructor(
        range: TypeConstraintRange | undefined,
        length: TypeConstraintLength | undefined,
        resolution: TypeConstraintResolution | undefined,
        precision: number,
        scale: number,
        bitWidth: number
    ) {
        this.range = range;
        this.length = length;
        this.resolution = resolution;
        this.precision = precision;
        this.scale = scale;
        this.bitWidth = bitWidth;
    }
}

export class Primitive implements IRenderable {
    primitive: TypePrimitive;
    constraints: TypeConstraint[];

    constructor(primitive: TypePrimitive, constraints: TypeConstraint[]) {
        this.primitive = primitive;
        this.constraints = constraints ?? [];
    }

    private constraintStr(): string {
        const constraint = this.constraints?.any() ? this.constraints[0] : null;
        const isPresentAndNumber = (n: number | undefined) => n && !isNaN(n);

        const lengthStr = (length: TypeConstraintLength) => {
            if (
                isPresentAndNumber(length.max) &&
                isPresentAndNumber(length.min)
            ) {
                return `(${length.min}..${length.max})`;
            } else if (isPresentAndNumber(length.max)) {
                return `(${length.max})`;
            } else if (isPresentAndNumber(length.min)) {
                return `(${length.min}..)`;
            }
            return "";
        };

        if (constraint) {
            if (
                isPresentAndNumber(constraint.precision) &&
                isPresentAndNumber(constraint.scale)
            )
                return `(${constraint.precision}.${constraint.scale})`;
            if (constraint.length) return lengthStr(constraint.length);
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

export class Enum implements IRenderable {
    members: EnumValue[];

    constructor(members: EnumValue[]) {
        this.members = members ?? [];
    }

    toSysl(): string {
        return `${indent(
            this.members
                .select(e => e.toSysl())
                .toArray()
                .join("\n")
        )}`;
    }
}

export class EnumValue implements IRenderable {
    name: string;
    value: number;

    constructor(name: string, value: number) {
        this.name = name;
        this.value = value;
    }

    toSysl(): string {
        return `${this.name}: ${this.value}`;
    }
}

export class Alias implements IRenderable {
    type:
        | Primitive
        | Struct
        | TypeDecorator<IRenderable>
        | Alias
        | Enum
        | Union
        | Reference;

    constructor(
        type:
            | Primitive
            | Struct
            | TypeDecorator<IRenderable>
            | Alias
            | Enum
            | Union
            | Reference
    ) {
        this.type = type;
    }

    toSysl(): string {
        return `${indent(this.type.toSysl())}`;
    }
}

export class Union implements IRenderable {
    types: Type[];

    constructor(types: Type[]) {
        this.types = types ?? [];
    }

    toSysl(): string {
        return `${indent(
            this.types
                .select(o => o.toSysl())
                .toArray()
                .join("\n")
        )}`;
    }
}

// The type of `!type`, `!table` and endpoints.
export class Struct implements IRenderable {
    elements: Type[];

    constructor(elements: Type[]) {
        this.elements = elements ?? [];
    }

    toSysl(): string {
        return indent(
            this.elements
                .select(e => {
                    let sysl = addTags(
                        `${e.name} <: ${e.value.toSysl()}${e.opt ? "?" : ""}`,
                        e.tags
                    );
                    if (e.annos.any()) {
                        sysl += `:\n${indent(renderAnnos(e.annos))}`;
                    }
                    return sysl;
                })
                .toArray()
                .join("\n")
        );
    }
}

// Used for `sequence of`, `set of` or references to types in other places in the model
export class TypeDecorator<T extends IRenderable> implements IRenderable {
    innerType: T;
    kind: DecoratorKind;

    constructor(innerType: T, kind: DecoratorKind) {
        this.innerType = innerType;
        this.kind = kind;
    }

    toSysl(): string {
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
    Reference,
}

export class Reference implements IRenderable {
    name: string;

    constructor(name: string) {
        this.name = name;
    }

    toSysl(): string {
        return this.name;
    }
}

export class Type implements IDescribable, ILocational, IRenderable {
    discriminator: string;
    name: string;
    opt: boolean;
    value:
        | Primitive
        | Struct
        | TypeDecorator<IRenderable>
        | Alias
        | Enum
        | Union
        | Reference;
    locations: Location[];
    tags: Tag[];
    annos: Annotation[];

    constructor(
        discriminator: string,
        name: string,
        opt: boolean,
        value:
            | Primitive
            | Struct
            | TypeDecorator<IRenderable>
            | Alias
            | Enum
            | Union
            | Reference,
        locations: Location[],
        tags: Tag[],
        annos: Annotation[]
    ) {
        this.discriminator = discriminator;
        this.name = name;
        this.opt = opt;
        this.value = value;
        this.locations = locations;
        this.tags = tags;
        this.annos = annos;
    }

    toSysl(): string {
        // Definition rendering
        if (this.discriminator.length > 0) {
            let sysl = `${addTags(
                `${this.discriminator} ${this.name}`,
                this.tags
            )}:`;
            if (this.annos.any()) {
                sysl += `\n${indent(renderAnnos(this.annos))}`;
            }
            return (sysl += `\n${this.value.toSysl()}${this.opt ? "?" : ""}`);
        }
        // Field rendering
        else {
            let sysl = addTags(
                `${this.name}${this.value.toSysl()}${this.opt ? "?" : ""}`,
                this.tags
            );
            if (this.annos.any()) {
                sysl += `:\n${renderAnnos(this.annos)}`;
            }
            return sysl;
        }
    }
}
