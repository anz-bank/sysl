import { Location } from "../common/location";
import { Annotation, Tag } from "./attribute";

/** An element that can be serialized to Sysl source. */
export interface IRenderable {
    toSysl(): string;
}

/** An element that can be annotated and tagged with additional metadata. */
export interface IDescribable {
    tags: Tag[];
    annos: Annotation[];
}

/** An element that has source context locating it in one or more sourcefiles. */
export interface ILocational {
    locations: Location[];
}
