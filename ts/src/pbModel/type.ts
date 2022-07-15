import "reflect-metadata";
import {
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
import { PbAppName } from "./appname";
import { PbAttribute } from "./attribute";
import { serializerFor, serializerForFn } from "./serialize";

@jsonObject
export class PbValue {
    @jsonMember b?: boolean;
    @jsonMember d?: number;
    @jsonMember s?: string;

    toModel(): ValueType {
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
        return new TypeConstraintRange({
            min: Number(this.min?.toModel()),
            max: Number(this.max?.toModel()),
        });
    }
}

@jsonObject
export class PbTypeConstraintLength {
    @jsonMember({ deserializer: (numStr): number => Number(numStr) })
    min?: number;
    @jsonMember({ deserializer: (numStr): number => Number(numStr) })
    max?: number;

    toModel(): TypeConstraintLength {
        return new TypeConstraintLength({ min: this.min, max: this.max });
    }
}

/** e.g.: 3 decimal places = {base = 10, index = -3} */
@jsonObject
export class PbTypeConstraintResolution {
    @jsonMember base?: number;
    @jsonMember index?: number;

    toModel(): TypeConstraintResolution {
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
export class PbTypeDefList {
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

    hasValue(): boolean {
        return !!(
            this.primitive ||
            this.relation ||
            this.typeRef ||
            this.set ||
            this.sequence ||
            this.enum ||
            this.tuple ||
            this.list?.type?.length ||
            this.map?.size
        );
    }

    static defsToModel(attrDefs: Map<string, PbTypeDef>): Type[] {
        return sortLocationalArray(
            Array.from(attrDefs).map(([key, value]) => value.toModel(key))
        );
    }

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
                "Error converting type: " + name + " " + JSON.stringify(this)
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
            opt: this.opt,
            value,
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        });
    }
}