import { IRenderable } from "./common";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Primitive } from "./primitive";
import { FieldValue } from "./field";

/** Used for `sequence of` and `set of` prefixes on data types. */
export class CollectionDecorator implements IRenderable {
    public innerType: Primitive | ElementRef;

    constructor(innerType: FieldValue, public isSet: boolean) {
        // Unwrap list of sequences/sets to prevent double decorator, e.g.:  myField(1..1) <: set of int
        if (innerType instanceof CollectionDecorator) {
            this.innerType = innerType.innerType;
            this.isSet = innerType.isSet;
        } else {
            this.innerType = innerType;
        }
    }

    toSysl(compact: boolean = false, parentRef?: ElementRef): string {
        return (this.isSet ? "set of " : "sequence of ") + this.innerType.toSysl(compact, parentRef);
    }

    toString(): string {
        return this.toSysl();
    }

    clone(context: CloneContext = new CloneContext()): CollectionDecorator {
        return new CollectionDecorator(this.innerType.clone(context), this.isSet);
    }
}
