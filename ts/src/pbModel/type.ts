import "reflect-metadata";
import {
    jsonArrayMember,
    jsonMapMember,
    jsonMember,
    jsonObject,
    jsonSetMember,
} from "typedjson";
import { Location } from "../common/location";
import { sortLocationalArray } from "../common/sort";
import {
    Element,
    ElementRef,
    GenericElement,
    GenericValue,
    Struct,
    Type,
    TypeConstraint,
    TypeConstraintLength,
    TypeConstraintRange,
    TypeConstraintResolution,
    Union,
    ValueType,
} from "../model";
import { CollectionDecorator } from "../model/decorator";
import { Enum, EnumValue } from "../model/enum";
import { Field, FieldValue } from "../model/field";
import { Primitive, TypePrimitive } from "../model/primitive";
import { PbAppName } from "./appname";
import { getAnnos, getTags, PbAttribute } from "./attribute";
import { serializerFor, serializerForFn } from "./serialize";
import { PbTypeDefList } from "./type_list";

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

    toModel(): Struct {
        return new Struct(PbTypeDef.defsToFields(this.attrDefs));
    }
}

@jsonObject
export class PbTypeDefStruct {
    @jsonMapMember(String, () => PbTypeDef, serializerForFn(() => PbTypeDef))
    attrDefs!: Map<string, PbTypeDef>;

    toModel(): Struct {
        return new Struct(PbTypeDef.defsToFields(this.attrDefs));
    }
}

@jsonObject
export class PbTypeDefUnion {
    @jsonArrayMember(() => PbTypeDef) type!: PbTypeDef[];

    toModel(): Union {
        const values = this.type.map(t => t.toValue());

        if (
            values.every(
                v =>
                    v instanceof Primitive ||
                    v instanceof ElementRef ||
                    v instanceof CollectionDecorator
            )
        )
            return new Union(values as FieldValue[]);

        throw new Error(
            "Cannot produce Union, some members are not a Primitive/Reference/TypeDecorator: " +
                values.map(v => v.constructor.name).join(", ")
        );
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

    static defsToFields(attrDefs: Map<string, PbTypeDef>): Field[] {
        return sortLocationalArray(
            Array.from(attrDefs).map(([key, value]) => value?.toField(key))
        );
    }

    toValue(): GenericValue {
        if (this.primitive)
            return new Primitive(
                this.primitive,
                (this.constraint ?? []).map(c => c.toModel())
            );
        if (this.relation) return this.relation.toModel();
        if (this.typeRef) return this.typeRef.ref.toModel();
        if (this.set) return new CollectionDecorator(this.set.toModel(), true);
        if (this.sequence)
            return new CollectionDecorator(this.sequence.toModel(), false);
        if (this.list)
            return new CollectionDecorator(this.list.toModel(), false);
        if (this.enum) return this.enum.toModel();
        if (this.tuple) return this.tuple.toModel();
        if (this.oneOf) return this.oneOf.toModel();

        throw new Error(`Error converting type: ${JSON.stringify(this)}`);
    }

    // `isInner` specifies whether a type exists within something else and is not a type definition.
    // It is true by default, meaning the type is a nested definition or a parameter.
    // When it is false, it means it is a top level `Type` definition and therefore may be an `Alias`.
    toModel(name?: string | undefined, isInner: boolean = true): Element {
        let value = this.toValue();
        const params = {
            discriminator: "",
            name: name ?? "",
            optional: this.opt,
            value,
            locations: this.sourceContexts,
            tags: getTags(this.attrs),
            annos: getAnnos(this.attrs),
        };

        if (value instanceof Struct) {
            if (name)
                return new Type(name, !!this.relation, value.fields, params);

            params.discriminator = this.relation ? "!table" : "!type";
        }
        if (value instanceof Enum) params.discriminator = "!enum";
        if (value instanceof Union) params.discriminator = "!union";

        if (
            !isInner &&
            (value instanceof Primitive ||
                value instanceof CollectionDecorator ||
                value instanceof ElementRef)
        )
            params.discriminator = "!alias";

        if (!params.discriminator) {
            if (
                value instanceof Primitive ||
                value instanceof CollectionDecorator ||
                value instanceof ElementRef
            )
                return new Field(params.name, value, this.opt, params);
            else
                throw new Error(
                    "No discriminator but not a Field-compatible value. Type of value is " +
                        value.constructor.name
                );
        }

        return new GenericElement(params);
    }

    toField(name: string): Field {
        const field = this.toModel(name);
        if (!(field instanceof Field))
            throw new Error(
                "Cannot produce Field, requested PbTypeDef is not a Field"
            );
        return field;
    }
}
