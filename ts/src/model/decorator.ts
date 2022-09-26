import { IRenderable } from "./common";

/** Used for `sequence of` and `set of` prefixes on data types. */
// TODO: Once Reference is removed, restrict the innerType to only what can actually be decorated (e.g. disallow apps and fields).
export class CollectionDecorator implements IRenderable {
    constructor(public innerType: IRenderable, public isSet: boolean) { }

    toSysl(): string {
        return (this.isSet ? "set of " : "sequence of ") + this.innerType.toSysl();
    }
}
