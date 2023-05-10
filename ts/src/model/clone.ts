import { Annotation, Tag } from "./attribute";
import { Model } from "./model";
import { ILocational } from "./common";

export type ICloneable = { clone(context?: CloneContext): ICloneable };

export class CloneContext {
    constructor(public model?: Model, public filter: ModelFilter = ModelFilters.Default, public depth: number = 0) {
        filter.bind(this);
    }

    apply<T extends ICloneable>(item: T | undefined): T | undefined {
        if (!item) return undefined;
        return this.filter(this, item) ? (item.clone(this) as T) : undefined;
    }

    applyUnder<T extends ICloneable>(item: T | undefined): T | undefined {
        if (!item) return undefined;
        return this.recurse([item])[0];
    }

    recurse<T extends ICloneable>(arr: T[] | undefined): T[] {
        if (!arr) return [];
        const childFilter = new CloneContext(this.model, this.filter, this.depth + 1);
        return arr.map(i => childFilter.apply(i)).filter(i => i) as T[];
    }
}

export type ModelFilter = (context: CloneContext, item: ICloneable | ILocational) => boolean;

export class ModelFilters {
    static Default: ModelFilter = (_, __) => true;
    static ExcludeAnnosAndTags: ModelFilter = (_, item) => !(item instanceof Annotation || item instanceof Tag);
    static OnlyFromFile(syslPath: string): ModelFilter {
        return (_, item) => {
            if ("locations" in item) return !item.locations.length || item.locations.some(l => l.file == syslPath);
            return true;
        };
    }
}
