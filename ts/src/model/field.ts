import { Element, IElementParams } from "./element";
import { Primitive, TypePrimitive } from "./primitive";
import { CollectionDecorator } from "./decorator";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Type } from "./type";

export type FieldValue = Primitive | ElementRef | CollectionDecorator;

export class Field extends Element {
    public override get parent(): Type | undefined {
        return super.parent as Type;
    }
    public override set parent(app: Type | undefined) {
        super.parent = app;
    }

    constructor(name: string, public value: FieldValue, public optional: boolean = false, p?: IElementParams<Type>) {
        if (!name) throw new Error("Field name must be specified.");
        super(name, p?.locations ?? [], p?.annos ?? [], p?.tags ?? [], p?.model, p?.parent);
        this.attachSubitems();
    }

    public override toDto() {
        let value = this.value;

        let collectionType: "set" | "sequence" | undefined;
        if (value instanceof CollectionDecorator) {
            collectionType = value.isSet ? "set" : "sequence";
            value = value.innerType;
        }

        return {
            ...super.toDto(),
            optional: this.optional,
            collectionType,
            ref: value instanceof ElementRef ? value.toString() : undefined,
            primitive: value instanceof Primitive ? value.toString() : undefined,
            constraint: value instanceof Primitive ? value.constraintStr() : undefined,
        };
    }

    toRef(): ElementRef {
        // Will currently throw when this instance is a param or restParam due to not being supported by ElementRef.
        return this.parent!.toRef().with({ fieldName: this.name });
    }

    private toValue(): string {
        return this.value.toSysl(true, this.parent ? this.parent.toRef() : undefined) + (this.optional ? "?" : "");
    }

    override toSysl(omitAny: boolean = false): string {
        if (omitAny && this.value instanceof Primitive && this.value.primitive == TypePrimitive.ANY)
            return this.render("", "", undefined, false);

        return this.render(`${this.safeName} <:`, "", this.toValue(), false);
    }

    public toRestParam(): string {
        return `${this.name}=${this.toValue()}`;
    }

    override toString(): string {
        return `${this.name} <: ${this.toValue()}`;
    }

    clone(context = new CloneContext(this.model)): Field {
        return new Field(this.name, this.value.clone(context), this.optional, {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}
