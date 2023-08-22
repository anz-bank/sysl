import "reflect-metadata";
import { indent } from "../common/format";
import { Location } from "../common/location";
import { IChild, ILocational, IRenderable } from "./common";
import { CloneContext } from "./clone";
import { Element } from "./element";
import { Model } from "./model";

export type AnnoValue = string | AnnoValue[];

export type AnnotationParams = {
    locations?: Location[];
    model?: Model;
    parent?: Element;
};

export class Annotation implements IChild, ILocational, IRenderable {
    locations: Location[];
    parent?: Element;
    model?: Model;

    constructor(public name: string, public value: AnnoValue, { locations, model, parent }: AnnotationParams = {}) {
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

    toString(): string {
        return `@${this.name} = ...`;
    }

    clone(context = new CloneContext(this.model)): Annotation {
        return new Annotation(this.name, this.cloneValue(this.value), {
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }

    private cloneValue(value: AnnoValue): AnnoValue {
        if (typeof value === "string" || typeof value === "number") return value;
        return value.map(v => this.cloneValue(v));
    }
}

export type TagParams = {
    locations?: Location[];
    model?: Model;
    parent?: Element;
};

export class Tag implements IChild, ILocational, IRenderable {
    locations: Location[];
    parent?: Element;
    model?: Model;

    constructor(public name: string, { locations, model, parent }: TagParams = {}) {
        this.locations = locations ?? [];
        this.model = model ?? parent?.model;
        this.parent = parent;
    }

    toSysl(): string {
        return `~${this.name}`;
    }

    toString(): string {
        return `[${this.toSysl()}]`;
    }

    clone(context = new CloneContext(this.model)): Tag {
        return new Tag(this.name, {
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}
