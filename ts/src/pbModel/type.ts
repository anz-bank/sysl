import "reflect-metadata";
import { jsonArrayMember, jsonMapMember, jsonMember, jsonObject, jsonSetMember } from "typedjson";
import { Location } from "../common/location";
import { ElementRef, Type, TypeConstraint, Range, Union, ValueType, Alias, AppChild } from "../model";
import { CollectionDecorator } from "../model/decorator";
import { Enum, EnumValue } from "../model/enum";
import { Field } from "../model/field";
import { FieldValue } from "../model/fieldValue";
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

@jsonObject
export class PbTypeConstraint {
    @jsonMember range?: PbTypeConstraintRange;
    @jsonMember length?: PbTypeConstraintLength;
    @jsonMember precision?: number;
    @jsonMember scale?: number;
    @jsonMember bitWidth?: number;

    toModel(): TypeConstraint {
        return new TypeConstraint(
            this.length?.toModel(),
            this.precision,
            this.scale,
            this.bitWidth as 32 | 64 | undefined
        );
    }
}

@jsonObject
export class PbScope {
    @jsonArrayMember(String) path!: string[];
    @jsonMember appname!: PbAppName;

    toModel(parentRef?: ElementRef): ElementRef {
        const [typeName, fieldName] = this.path;
        let namespace = this.appname?.part.slice(0, -1) ?? [];
        let appName = this.appname?.part.at(-1) ?? "";

        if (!appName && parentRef?.appName) {
            appName = parentRef?.appName;
            namespace = [...parentRef.namespace];
        }

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

    getContext(): PbScope | undefined {
        return (
            this.typeRef?.context ||
            this.set?.typeRef?.context ||
            this.sequence?.typeRef?.context ||
            this.list?.type.getContext()
        );
    }

    static defsToFields(attrDefs: Map<string, PbTypeDef>, parentRef?: ElementRef): Field[] {
        return Array.from(attrDefs)
            .map(([key, value]) => value?.toField(value.getContext()?.toModel() ?? parentRef, key))
            .sort(Location.compareFirst);
    }

    toValue(parentRef: ElementRef | undefined): FieldValue {
        const collection = this.set || this.sequence || this.list;
        if (collection) return new CollectionDecorator(collection.toValue(parentRef), !!this.set);
        if (this.primitive) return new Primitive(this.primitive, this.constraint?.at(0)?.toModel());
        if (this.typeRef) return this.typeRef.ref.toModel(parentRef);
        if (this.noType) return Primitive.Any;
        throw new Error(`Error converting type: ${JSON.stringify(this)}`);
    }

    toAppChild(name: string = "", parentRef?: ElementRef): AppChild {
        const params = {
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
            locations: this.sourceContexts,
        };

        const type = this.tuple || this.relation;
        if (type) return new Type(name, !!this.relation, PbTypeDef.defsToFields(type.attrDefs, parentRef), params);
        if (this.enum) {
            // Original order of enum items is not serialized, so assume value (number) order, since it's most common.
            const values = [...this.enum.items].map(([k, v]) => new EnumValue(k, v)).sort((a, b) => a.value - b.value);
            return new Enum(name, values, params);
        }
        if (this.oneOf) {
            const values = this.oneOf.type.map(t => t.toValue(parentRef));
            return new Union(name, values, params);
        }
        return new Alias(name, this.toValue(parentRef), params);
    }

    toField(parentRef: ElementRef | undefined, name: string): Field {
        if (this.tuple || this.relation || this.enum || this.oneOf)
            throw new Error("Cannot produce Field, requested PbTypeDef is not a Field");

        const params = {
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
            locations: this.sourceContexts,
        };

        return new Field(name, this.toValue(parentRef), this.opt, params);
    }
}
