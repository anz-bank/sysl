import { IRenderable } from "./common";
import { ElementRef } from "./elementRef";
import { Primitive } from "./primitive";
import { FieldValue } from "./fieldValue";

/** Used for `sequence of` and `set of` prefixes on data types. */
export class CollectionDecorator implements IRenderable {
    public readonly innerType: Primitive | ElementRef;

    constructor(innerType: FieldValue, public readonly isSet: boolean) {
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
}
