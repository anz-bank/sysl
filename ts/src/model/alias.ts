import { indent } from "../common/format";
import { IRenderable } from "./common";
import { GenericValue } from "./genericElement";

export class Alias implements IRenderable {
    type: GenericValue;

    constructor(type: GenericValue) {
        this.type = type;
    }

    toSysl(): string {
        return `${indent(this.type.toSysl())}`;
    }
}
