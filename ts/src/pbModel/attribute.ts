import "reflect-metadata";
import { jsonArrayMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../common/location";
import { Annotation, AnnoValue, Tag } from "../model/attribute";
import { Deserializer } from "typedjson/lib/types/deserializer";

@jsonObject
export class PbAttributeArray {
    @jsonArrayMember(() => PbAttribute) elt!: PbAttribute[];
}

@jsonObject
export class PbAttribute {
    @jsonMember({ deserializer: PbAttribute.deserializeString }) s?: string;
    @jsonMember n?: number;
    @jsonMember(PbAttributeArray) a?: PbAttributeArray;
    @jsonArrayMember(Location) sourceContexts!: Location[];

    /**
     * Returns the string value to set on the attribute given the value read from the input.
     *
     * String attribute values in Sysl source use backslashes to escape quotes and thus also other backslashes. These
     * strings are exported literally (as-is) to JSON (e.g. if you see four backslashes in Sysl, there will be four
     * backslashes in the serialized JSON). However when loading the JSON into this model, the backslashes are
     * interpreted as escapes, and the unescaped values are loaded. In general this is fine, except for escaped
     * backslashes which, when unescaped, become single backslashes which begin an escape sequence with the next
     * character. Therefore we re-esacpe the backslashes to avoid accidentally escaping something else with them.
     */
    static deserializeString(input?: string): string | undefined {
        return input?.replaceAll("\\", "\\\\");
    }

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
        return new Tag(tagName, { locations: this.sourceContexts });
    }

    toAnno(name: string): Annotation {
        return new Annotation(name, this.getAnnoValue(), { locations: this.sourceContexts ?? [] });
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
