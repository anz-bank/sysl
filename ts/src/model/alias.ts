import { ElementRef } from "./common";
import { CloneContext } from "./clone";
import { CollectionDecorator } from "./decorator";
import { Element, IElementParams } from "./element";
import { Primitive } from "./primitive";

export class Alias extends Element {
    constructor(name: string, public value: AliasValue, { annos, tags, locations, parent, model }: IElementParams) {
        super(name, locations ?? [], annos ?? [], tags ?? [], model, parent);
        this.attachSubitems();
    }

    toSysl(): string {
        return this.render("!alias", this.value.toSysl());
    }

    override toString(): string {
        return `!Alias ${this.safeName}`;
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    clone(context = new CloneContext(this.model)): Alias {
        return new Alias(this.name, this.value.clone(), {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}

export type AliasValue = Primitive | CollectionDecorator | ElementRef;
