import { Serializable, TypedJSON } from "typedjson";

export function serializeMap<T>(mapObject: Map<string, T>): {
    [key: string]: T;
} {
    const obj: { [key: string]: T } = {};
    mapObject.forEach((val, key) => {
        obj[key as string] = val;
    });
    return obj;
}

export function deserializeMap<T>(
    rootType: Serializable<T>,
    stringifiedMapObject: { [key: string]: T }
) {
    const map = new Map<string, T>();
    if (stringifiedMapObject) {
        Object.entries(stringifiedMapObject).forEach(
            ([type, displayFields]) => {
                map.set(type, TypedJSON.parse(displayFields, rootType) as T);
            }
        );
    }
    return map;
}

export function serializerFor<T>(type: Serializable<T>, nestedUnder?: string) {
    const deserializer = (stringifiedMapObject: { [key: string]: T }) =>
        deserializeMap(type, stringifiedMapObject);
    const nestedDeserializer = (wrapper: any) =>
        wrapper ? deserializer(wrapper[nestedUnder!]) : wrapper;
    return {
        serializer: serializeMap,
        deserializer: nestedUnder ? nestedDeserializer : deserializer,
    };
}
