import "reflect-metadata";
import { indent, safeName } from "../common/format";
import { Location } from "../common/location";
import { Annotation, Tag } from "./attribute";
import {
    IElement,
    IElementParams,
    IRenderable,
    setParentAndModelDeep,
} from "./common";
import { AppName, Model } from "./model";
import { addTags, renderAnnos } from "./renderers";

export type TypeConstraintRangeParams = {
    min?: number | undefined;
    max?: number | undefined;
};

export class TypeConstraintRange {
    min: number | undefined;
    max: number | undefined;

    constructor({ min, max }: TypeConstraintRangeParams) {
        this.min = min;
        this.max = max;
    }
}

export type TypeConstraintLengthParams = {
    min?: number | undefined;
    max?: number | undefined;
};

export class TypeConstraintLength {
    min: number | undefined;
    max: number | undefined;

    constructor({ min, max }: TypeConstraintLengthParams) {
        this.min = min;
        this.max = max;
    }
}

export type TypeConstraintResolutionParams = {
    base?: number | undefined;
    index?: number | undefined;
};

/** e.g.: 3 decimal places = {base = 10, index = -3} */
export class TypeConstraintResolution {
    base: number | undefined;
    index: number | undefined;

    constructor({ base, index }: TypeConstraintResolutionParams) {
        this.base = base;
        this.index = index;
    }
}

export type TypeConstraintParams = {
    range?: TypeConstraintRange;
    length?: TypeConstraintLength;
    resolution?: TypeConstraintResolution;
    precision?: number;
    scale?: number;
    bitWidth?: number;
};

export class TypeConstraint {
    range: TypeConstraintRange | undefined;
    length: TypeConstraintLength | undefined;
    resolution: TypeConstraintResolution | undefined;
    precision?: number;
    scale?: number;
    bitWidth?: number;

    constructor({
        range,
        length,
        resolution,
        precision,
        scale,
        bitWidth,
    }: TypeConstraintParams) {
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

    constructor(primitive: TypePrimitive, constraints?: TypeConstraint[]) {
        this.primitive = primitive;
        this.constraints = constraints ?? [];
    }

    private constraintStr(): string {
        const constraint = this.constraints?.length
            ? this.constraints[0]
            : null;
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

export class Enum implements IRenderable {
    members: EnumValue[];

    constructor(members: EnumValue[]) {
        this.members = members ?? [];
    }

    toSysl(): string {
        return `${indent(this.members.map(e => e.toSysl()).join("\n"))}`;
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
    type: TypeValue;

    constructor(type: TypeValue) {
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
        return `${indent(this.types.map(o => o.toSysl()).join("\n"))}`;
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
                .map(e => {
                    let sysl = addTags(
                        `${safeName(e.name)} <: ${e.value.toSysl()}${
                            e.optional ? "?" : ""
                        }`,
                        e.tags
                    );
                    if (e.annos.length) {
                        sysl += `:\n${indent(renderAnnos(e.annos))}`;
                    }
                    return sysl;
                })
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
    List,
}

export type ReferenceParams = {
    appName: AppName;
    typeName: string;
    fieldName?: string;
};

export class Reference implements IRenderable {
    /** The name of the referenced application. */
    appName?: AppName;
    /** The name of the referenced type. */
    typeName: string;
    /** The name of the referenced field (foreign key only). */
    fieldName?: string;

    constructor({ appName, typeName, fieldName }: ReferenceParams) {
        this.appName = appName;
        this.typeName = typeName;
        this.fieldName = fieldName;
    }

    toSysl(): string {
        let sysl = safeName(this.typeName);
        if (this.appName) {
            sysl = `${this.appName?.toSysl()}.${sysl}`;
        }
        if (this.fieldName) {
            sysl = `${sysl}.${safeName(this.fieldName)}`;
        }
        return sysl;
    }
}

export type TypeValue =
    | Primitive
    | Struct
    | TypeDecorator<IRenderable>
    | Alias
    | Enum
    | Union
    | Reference;

export type TypeParams = IElementParams & {
    discriminator?: string;
    name: string;
    value: TypeValue;
    optional?: boolean;
};

export class Type implements IElement {
    discriminator?: string;
    name: string;
    optional: boolean;
    value: TypeValue;
    annos: Annotation[];
    tags: Tag[];
    locations: Location[];
    parent?: IElement;
    model?: Model;

    constructor({
        discriminator,
        name,
        optional: optional,
        value,
        annos,
        tags,
        locations,
        parent,
        model,
    }: TypeParams) {
        this.discriminator = discriminator;
        this.name = name;
        this.optional = optional ?? false;
        this.value = value;
        this.annos = annos ?? [];
        this.tags = tags ?? [];
        this.locations = locations ?? [];
        this.parent = parent;
        this.model = model;

        setParentAndModelDeep(this, this.children(), this.annos, this.tags);
    }

    /** Returns an array of child types (i.e. fields) nested in this statement's {@code value}. */
    children(): Type[] {
        if (this.value && "elements" in this.value) {
            return this.value.elements;
        }
        return [];
    }

    toSysl(): string {
        const optStr = this.optional ? "?" : "";
        // Definition rendering
        if (this.discriminator) {
            let sysl = `${addTags(
                `${this.discriminator} ${safeName(this.name)}`,
                this.tags
            )}:`;
            if (this.annos.length) {
                sysl += `\n${indent(renderAnnos(this.annos))}`;
            }
            return (sysl += `\n${this.value.toSysl()}${optStr}`);
        }
        // Field rendering
        else {
            let sysl = `${safeName(this.name)}${this.value.toSysl()}${optStr}`;
            sysl = addTags(sysl, this.tags);
            if (this.annos.length) {
                sysl += `:\n${renderAnnos(this.annos)}`;
            }
            return sysl;
        }
    }
}
