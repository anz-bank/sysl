import { Annotation, Tag } from "../model/attribute";
import { ILocational } from "../model/common";
import { PbAttribute } from "../pbModel/attribute";

const tagsAttrName = "patterns";

function getStart<T extends ILocational>(item: T): number {
    if (!item.locations?.length) {
        return 0;
    }
    const firstLoc = item.locations[0];
    return firstLoc.start ? firstLoc.start.line : firstLoc.end.line;
}

export function sortLocationalArray<T extends ILocational>(array: T[]): T[] {
    return array.sort((i1, i2) => {
        return getStart(i1) - getStart(i2);
    });
}

export function getTags(attrs: Map<string, PbAttribute>): Tag[] {
    const tagAttr = Array.from(attrs).find(
        ([key]) => key === tagsAttrName
    )?.[1];
    if (!tagAttr) {
        return [];
    }
    if (!tagAttr.a) {
        throw new Error("Tags attribute must have an array value");
    }
    return tagAttr.a.elt.map(e => e.toTag());
}

export function getAnnos(attrs: Map<string, PbAttribute>): Annotation[] {
    return Array.from(attrs)
        .filter(([key]) => key !== tagsAttrName)
        .map(([key, value]) => value.toAnno(key));
}
