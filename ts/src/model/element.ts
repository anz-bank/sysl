import { Location } from "../common";
import { Annotation, Tag } from "./attribute";
import { IChild, ILocational, IRenderable } from "./common";
import { Model } from "./model";

/**
 * An object in a Sysl model that can have nested objects (children, annotations and tags).
 */
export abstract class Element implements ILocational, IRenderable, IChild {
    constructor(
        public name: string,
        public locations: Location[],
        public annos: Annotation[],
        public tags: Tag[],
        public model?: Model,
        public parent?: Element
    ) {}
    abstract toSysl(): string;
}

export abstract class ParentElement<TChild extends Element> extends Element {
    abstract get children(): TChild[];
}

/** Common set of properties that are received by all {@link IElement} constructors. */
export type IElementParams = {
    annos?: Annotation[];
    tags?: Tag[];
    parent?: Element;
    locations?: Location[];
    model?: Model;
};

/**
 * Sets each of {@code children}'s {@code parent} and {@code model} properties to {@code parent}
 * and its {@code model}.
 */
export function setParentAndModel(parent: Element, children: IChild[]) {
    children.forEach(child => (child.parent = parent));
    setModel(parent.model, ...children);
}

/** Sets each of {@code children}'s model to {@code model}. */
export function setModel(model?: Model, ...children: IChild[]) {
    children.forEach(child => (child.model = model));
}

export function setParentAndModelDeep(
    parent: Element,
    ...childrenArrays: IChild[][]
) {
    childrenArrays.forEach(children => {
        setParentAndModel(parent, children);
    });
}
export function setModelDeep(model?: Model, ...childrenArrays: IChild[][]) {
    childrenArrays.forEach(children => {
        setModel(model, ...children);
    });
}
