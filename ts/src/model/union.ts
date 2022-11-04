import { indent } from "../common/format";
import { IRenderable } from "./common";
import { FieldValue } from "./field";

export class Union implements IRenderable {
    constructor(public members: FieldValue[]) {}

    toSysl(): string {
        return `${indent(this.members.length ? this.members.map(o => o.toSysl()).join("\n") : "...")}`;
    }
}
