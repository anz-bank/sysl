import "reflect-metadata";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams, IParentElement } from "./element";
import { Field } from "./field";
import { Application } from "./application";

export class Type extends Element implements IParentElement<Field> {
    override get parent(): Application | undefined {
        return super.parent as Application;
    }
    override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(
        name: string,
        public isTable: boolean = false,
        public children: Field[] = [],
        p?: IElementParams<Application>
    ) {
        super(name, p?.locations ?? [], p?.annos ?? [], p?.tags ?? [], p?.model, p?.parent);
        this.attachSubitems();
    }

    toSysl(): string {
        return this.render(this.isTable ? "!table" : "!type", this.children);
    }

    override toDto() {
        return { ...super.toDto(), children: this.children.map(e => e.toDto()) };
    }

    static fromDto(dto: ReturnType<Type["toDto"]>): Type {
        return new Type(dto.name, false, dto.children.map(Field.fromDto), Element.paramsFromDto(dto));
    }

    toRef(): ElementRef {
        return this.parent!.toRef().with({ typeName: this.name });
    }

    override toString(): string {
        return `${this.isTable ? "!table" : "!type"} ${this.safeName}`;
    }

    clone(context = new CloneContext(this.model)): Type {
        const params = {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        };

        return new Type(this.name, this.isTable, context.recurse(this.children), params);
    }
}
