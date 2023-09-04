/**
 * Creates a flattened array of values by recursively running each element in {@link collection} through
 * {@link iteratee} and flattening the mapped results.
 * @param collection The collection to iterate over.
 * @param iteratee The function invoked per iteration.
 * @returns Returns the new flattened array.
 */
export function flatMapDeep<T>(collection: T[], iteratee: (item: T) => T[]): T[] {
    const children = collection.flatMap(iteratee);
    if (!children.length) return collection;
    return [...collection, ...flatMapDeep(children, iteratee)];
}

/**
 * Provides support for lazy initialization.
 */
export class Lazy<T> {
    #value?: T;
    #valueCreated = false;

    /**
     * Initializes a new instance of the {@link Lazy} class. When lazy initialization occurs, the specified
     * initialization function is used.
     * @param valueFactory The function that is invoked to produce the lazily initialized value when it is needed.
     */
    constructor(public valueFactory: () => T) {}

    /** Gets the lazily initialized value of the current {@link Lazy} instance. */
    public get value(): T {
        if (!this.#valueCreated) {
            this.#value ??= this.valueFactory();
            this.#valueCreated = true;
        }
        return this.#value!;
    }
}
