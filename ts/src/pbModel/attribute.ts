import "reflect-metadata";
import { jsonArrayMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../common/location";
import { Annotation, AnnoValue, Tag } from "../model/attribute";

@jsonObject
export class PbAttributeArray {
    @jsonArrayMember(() => PbAttribute) elt!: PbAttribute[];
}

@jsonObject
export class PbAttribute {
    @jsonMember s?: string;
    @jsonMember n?: number;
    @jsonMember(PbAttributeArray) a?: PbAttributeArray;
    @jsonArrayMember(Location) sourceContexts!: Location[];

    getAnnoValue(): AnnoValue {
        if (this.s != null) {
            return this.s;
        } else if (this.n != null) {
            return this.n;
        } else if (this.a != null) {
            return this.a.elt?.map(i => i.getAnnoValue()) ?? [];
        } else {
            throw new Error("Missing attribute value: " + JSON.stringify(this));
        }
    }

    toTag(): Tag {
        const tagName = this.getAnnoValue();
        if (typeof tagName != "string") throw new Error("Cannot create tag with non-string name");
        return new Tag({
            name: tagName,
            locations: this.sourceContexts ?? [],
        });
    }

    toAnno(name: string): Annotation {
        return new Annotation({
            name,
            locations: this.sourceContexts ?? [],
            value: this.getAnnoValue(),
        });
    }
}

const tagsAttrName = "patterns";

export function getTags(attrs?: Map<string, PbAttribute>): Tag[] {
    if (!attrs) {
        return [];
    }
    const tagAttr = Array.from(attrs).find(([key]) => key === tagsAttrName)?.[1];
    if (!tagAttr) {
        return [];
    }
    if (!tagAttr.a) {
        throw new Error("Tags attribute must have an array value");
    }
    return tagAttr.a.elt.map(e => e.toTag());
}

export function getAnnos(attrs?: Map<string, PbAttribute>): Annotation[] {
    if (!attrs) {
        return [];
    }
    return Array.from(attrs)
        .filter(([key]) => key !== tagsAttrName)
        .map(([key, value]) => value.toAnno(key));
}
