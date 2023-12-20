import { Element, IElementParams } from "./element";
import { Primitive, TypePrimitive } from "./primitive";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Type } from "./type";
import { FieldValue } from "./fieldValue";

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
        return {
            ...super.toDto(),
            ...FieldValue.toDto(this.value),
            optional: this.optional,
        };
    }

    static fromDto(dto: ReturnType<Field["toDto"]>): Field {
        return new Field(dto.name, FieldValue.fromDto(dto), dto.optional, Element.paramsFromDto(dto));
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
        return new Field(this.name, this.value, this.optional, {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}
