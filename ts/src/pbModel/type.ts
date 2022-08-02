import "reflect-metadata";
import {
    defaultTypeEmitter,
    jsonArrayMember,
    jsonMapMember,
    jsonMember,
    jsonObject,
    jsonSetMember,
} from "typedjson";
import { Location } from "../common/location";
import { getAnnos, getTags, sortLocationalArray } from "../common/sort";
import { AppName, ValueType } from "../model";
import { IRenderable } from "../model/common";
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
    Union,
} from "../model/type";
import type { PbTypeDefList } from "./type_list";
import { PbAppName } from "./appname";
import { PbAttribute } from "./attribute";
import { serializerFor, serializerForFn } from "./serialize";

@jsonObject
export class PbValue {
    @jsonMember b?: boolean;
    @jsonMember d?: number;
    @jsonMember s?: string;

    isEmpty(): boolean {
        return this.b == null && this.d == null && this.s == null;
    }

    toModel(): ValueType | undefined {
        if (this.b != null) {
            return this.b;
        } else if (this.d != null) {
            return this.d;
        } else if (this.s != null) {
            return this.s;
        } else {
            return undefined;
        }
    }
}

@jsonObject
export class PbTypeConstraintRange {
    @jsonMember min?: PbValue;
    @jsonMember max?: PbValue;

    isEmpty(): boolean {
        return this.min == null && this.max == null;
    }

    toModel(): TypeConstraintRange | undefined {
        if (this.isEmpty()) {
            return undefined;
        }
        const min = this.min?.toModel();
        const max = this.max?.toModel();
        return new TypeConstraintRange({
            min: min ? Number(min) : undefined,
            max: max ? Number(max) : undefined,
        });
    }
}

@jsonObject
export class PbTypeConstraintLength {
    @jsonMember({ deserializer: (numStr): number => Number(numStr) })
    min?: number;
    @jsonMember({ deserializer: (numStr): number => Number(numStr) })
    max?: number;

    isEmpty(): boolean {
        return (
            (this.min == null || isNaN(this.min)) &&
            (this.max == null || isNaN(this.max))
        );
    }

    toModel(): TypeConstraintLength | undefined {
        if (this.isEmpty()) {
            return undefined;
        }
        return new TypeConstraintLength({ min: this.min, max: this.max });
    }
}

/** e.g.: 3 decimal places = {base = 10, index = -3} */
@jsonObject
export class PbTypeConstraintResolution {
    @jsonMember base?: number;
    @jsonMember index?: number;

    isEmpty(): boolean {
        return (
            (this.base == null || isNaN(this.base)) &&
            (this.index == null || isNaN(this.index))
        );
    }

    toModel(): TypeConstraintResolution | undefined {
        if (this.isEmpty()) {
            return undefined;
        }
        return new TypeConstraintResolution({ ...this });
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
        return new TypeConstraint({
            range: this.range?.toModel(),
            length: this.length?.toModel(),
            resolution: this.resolution?.toModel(),
            precision: this.precision,
            scale: this.scale,
            bitWidth: this.bitWidth,
        });
    }
}

@jsonObject
export class PbScope {
    @jsonArrayMember(String) path!: string[];
    @jsonMember appname!: PbAppName;

    toModel(): Reference {
        const [typeName, fieldName] = this.path;
        return new Reference({
            appName: this.appname?.part && new AppName(this.appname.part),
            typeName,
            fieldName,
        });
    }
}

@jsonObject
export class PbScopedRef {
    /** The context in which the ref appeared. */
    @jsonMember context!: PbScope;
    /** The target of the ref. */
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
    @jsonMapMember(String, () => PbTypeDef, serializerForFn(() => PbTypeDef))
    attrDefs!: Map<string, PbTypeDef>;

    toModel(): Struct {
        return new Struct(PbTypeDef.defsToModel(this.attrDefs));
    }
}

@jsonObject
export class PbTypeDefTuple {
    @jsonMapMember(String, () => PbTypeDef, serializerForFn(() => PbTypeDef))
    attrDefs!: Map<string, PbTypeDef>;

    toModel(): Struct {
        return new Struct(PbTypeDef.defsToModel(this.attrDefs));
    }
}

@jsonObject
export class PbTypeDefUnion {
    @jsonArrayMember(() => PbTypeDef) type!: PbTypeDef[];

    toModel(): Union {
        return new Union(this.type.map(t => t.toModel()));
    }
}

@jsonObject
export class PbTypeDefEnum {
    @jsonMapMember(String, Number, serializerFor(String))
    items!: Map<string, number>;

    toModel(): Enum {
        return new Enum(
            Array.from(this.items).map(
                ([key, value]) => new EnumValue(key, value)
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
    @jsonMember list?: PbTypeDefList;
    @jsonArrayMember(PbTypeConstraint) constraint?: PbTypeConstraint[];
    @jsonMember opt!: boolean;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMember enum?: PbTypeDefEnum;
    @jsonMember tuple?: PbTypeDefTuple;
    // FIXME: Defaults to an empty map even if not defined in the deserialized source.
    @jsonMapMember(PbTypeDef, PbTypeDef, serializerFor(PbTypeDef))
    map?: Map<PbTypeDef, PbTypeDef>;
    @jsonMember oneOf?: PbTypeDefUnion;
    @jsonMember noType?: Object;
    @jsonMapMember(String, PbAttribute, serializerFor(PbAttribute)) attrs!: Map<
        string,
        PbAttribute
    >;

    hasValue(): boolean {
        return !!(
            this.primitive ||
            this.relation ||
            this.typeRef ||
            this.set ||
            this.sequence ||
            this.enum ||
            this.tuple ||
            this.list ||
            this.map?.size
        );
    }

    static defsToModel(attrDefs: Map<string, PbTypeDef>): Type[] {
        return sortLocationalArray(
            Array.from(attrDefs).map(([key, value]) => value?.toModel(key))
        );
    }

    // `isInner` specifies whether a type exists within something else and is not a type definition.
    // It is true by default, meaning the type is a nested definition or a parameter.
    // When it is false, it means it is a top level `Type` definition and therefore may be an
    // `Alias`.
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
                (this.constraint ?? []).map(c => c.toModel())
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
        } else if (this.list) {
            value = new TypeDecorator<Type>(
                this.list.toModel(),
                DecoratorKind.List
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
            throw new Error(
                `Error converting type: ${name} ${JSON.stringify(this)}`
            );
        }
        // Catch the case that this is a top level definition and we failed to assign a discriminator above.
        if (!isInner && discriminator == "") {
            discriminator = "!alias";
            value = new Alias(value);
        }
        return new Type({
            discriminator,
            name: name ?? "",
            optional: this.opt,
            value,
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        });
    }
}
