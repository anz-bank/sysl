import { IRenderable } from "./common";
import { Field } from "./field";

// The type of `!type`, `!table` and endpoints.
export class Struct implements IRenderable {
    fields: Field[];

    constructor(elements: Field[]) {
        this.fields = elements ?? [];
    }

    toSysl(): string {
        throw new Error("Not implemented");
    }
}
