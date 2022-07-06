import { IEnumerable } from "linq-to-typescript";
import { Annotation, Tag } from "./model/attribute";
import { ILocational } from "./model/common";
import { PbAttribute } from "./pbModel/attribute";

function getStart<T extends ILocational>(item: T): number {
    if (item.locations) {
        return item.locations.first().start
            ? item.locations.first().start.line
            : item.locations.first().end.line;
    }
    return 0;
}

export function sortLocationalArray<T extends ILocational>(array: T[]): T[] {
    return array.sort((i1, i2) => {
        return getStart(i1) - getStart(i2);
    });
}

export function getTags(
    attrs: IEnumerable<[string, PbAttribute]> | undefined
): Tag[] {
    const tagAttr = attrs?.any()
        ? attrs.firstOrDefault(a => a[0] === "patterns")
        : undefined;
    if (tagAttr) {
        if (tagAttr[1].a) {
            return tagAttr[1].a.elt.select(e => e.toTag()).toArray() ?? [];
        } else {
            throw new Error("Tags attribute must have an array value");
        }
    }
    return [];
}

export function getAnnos(
    attrs: Map<string, PbAttribute> | undefined
): Annotation[] {
    return attrs?.any()
        ? attrs
              .where(a => a[0] != "patterns")
              .select(a => a[1].toAnno(a[0]))
              .toArray()
        : [];
}
