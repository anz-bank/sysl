import { IEnumerable } from "linq-to-typescript";

declare global {
    interface Array<T> extends IEnumerable<T> {}
    interface Map<K, V> extends IEnumerable<[K, V]> {}
    interface Set<T> extends IEnumerable<T> {}
    interface String extends IEnumerable<string> {}
}
