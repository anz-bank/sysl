import "reflect-metadata";
import { indent } from "../common/format";
import { Location } from "../common/location";
import { IChild, IElement, ILocational, IRenderable } from "./common";
import { Model } from "./model";

export type AnnoValue = string | number | AnnoValue[];

export type AnnotationParams = {
    name: string;
    value: AnnoValue;
    locations?: Location[];
    model?: Model;
};

export class Annotation implements IChild, ILocational, IRenderable {
    value: AnnoValue;
    name: string;
    locations: Location[];
    parent?: IElement;
    model?: Model;

    constructor({ name, value, locations, model }: AnnotationParams) {
        this.name = name;
        this.value = value;
        this.locations = locations ?? [];
        this.model = model;
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
                    .map((item, index) =>
                        index > 0
                            ? valueString(item)
                            : valueString(item).trimStart()
                    )
                    .join(",")}]`;
            }
        }
        return `${this.name} =${valueString(this.value)}`;
    }
}

export type TagParams = {
    value: AnnoValue;
    locations?: Location[];
    model?: Model;
};

export class Tag implements IChild, ILocational, IRenderable {
    value: AnnoValue;
    locations: Location[];
    parent?: IElement;
    model?: Model;

    constructor({ value, locations, model }: TagParams) {
        this.value = value;
        this.locations = locations ?? [];
        this.model = model;
    }

    toSysl(): string {
        return `~${this.value}`;
    }
}
