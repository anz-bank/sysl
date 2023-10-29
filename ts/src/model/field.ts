import { Element, IElementParams } from "./element";
import { Primitive } from "./primitive";
import { CollectionDecorator } from "./decorator";
import { ElementRef } from "./common";
import { CloneContext } from "./clone";

export type FieldValue = Primitive | ElementRef | CollectionDecorator;

export class Field extends Element {
    constructor(name: string, public value: FieldValue, public optional: boolean = false, p?: IElementParams) {
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
        return this.parent!.toRef().with({ fieldName: this.name });
    }

    override toSysl(): string {
        let value: string = this.value.toSysl(true, this.parent ? this.toRef() : undefined);
        value += this.optional ? "?" : "";
        //TODO: Everyone who uses Field with empty name should use FieldValue instead.
        let name = this.name ? `${this.safeName} <:` : "";
        return this.render(name, "", value, false);
    }

    override toString(): string {
        return `${this.name ? this.name : "[param]"} <: ${this.value.toSysl(true)}${this.optional ? "?" : ""}`;
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
