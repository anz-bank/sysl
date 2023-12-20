import { IRenderable } from "./common";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams } from "./element";
import { Application } from "./application";

export class Enum extends Element {
    public override get parent(): Application | undefined {
        return super.parent as Application;
    }
    public override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(name: string, public children: EnumValue[], p: IElementParams<Application>) {
        super(name, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
    }

    toSysl(): string {
        return this.render("!enum", this.children);
    }

    override toString(): string {
        return `!enum ${this.safeName}`;
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    override toDto() {
        return { ...super.toDto(), children: Object.fromEntries(this.children.map(c => [c.name, c.value])) };
    }

    static fromDto(dto: ReturnType<Enum["toDto"]>): Enum {
        return new Enum(
            dto.name,
            Object.entries(dto.children).map(([n, v]) => new EnumValue(n, v)),
            Element.paramsFromDto(dto)
        );
    }

    clone(context = new CloneContext(this.model)): Enum {
        const params = {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        };
        return new Enum(this.name, context.recurse(this.children), params);
    }
}

export class EnumValue implements IRenderable {
    name: string;
    value: number;

    constructor(name: string, value: number) {
        this.name = name;
        this.value = value;
    }

    toSysl(): string {
        return `${this.name}: ${this.value}`;
    }

    toString(): string {
        return this.toSysl();
    }

    clone(_context: CloneContext): EnumValue {
        return new EnumValue(this.name, this.value);
    }
}
