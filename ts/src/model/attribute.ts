import "reflect-metadata";
import { Location } from "../common/location";
import { indent } from "../common/format";
import { ILocational, IRenderable } from "./common";

export type AnnoValue = string | number | AnnoValue[];

export type AnnotationParams = {
    name: string;
    value: AnnoValue;
    locations?: Location[];
};

export class Annotation implements ILocational, IRenderable {
    value: AnnoValue;
    name: string;
    locations: Location[];

    constructor({ name, value, locations }: AnnotationParams) {
        this.name = name;
        this.value = value;
        this.locations = locations ?? [];
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
                    return ` "${v}"`;
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
};

export class Tag implements ILocational, IRenderable {
    value: AnnoValue;
    locations: Location[];

    constructor({ value, locations }: TagParams) {
        this.value = value;
        this.locations = locations ?? [];
    }

    toSysl(): string {
        return `~${this.value}`;
    }
}
