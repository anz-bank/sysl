import { IEnumerable } from "linq-to-typescript";
import { Serializable, TypedJSON } from "typedjson";
import { Annotation, Tag } from "./model/attribute";
import { ComplexType } from "./model/common";
import { Element } from "./model/element";
import { Application } from "./model/model";
import { Endpoint, Statement } from "./model/statement";
import { Type } from "./model/type";
import { PbAttribute } from "./pbModel/attribute";

export function serializeMap<T>(mapObject: Map<string, T>): { [key: string]: T } {
    const obj: { [key: string]: T; } = {};
    mapObject.forEach((val, key) => {
        obj[key as string] = val;
    });
    return obj;
}

export function deserializeMap<T>(rootType: Serializable<T>, stringifiedMapObject: { [key: string]: T }) {
    const map = new Map<string, T>();
    if (stringifiedMapObject) {
        Object.entries(stringifiedMapObject)
            .forEach(([type, displayFields]) => {
                map.set(
                    type,
                    TypedJSON.parse(displayFields, rootType) as T
                );
            });
    }
    return map;
}

export function serializerFor<T>(type: Serializable<T>, nestedUnder?: string) {
    const deserializer = (stringifiedMapObject: { [key: string]: T; }) => deserializeMap(type, stringifiedMapObject);
    const nestedDeserializer = (wrapper: any) => wrapper ? deserializer(wrapper[nestedUnder!]) : wrapper;
    return { serializer: serializeMap, deserializer: nestedUnder ? nestedDeserializer : deserializer };
}

export function indent(text: string): string {
    return `${text.split("\n").select(l => `    ${l}`).toArray().join("\n")}`;
}

function getStart(item: Application | Tag | Annotation | Statement | Endpoint | Type): number {
    if (item.locations) {
        return item.locations.first().start ? item.locations.first().start.line : item.locations.first().end.line;
    }
    return 0;
}

export function sortElements<T extends ComplexType>(elements: Element<T>[]): Element<T>[] {
    return elements.sort((i1, i2) => {
        return getStart(i1.content) - getStart(i2.content);
    })
}

export function joinedAppName(name: string[], compact: boolean = false): string {
    return name.join(compact ? "::" : " :: ");
}

export function getTags(attrs: IEnumerable<[string, PbAttribute]> | undefined): Tag[] {
    const tagAttr = attrs?.any() ? attrs.firstOrDefault(a => a[0] === "patterns") : undefined;
    if (tagAttr) {
        if (tagAttr[1].a) {
            return tagAttr[1].a.elt.select(e => e.toTag()).toArray() ?? [];
        }
        else {
            throw new Error("Tags attribute must have an array value");
        }
    }
    return [];
}

export function getAnnos(attrs: Map<string, PbAttribute> | undefined): Annotation[] {
    return attrs?.any() ? attrs.where(a => a[0] != "patterns").select(a => a[1].toAnno(a[0])).toArray() : [];
}
