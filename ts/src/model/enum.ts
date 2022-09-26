import { indent } from "../common/format";
import { IRenderable } from "./common";

export class Enum implements IRenderable {
    members: EnumValue[];

    constructor(members: EnumValue[]) {
        this.members = members ?? [];
    }

    toSysl(): string {
        return `${indent(this.members.map(e => e.toSysl()).join("\n"))}`;
    }
}

export class EnumValue implements IRenderable {
    name: string;
    value: number;

    constructor(name: string, value: number) {
        this.name = name;
        this.value = value;
    }

    toSysl(): string {
        return `${this.name}: ${this.value}`;
    }
}
