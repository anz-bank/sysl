import "reflect-metadata";
import { ElementRef } from "./common";
import { CloneContext } from "./clone";
import { IElementParams, ParentElement } from "./element";
import { Field } from "./field";

export class Type extends ParentElement<Field> {
    constructor(name: string, public isTable: boolean = false, public children: Field[] = [], p?: IElementParams) {
        super(name, p?.locations ?? [], p?.annos ?? [], p?.tags ?? [], p?.model, p?.parent);
        this.attachSubitems();
    }

    toRef(): ElementRef {
        return this.parent!.toRef().with({ typeName: this.name });
    }

    toSysl(): string {
        return this.render(this.isTable ? "!table" : "!type", this.children);
    }

    override toString(): string {
        return `${this.isTable ? "!table" : "!type"} ${this.safeName}`;
    }

    clone(context = new CloneContext(this.model)): Type {
        const params = {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
        };

        return new Type(this.name, this.isTable, context.recurse(this.children), params);
    }
}
