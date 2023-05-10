import { IRenderable } from "./common";
import { CloneContext } from "./clone";
import { Element } from "./element";

/** Used for `sequence of` and `set of` prefixes on data types. */
// TODO: Change inner type to be a 'Primitive | ElementRef'
export class CollectionDecorator implements IRenderable {
    constructor(public innerType: Element, public isSet: boolean) {}

    toSysl(): string {
        return (this.isSet ? "set of " : "sequence of ") + this.innerType.toSysl();
    }

    toString(): string {
        return this.toSysl();
    }

    clone(context: CloneContext = new CloneContext()): CollectionDecorator {
        return new CollectionDecorator(this.innerType.clone(context), this.isSet);
    }
}
