import "reflect-metadata";
import { indent } from "../common/format";
import { Location } from "../common/location";
import { IChild, ILocational, IRenderable } from "./common";
import { Element } from "./element";
import { Model } from "./model";

export type AnnoValue = string | number | AnnoValue[];

export type AnnotationParams = {
    name: string;
    value: AnnoValue;
    locations?: Location[];
    model?: Model;
    parent?: Element;
};

export class Annotation implements IChild, ILocational, IRenderable {
    value: AnnoValue;
    name: string;
    locations: Location[];
    parent?: Element;
    model?: Model;

    constructor({ name, value, locations, model, parent }: AnnotationParams) {
        // TODO: Check validity of name, throw if invalid.
        this.name = name;
        this.value = value;
        this.locations = locations ?? [];
        this.model = model ?? parent?.model;
        this.parent = parent;
    }

    toSysl(): string {
        function valueString(v: AnnoValue): string {
            if (typeof v === "string") {
                if (v.includes("\n")) {
                    const lines = v
                        .trimEnd()
                        .split("\n")
                        .map(line => (line ? ` ${line}` : ""));
                    return `:\n` + indent(`|${lines.join("\n|")}`);
                } else {
                    return ` "${v.replaceAll(`"`, `\\"`)}"`;
                }
            } else if (typeof v === "number") {
                return ` "${v}"`;
            } else {
                return ` [${v
                    .map((item, index) => (index > 0 ? valueString(item) : valueString(item).trimStart()))
                    .join(",")}]`;
            }
        }
        return `${this.name} =${valueString(this.value)}`;
    }
}

export type TagParams = {
    name: string;
    locations?: Location[];
    model?: Model;
    parent?: Element;
};

export class Tag implements IChild, ILocational, IRenderable {
    name: string;
    locations: Location[];
    parent?: Element;
    model?: Model;

    constructor({ name, locations, model, parent }: TagParams) {
        this.name = name;
        this.locations = locations ?? [];
        this.model = model ?? parent?.model;
        this.parent = parent;
    }

    toSysl(): string {
        return `~${this.name}`;
    }
}
