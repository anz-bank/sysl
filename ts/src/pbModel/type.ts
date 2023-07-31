import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject, jsonSetMember } from "typedjson";
import { Location } from "../common/location";
import { sortByLocation } from "../common/sort";
import { Element, ElementRef, Type, TypeConstraint, Range, DecimalResolution, Union, ValueType, Alias } from "../model";
import { CollectionDecorator } from "../model/decorator";
import { Enum, EnumValue } from "../model/enum";
import { Field, FieldValue } from "../model/field";
import { Primitive, TypePrimitive } from "../model/primitive";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor, serializerForFn } from "./serialize";
import { PbTypeDefList } from "./type_list";

function noNaN(n?: number): number | undefined {
    return isNaN(n as number) ? undefined : n;
}

@jsonObject
export class PbValue {
    @jsonMember b?: boolean;
    @jsonMember d?: number;
    @jsonMember s?: string;

    toModel(): ValueType | undefined {
        return this.b || this.d || this.s;
    }

    toNumber(): number | undefined {
        if (this.b || this.s) throw new Error("Requested number of a PbValue with string/bool");
        return noNaN(this.d);
    }
}

@jsonObject
export class PbTypeConstraintRange {
    @jsonMember min?: PbValue;
    @jsonMember max?: PbValue;

    toModel(): Range | undefined {
        return this.min || this.max ? new Range(this.min?.toNumber(), this.max?.toNumber()) : undefined;
    }
}

@jsonObject
export class PbTypeConstraintLength {
    @jsonMember({ deserializer: (numStr): number | undefined => noNaN(Number(numStr)) })
    min?: number;
    @jsonMember({ deserializer: (numStr): number | undefined => noNaN(Number(numStr)) })
    max?: number;

    toModel(): Range | undefined {
        return this.min || this.max ? new Range(this.min, this.max) : undefined;
    }
}

/** e.g.: 3 decimal places = {base = 10, index = -3} */
@jsonObject
export class PbTypeConstraintResolution {
    @jsonMember base?: number;
    @jsonMember index?: number;

    toModel(): DecimalResolution | undefined {
        return this.base || this.index ? new DecimalResolution(noNaN(this.base), noNaN(this.index)) : undefined;
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

    toModel(): ElementRef {
        const [typeName, fieldName] = this.path;
        const namespace = this.appname?.part.slice(0, -1) ?? [];
        const appName = this.appname?.part.at(-1) ?? "";

        return new ElementRef(namespace, appName, typeName, fieldName);
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
}

@jsonObject
export class PbTypeDefStruct {
    @jsonMapMember(String, () => PbTypeDef, serializerForFn(() => PbTypeDef))
    attrDefs!: Map<string, PbTypeDef>;
}

@jsonObject
export class PbTypeDefUnion {
    @jsonArrayMember(() => PbTypeDef) type: PbTypeDef[] = [];
}

@jsonObject
export class PbTypeDefEnum {
    @jsonMapMember(String, Number, serializerFor(String))
    items!: Map<string, number>;
}

@jsonObject
export class PbTypeDef {
    @jsonMember primitive?: TypePrimitive;
    @jsonMember relation?: PbTypeRelation;
    @jsonMember typeRef?: PbScopedRef;
    @jsonMember set?: PbTypeDef;
    @jsonMember sequence?: PbTypeDef;
    @jsonMember(() => PbTypeDefList) list?: PbTypeDefList;
    @jsonArrayMember(PbTypeConstraint) constraint?: PbTypeConstraint[];
    @jsonMember opt!: boolean;
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMember enum?: PbTypeDefEnum;
    @jsonMember tuple?: PbTypeDefStruct;
    // FIXME: Defaults to an empty map even if not defined in the deserialized source.
    @jsonMapMember(PbTypeDef, PbTypeDef, serializerFor(PbTypeDef))
    map?: Map<PbTypeDef, PbTypeDef>;
    @jsonMember oneOf?: PbTypeDefUnion;
    @jsonMember noType?: Object;
    @jsonMapMember(String, PbAttribute, serializerFor(PbAttribute)) attrs!: Map<string, PbAttribute>;

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

    static defsToFields(attrDefs: Map<string, PbTypeDef>): Field[] {
        return sortByLocation(Array.from(attrDefs).map(([key, value]) => value?.toField(key)));
    }

    toValue(): FieldValue {
        const collection = this.set || this.sequence || this.list;
        if (collection) return new CollectionDecorator(collection.toModel(), !!this.set);
        if (this.primitive) return new Primitive(this.primitive, this.constraint?.at(0)?.toModel());
        if (this.typeRef) return this.typeRef.ref.toModel();
        throw new Error(`Error converting type: ${JSON.stringify(this)}`);
    }

    // `isInner` specifies whether a type exists within something else and is not a type definition.
    // It is true by default, meaning the type is a nested definition or a parameter.
    // When it is false, it means it is a top level `Type` definition and therefore may be an `Alias`.
    toModel(name: string = "", isInner: boolean = true): Element {
        const params = {
            tags: sortByLocation(getTags(this.attrs)),
            annos: sortByLocation(getAnnos(this.attrs)),
            locations: this.sourceContexts,
        };

        const type = this.tuple || this.relation;
        if (type) return new Type(name, !!this.relation, sortByLocation(PbTypeDef.defsToFields(type.attrDefs)), params);
        if (this.enum) {
            // Original order of enum items is not serialized, so assume value (number) order, since it's most common.
            const values = [...this.enum.items].map(([k, v]) => new EnumValue(k, v)).sort((a, b) => a.value - b.value);
            return new Enum(name, values, params);
        }
        if (this.oneOf) {
            const values = this.oneOf.type.map(t => t.toValue());
            return new Union(name, values, params);
        }
        return isInner ? new Field(name, this.toValue(), this.opt, params) : new Alias(name, this.toValue(), params);
    }

    toField(name: string): Field {
        const field = this.toModel(name);
        if (!(field instanceof Field)) throw new Error("Cannot produce Field, requested PbTypeDef is not a Field");
        return field;
    }
}
