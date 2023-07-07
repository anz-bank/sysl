import { Annotation, Tag } from "./attribute";
import { Model } from "./model";
import { ILocational } from "./common";

/**
 * Represents a type that is cloneable. The `clone` function is expected to return a deep copy of the item it was called
 * on.
 */
export type ICloneable = { clone(context?: CloneContext): ICloneable };

/**
 * Holds context information for a clone operation at a specific depth. When cloning reaches nested objects, a new
 * instance of {@link CloneContext} is created with a greater depth.
 */
export class CloneContext {
    /**
     * Creates a new instance of the {@link CloneContext} class.
     * @param model The target {@link Model} where items are being cloned to. The {@link ILocational.model} property
     * will be set to this value for newly cloned objects. If `undefined`, it will no be set.
     * @param filter The {@link ModelFilter} used to filter out certain items when cloning. Often used filters are
     * defined on {@link ModelFilters}.
     * @param depth The current depth of the cloning process, which is incremented every time {@link recurse} is called.
     * The value is usually 0 at the {@link Model} level, 1 at the {@link Application} level, 2 at the {@link Type} or
     * {@link Endpoint} level, etc. Useful for making decision at the filter or when you want to draw a tree of objects.
     * @param keepLocation Specify `true` to reserve the {@link ILocational.locations} of cloned items. By default
     * location data isn't copied because the cloned items represent something that hans't been persisted to a file yet.
     * But sometimes location data should be kept, such as when using clone to create a filtered view of a model.
     */
    constructor(
        public model?: Model,
        public filter = ModelFilters.Default,
        public depth = 0,
        public keepLocation = false
    ) {
        filter.bind(this);
    }

    /**
     * Recursively clones a single item. Does not increment the depth.
     * @param item The item to clone recursively. Each item will be filtered, then have {@link ICloneable.clone}
     * called.
     * @returns The cloned item, or `undefined` if either the item was filtered out, or `undefined` was supplied as the
     * argument.
     */
    apply<T extends ICloneable>(item: T | undefined): T | undefined {
        if (!item) return undefined;
        return this.filter(this, item) ? (item.clone(this) as T) : undefined;
    }

    /**
     * Recursively clones a single item, under a new clone context with an incremented depth.
     * @param item The item to clone recursively. Each item will be filtered, then have {@link ICloneable.clone}
     * called.
     * @returns The cloned item, or `undefined` if either the item was filtered out, or `undefined` was supplied as the
     * argument.
     */
    applyUnder<T extends ICloneable>(item: T | undefined): T | undefined {
        if (!item) return undefined;
        return this.recurse([item])[0];
    }

    /**
     * Recursively clones an array of items, under a new clone context with an incremented depth.
     * @param arr An array of items to clone recursively. Each item will be filtered, then have {@link ICloneable.clone}
     * called. If `undefined` was specified, it will be treated like an empty array.
     * @returns An array of cloned items. If no item matches the filter, an empty array will be returned.
     */
    recurse<T extends ICloneable>(arr: T[] | undefined): T[] {
        if (!arr) return [];
        const childContext = new CloneContext(this.model, this.filter, this.depth + 1, this.keepLocation);
        return arr.map(i => childContext.apply(i)).filter(i => i) as T[];
    }
}

/**
 * A function signature used to filter out items when cloning a model. The function accepts a context and the item to be
 * checked, and returns true if the item should be cloned, or false if the item should be skipped.
 */
export type ModelFilter = (context: CloneContext, item: ICloneable | ILocational) => boolean;

export class ModelFilters {
    /** Default filter, doesn't filter out anything. */
    static Default: ModelFilter = (_, __) => true;

    /** Filters out all annotations and tags. */
    static ExcludeAnnosAndTags: ModelFilter = (_, item) => !(item instanceof Annotation || item instanceof Tag);

    /**
     * Filters out everything that isn't from the specified file. If an item doesn't have any location data, it is kept.
     * @param syslPath The root-relative path of the file to keep items from. Call {@link Model.convertSyslPath} if
     * needed.
     */
    static OnlyFromFile(syslPath: string): ModelFilter {
        return (_, item) => {
            if ("locations" in item) return !item.locations.length || item.locations.some(l => l.file == syslPath);
            return true;
        };
    }
}
