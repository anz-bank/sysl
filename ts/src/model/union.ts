import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams } from "./element";
import { FieldValue } from "./fieldValue";
import { Application } from "./application";

export class Union extends Element {
    public override get parent(): Application | undefined {
        return super.parent as Application;
    }
    public override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(name: string, public children: FieldValue[], p: IElementParams<Application>) {
        super(name, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
    }

    toSysl(): string {
        return this.render("!union", this.children);
    }

    override toString(): string {
        return `!union ${this.safeName}`;
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    public override toDto() {
        return { ...super.toDto(), children: this.children.map(v => FieldValue.toDto(v)) };
    }

    static fromDto(dto: ReturnType<Union["toDto"]>): Union {
        return new Union(dto.name, dto.children.map(FieldValue.fromDto), Element.paramsFromDto(dto));
    }

    clone(context = new CloneContext(this.model)): Union {
        const params = {
            annos: context.recurse(this.annos),
            tags: context.recurse(this.tags),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        };

        return new Union(this.name, [...this.children], params);
    }
}
