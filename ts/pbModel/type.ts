import "reflect-metadata";
import {
    jsonArrayMember,
    jsonMapMember,
    jsonMember,
    jsonObject,
    jsonSetMember,
} from "typedjson";
import { Location } from "../location";
import {
    Alias,
    DecoratorKind,
    Enum,
    EnumValue,
    Primitive,
    Reference,
    Struct,
    Type,
    TypeConstraint,
    TypeConstraintLength,
    TypeConstraintRange,
    TypeConstraintResolution,
    TypeDecorator,
    TypePrimitive,
    TypeValue,
    Union,
} from "../model/type";
import { getAnnos, getTags, sortLocationalArray } from "../sort";
import { PbAttribute } from "./attribute";
import { PbAppName } from "./appname";
import { deserializeMap, serializeMap, serializerFor } from "./serialize";
import { joinedAppName } from "../format";
import { IRenderable } from "../model/common";

@jsonObject
export class PbValue {
    @jsonMember b?: boolean;
    @jsonMember d?: number;
    @jsonMember s?: string;

    toModel(): TypeValue {
        if (this.b) return this.b;
        else if (this.d) return this.d;
        else if (this.s) return this.s;
        else throw new Error("Missing Type value");
    }
}

@jsonObject
export class PbTypeConstraintRange {
    @jsonMember min?: PbValue;
    @jsonMember max?: PbValue;

    toModel(): TypeConstraintRange {
        return new TypeConstraintRange(
            this.min?.toModel(),
            this.max?.toModel()
        );
    }
}

@jsonObject
export class PbTypeConstraintLength {
    @jsonMember({
        deserializer: (numberString): number => Number(numberString),
    })
    min?: number;
    @jsonMember({
        deserializer: (numberString): number => Number(numberString),
    })
    max?: number;

    toModel(): TypeConstraintLength {
        return new TypeConstraintLength(this.min, this.max);
    }
}

/** e.g.: 3 decimal places = {base = 10, index = -3} */
@jsonObject
export class PbTypeConstraintResolution {
    @jsonMember base?: number;
    @jsonMember index?: number;

    toModel(): TypeConstraintResolution {
        return new TypeConstraintResolution(
            this.base ?? undefined,
            this.index ?? undefined
        );
    }
}

@jsonObject
export class PbTypeConstraint {
    @jsonMember range?: PbTypeConstraintRange;
    @jsonMember length?: PbTypeConstraintLength;
    @jsonMember resolution?: PbTypeConstraintResolution;
    @jsonMember precision!: number;
    @jsonMember scale!: number;
    @jsonMember bitWidth!: number;

    toModel(): TypeConstraint {
        return new TypeConstraint(
            this.range?.toModel(),
            this.length?.toModel(),
            this.resolution?.toModel(),
            this.precision,
            this.scale,
            this.bitWidth
        );
    }
}

@jsonObject
export class PbScope {
    @jsonArrayMember(String) path!: string[];
    @jsonMember appname!: PbAppName;

    toModel(): Reference {
        const name = `${
            this.appname?.part.any()
                ? joinedAppName(this.appname.part, false) + "."
                : ""
        }${this.path.join(".")}`;
        return new Reference(name);
    }
}

@jsonObject
export class PbScopedRef {
    /** The context in which the ref appeared. */
    @jsonMember context!: PbScope;
    @jsonMember ref!: PbScope;
}

@jsonObject
export class PbRelationKey {
    @jsonSetMember(String) attrName!: Set<string>;
}

@jsonObject
export class PbTypeRelation {
    @jsonMember primaryKey!: PbRelationKey;
    @jsonArrayMember(PbRelationKey) key!: PbRelationKey[];
    @jsonArrayMember(String) inject!: string[];
    @jsonMapMember(String, () => PbTypeDef, {
        serializer: (mapObject: Map<string, PbTypeDef>) =>
            serializeMap(mapObject),
        deserializer: (stringifiedMapObject: { [key: string]: PbTypeDef }) =>
            deserializeMap(PbTypeDef, stringifiedMapObject),
    })
    attrDefs!: Map<string, PbTypeDef>;

    toModel(): Struct {
        return new Struct(
            sortLocationalArray(
                this.attrDefs.select(a => a[1].toModel(a[0])).toArray()
            )
        );
    }
}

@jsonObject
export class PbTypeDefList {
    @jsonArrayMember(() => PbTypeDef) type!: PbTypeDef[];

    toModel(): Union {
        return new Union(this.type.select(t => t.toModel()).toArray());
    }
}

@jsonObject
export class PbTypeDefEnum {
    @jsonMapMember(String, Number, serializerFor(String))
    items!: Map<string, number>;

    toModel(): Enum {
        return new Enum(
            this.items.select(i => new EnumValue(i[0], i[1])).toArray()
        );
    }
}

@jsonObject
export class PbTypeDefTuple {
    @jsonMapMember(String, () => PbTypeDef, {
        serializer: (mapObject: Map<string, PbTypeDef>) =>
            serializeMap(mapObject),
        deserializer: (stringifiedMapObject: { [key: string]: PbTypeDef }) =>
            deserializeMap(PbTypeDef, stringifiedMapObject),
    })
    attrDefs!: Map<string, PbTypeDef>;

    toModel(): Struct {
        return new Struct(
            sortLocationalArray(
                this.attrDefs.select(a => a[1].toModel(a[0])).toArray()
            )
        );
    }
}

@jsonObject
export class PbTypeDef {
    @jsonMember primitive?: TypePrimitive;
    @jsonMember relation?: PbTypeRelation;
    @jsonMember typeRef?: PbScopedRef;
    @jsonMember set?: PbTypeDef;
    @jsonMember sequence?: PbTypeDef;
    @jsonArrayMember(PbTypeConstraint) constraint?: PbTypeConstraint[];
    @jsonMember opt!: boolean;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMember enum?: PbTypeDefEnum;
    @jsonMember tuple?: PbTypeDefTuple;
    @jsonMember list?: PbTypeDefList;
    @jsonMapMember(PbTypeDef, PbTypeDef, serializerFor(PbTypeDef)) map?: Map<
        PbTypeDef,
        PbTypeDef
    >;
    @jsonMember oneOf?: PbTypeDefList;
    @jsonMember noType?: Object;
    @jsonMapMember(String, PbAttribute, serializerFor(PbAttribute)) attrs!: Map<
        string,
        PbAttribute
    >;

    // `isInner` specifies whether a type exists within something else and is not a type definition.
    // It is true by default, meaning the type is a nested definition or a parameter.
    // When it is false, it means it is a top level `Type` definition and therefore may be an `Alias`
    toModel(name?: string | undefined, isInner: boolean = true): Type {
        let value:
            | Primitive
            | Struct
            | TypeDecorator<IRenderable>
            | Alias
            | Enum
            | Union
            | Reference;
        let discriminator = "";
        if (this.primitive) {
            value = new Primitive(
                this.primitive,
                this.constraint?.select(c => c.toModel()).toArray() ?? []
            );
        } else if (this.relation) {
            discriminator = "!table";
            value = this.relation.toModel();
        } else if (this.typeRef) {
            value = new TypeDecorator<Reference>(
                this.typeRef.ref.toModel(),
                DecoratorKind.Reference
            );
        } else if (this.set) {
            value = new TypeDecorator<Type>(
                this.set.toModel(),
                DecoratorKind.Set
            );
        } else if (this.sequence) {
            value = new TypeDecorator<Type>(
                this.sequence.toModel(),
                DecoratorKind.Sequence
            );
        } else if (this.enum) {
            discriminator = "!enum";
            value = this.enum.toModel();
        } else if (this.tuple) {
            discriminator = "!type";
            value = this.tuple.toModel();
        } else if (this.oneOf) {
            discriminator = "!union";
            value = this.oneOf.toModel();
        } else {
            throw new Error("Error converting type.");
        }
        // Catch the case that this is a top level definition and we failed to assign a discriminator above.
        if (!isInner && discriminator == "") {
            discriminator = "!alias";
            value = new Alias(value);
        }
        return new Type(
            discriminator,
            name ?? "",
            this.opt,
            value,
            this.sourceContexts,
            getTags(this.attrs),
            getAnnos(this.attrs)
        );
    }
}
