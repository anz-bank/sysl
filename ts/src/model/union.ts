import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams } from "./element";
import { FieldValue } from "./field";
import { Application } from "./application";

export class Union extends Element {
    public override get parent(): Application | undefined {
        return super.parent as Application;
    }
    public override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(name: string, public members: FieldValue[], p: IElementParams<Application>) {
        super(name, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
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
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        };

        return new Union(this.name, context.recurse(this.members), params);
    }
}
