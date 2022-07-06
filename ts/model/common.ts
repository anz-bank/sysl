import { Location } from "../location";
import { Annotation, Tag } from "./attribute";

export interface IRenderable {
    toSysl(): string;
}

export interface IDescribable {
    tags: Tag[];
    annos: Annotation[];
}

export interface ILocational {
    locations: Location[];
}
