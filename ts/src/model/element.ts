import { Location } from "../common";
import { Annotation, AnnoValue, Tag } from "./attribute";
import { ElementRef, IChild, ILocational, IRenderable } from "./common";
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

    /**
     * Returns an {@link ElementRef} that references this element.
     */
    abstract toRef(): ElementRef;

    /**
     * Tries to retrieve the specified annotation from this element. If it's not found, `undefined` is returned.
     * If multiple annotations of the same name exist, the first is returned.
     * @param name The name of the annotation to retrieve.
     * @returns The {@link Annotation} of the specified name, or `undefined` if no annotation with that name exists.
     */
    public tryAnno(name: string): Annotation | undefined {
        return this.annos.find(a => a.name == name);
    }

    /**
     * Retrieves the specified annotation from this element. If it's not found, an error is thrown.
     * @param name The name of the annotation to retrieve.
     * @returns The {@link Annotation} of the specified name.
     * @throws {@link Error} Thrown if no annotation matching the specified name is found.
     */
    public anno(name: string): Annotation {
        const anno = this.tryAnno(name);
        if (!anno) throw new Error(`No annotation named '${name}' was found on element '${this.name}'.`);
        return anno;
    }

    /**
     * Sets the specified annotation to the specified value, or removes it if the value specified is `undefined`.
     * If the annotation doesn't already exist on this element, a new one will be inserted with an best effort to
     * place it alphabetically.
     * @param name The name of the annotation to set.
     * @param value The value of the annotation to set, or `undefined` to remove the annotation.
     * @returns The {@link Annotation} that was updated, inserted or removed, or `undefined` if attempting to remove
     * an annotation that doesn't exist.
     */
    public setAnno(name: string, value: AnnoValue): Annotation | undefined {
        let anno = this.tryAnno(name);

        if (value == undefined) {
            if (anno) {
                this.annos = this.annos.filter(a => a !== anno);
                anno.parent = undefined;
            }
        }
        else {
            
            if (anno) {
                anno.value = value;
            }
            else {
                anno = new Annotation({ name, value });
                this.insertAnnoOrdered(anno);
            }
        }

        return anno;
    }

    /**
     * Makes a best effort to insert the annotation in an alphabetical position, by searching from the bottom.
     * Also sets the parent and model. Does not check if this annotation already exists on this element.
     * @param anno The {@link Annotation} to insert.
     * @returns The index at which the annotation was inserted.
     */
    public insertAnnoOrdered(anno: Annotation): number {
        anno.parent = this;
        anno.model = this.model;

        for (let i = this.annos.length - 1; i >= 0; i--) {
            if (anno.name >= this.annos[i].name) {
                this.annos.splice(i+1, 0, anno);
                return i+1;
            }
        }

        this.annos.splice(0, 0, anno);
        return 0;
    }

    /**
     * Tries to retrieve the specified tag from this element. If it's not found, `undefined` is returned.
     * If multiple tags of the same name exist, the first is returned.
     * @param name The name of the tag to retrieve.
     * @returns The {@link Tag} of the specified name, or `undefined` if no tag with that name exists.
     */
    public tryTag(name: string): Tag | undefined {
        return this.tags.find(a => a.name == name);
    }

    /**
     * Retrieves the specified tag from this element. If it's not found, an error is thrown.
     * @param name The name of the tag to retrieve.
     * @returns The {@link Tag} of the specified name.
     * @throws {@link Error} Thrown if no tag matching the specified name is found.
     */
    public tag(name: string): Tag {
        const tag = this.tryTag(name);
        if (!tag) throw new Error(`No tag named '${name}' was found on element '${this.name}'.`);
        return tag;
    }

    /**
     * Sets the specified tag to the specified value. If the tag doesn't already exist on this element,
     * a new one will be inserted as the last tag.
     * @param name The name of the tag to set.
     * @returns The {@link Tag} that was updated or inserted.
     */
    public setTag(name: string): Tag {
        let tag = this.tryTag(name);
        if (!tag) {
            tag = new Tag({name, parent: this});
            this.tags.push(tag)
        }
        return tag;
    }

    /**
     * Removes the specified tag from the element, if it exists.
     * @param name The name of the tag to remove.
     * @returns The {@link Tag} that was removed, or `undefined` if no tag with that name was found.
     */
    public removeTag(name: string): Tag | undefined {
        let tag = this.tryTag(name);
        if (tag) {
            this.tags = this.tags.filter(t => t !== tag);
            tag.parent = undefined;
        }
        return tag;
    }
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
 * Sets each of `children`'s `parent` and `model` properties to `parent`
 * and its `model`.
 */
export function setParentAndModel(parent: Element, children: IChild[]) {
    children.forEach(child => (child.parent = parent));
    setModel(parent.model, ...children);
}

/** Sets each of `children`'s model to `model`. */
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
