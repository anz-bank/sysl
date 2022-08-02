import { Location } from "../common/location";
import { Annotation, Tag } from "./attribute";
import { Model } from "./model";

/** An element that can be serialized to Sysl source. */
export interface IRenderable {
    toSysl(): string;
}

/**
 * An element that has source context locating it in one or more sourcefiles.
 *
 * All objects with locations in Sysl source exist in a Sysl model. The {@code model} property is
 * a reference to that model.
 *
 * {@code model} may be undefined if an object has been created detached from a model.
 */
export interface ILocational {
    locations: Location[];
    model?: Model;
}

/**
 * An object that has a parent in the Sysl model.
 *
 * Follow the chain of {@code parent} properties should lead to an {@link Application} which will
 * have a falsey parent.
 *
 * {@code parent} can also be undefined if an object has been created detached from a model.
 */
export interface IChild {
    parent?: IElement;
}

/**
 * An object in a Sysl model that can have nested objects (children, annotations and tags).
 */
export interface IElement extends ILocational, IRenderable, IChild {
    annos: Annotation[];
    tags: Tag[];
}

/** Common set of properties that are received by all {@link IElement} constructors. */
export type IElementParams = {
    annos?: Annotation[];
    tags?: Tag[];
    parent?: IElement;
    locations?: Location[];
    model?: Model;
};

export function setParent(parent: IElement, children: IChild[]) {
    children.forEach(child => (child.parent = parent));
}

export function setParentDeep(parent: IElement, ...childrenArrays: IChild[][]) {
    childrenArrays.forEach(children => setParent(parent, children));
}
