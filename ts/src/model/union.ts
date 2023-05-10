import { ElementRef } from "./common";
import { CloneContext } from "./clone";
import { Element, IElementParams } from "./element";
import { FieldValue } from "./field";

export class Union extends Element {
    constructor(name: string, public members: FieldValue[], { annos, tags, locations, parent, model }: IElementParams) {
        super(name, locations ?? [], annos ?? [], tags ?? [], model, parent);
        this.attachSubitems();
    }

    toSysl(): string {
        return this.render("!union", this.members);
    }

    override toString(): string {
        return `!union ${this.safeName}`;
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    clone(context = new CloneContext(this.model)): Union {
        const params = {
            annos: context.recurse(this.annos),
            tags: context.recurse(this.tags),
            model: context.model ?? this.model,
        };

        return new Union(this.name, context.recurse(this.members), params);
    }
}
