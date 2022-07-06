import "reflect-metadata";
import { Location } from "../location";
import { indent } from "../format";
import { ILocational, IRenderable } from "./common";

export type AnnoValue = string | number | AnnoValue[];

export class Annotation implements ILocational, IRenderable {
    value: AnnoValue;
    name: string;
    locations: Location[];

    constructor(name: string, locations: Location[], value: AnnoValue) {
        this.name = name;
        this.value = value;
        this.locations = locations;
    }

    toSysl(): string {
        function valueString(v: AnnoValue): string {
            if (typeof v == "string") {
                if (v.contains("\n")) {
                    return (
                        `:\n` +
                        indent(`| ${v.trimEnd().split("\n").join("\n| ")}`)
                    );
                } else {
                    return ` "${v}"`;
                }
            } else if (typeof v == "number") {
                return ` "${v}"`;
            } else {
                return ` [${v
                    .select(
                        (item, index) =>
                            `${
                                index > 0
                                    ? valueString(item)
                                    : valueString(item).trimStart()
                            }`
                    )
                    .toArray()
                    .join(",")}]`;
            }
        }
        return `${this.name} =${valueString(this.value)}`;
    }
}

export class Tag implements ILocational, IRenderable {
    value: AnnoValue;
    locations: Location[];

    constructor(value: AnnoValue, locations: Location[]) {
        this.value = value;
        this.locations = locations;
    }

    toSysl(): string {
        return `~${this.value}`;
    }
}
