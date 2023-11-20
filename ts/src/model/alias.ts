import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { CollectionDecorator } from "./decorator";
import { Element, IElementParams } from "./element";
import { Primitive } from "./primitive";
import { Application } from "./application";

export class Alias extends Element {
    public override get parent(): Application | undefined {
        return super.parent as Application;
    }
    public override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(name: string, public value: AliasValue, p: IElementParams<Application>) {
        super(name, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
        this.attachSubitems();
    }

    toSysl(): string {
        return this.render("!alias", this.value.toSysl(true, this.parent?.toRef()));
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
