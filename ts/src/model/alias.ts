import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams } from "./element";
import { Application } from "./application";
import { FieldValue } from "./fieldValue";

export class Alias extends Element {
    public override get parent(): Application | undefined {
        return super.parent as Application;
    }
    public override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(name: string, public value: FieldValue, p: IElementParams<Application>) {
        super(name, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
        this.attachSubitems();
    }

    toSysl(): string {
        return this.render("!alias", this.value.toSysl(true, this.parent?.toRef()));
    }

    override toString(): string {
        return `!alias ${this.safeName}`;
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    public override toDto() {
        return { ...super.toDto(), value: FieldValue.toDto(this.value) };
    }

    static fromDto(dto: ReturnType<Alias["toDto"]>): Alias {
        return new Alias(dto.name, FieldValue.fromDto(dto.value), Element.paramsFromDto(dto));
    }

    clone(context = new CloneContext(this.model)): Alias {
        return new Alias(this.name, this.value, {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}
