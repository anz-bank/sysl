import { Location } from "../common";
import { indent, toSafeName } from "../common/format";
import { Annotation, AnnoValue, Tag } from "./attribute";
import { ElementRef, IChild, ILocational, IRenderable } from "./common";
import { CloneContext, ICloneable } from "./clone";
import { Model } from "./model";
import { addTags, renderAnnos } from "./renderers";

/**
 * An object in a Sysl model that can have nested objects (children, annotations and tags).
 */
export abstract class Element implements ILocational, IRenderable, IChild, ICloneable {
    constructor(
        public name: string,
        public locations: Location[],
        public annos: Annotation[],
        public tags: Tag[],
        public model?: Model,
        public parent?: Element
    ) {}

    abstract toSysl(): string;
    abstract toString(): string;

    public get safeName(): string {
        return toSafeName(this.name);
    }

    /**
     * Returns an {@link ElementRef} that references this element.
     */
    abstract toRef(): ElementRef;

    /**
     * Tries to find the specified annotation from this element. If it's not found, `undefined` is returned.
     * If multiple annotations of the same name exist, the first is returned.
     * @param name The name of the annotation to retrieve.
     * @returns The {@link Annotation} of the specified name, or `undefined` if no annotation with that name exists.
     */
    public findAnno(name: string): Annotation | undefined {
        return this.annos.find(a => a.name == name);
    }

    /**
     * Retrieves the specified annotation from this element. If it's not found, an error is thrown.
     * @param name The name of the annotation to retrieve.
     * @returns The {@link Annotation} of the specified name.
     * @throws {@link Error} Thrown if no annotation matching the specified name is found.
     */
    public getAnno(name: string): Annotation {
        const anno = this.findAnno(name);
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
        let anno = this.findAnno(name);

        if (value == undefined) {
            if (anno) {
                this.annos = this.annos.filter(a => a !== anno);
                anno.parent = undefined;
            }
        } else {
            if (anno) {
                anno.value = value;
            } else {
                anno = new Annotation(name, value);
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
                this.annos.splice(i + 1, 0, anno);
                return i + 1;
            }
        }

        this.annos.splice(0, 0, anno);
        return 0;
    }

    /**
     * Tries to find the specified tag from this element. If it's not found, `undefined` is returned.
     * If multiple tags of the same name exist, the first is returned.
     * @param name The name of the tag to retrieve.
     * @returns The {@link Tag} of the specified name, or `undefined` if no tag with that name exists.
     */
    public findTag(name: string): Tag | undefined {
        return this.tags.find(a => a.name == name);
    }

    /**
     * Retrieves the specified tag from this element. If it's not found, an error is thrown.
     * @param name The name of the tag to retrieve.
     * @returns The {@link Tag} of the specified name.
     * @throws {@link Error} Thrown if no tag matching the specified name is found.
     */
    public getTag(name: string): Tag {
        const tag = this.findTag(name);
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
        let tag = this.findTag(name);
        if (!tag) {
            tag = new Tag(name, { parent: this });
            this.tags.push(tag);
        }
        return tag;
    }

    /**
     * Removes the specified tag from the element, if it exists.
     * @param name The name of the tag to remove.
     * @returns The {@link Tag} that was removed, or `undefined` if no tag with that name was found.
     */
    public removeTag(name: string): Tag | undefined {
        let tag = this.findTag(name);
        if (tag) {
            this.tags = this.tags.filter(t => t !== tag);
            tag.parent = undefined;
        }
        return tag;
    }

    public abstract clone(context?: CloneContext): Element;

    /**
     * Ensures the `.parent` and `.model` properties of this instance and its model are set for all subitems: `annos`,
     * `tags`, `children` and any additional subitems specified.
     * @param extraSubitems Additional children to ensure attachment.
     */
    protected attachSubitems(extraSubitems: IChild[] = []) {
        [
            ...this.annos,
            ...this.tags,
            ...((this as unknown as ParentElement<Element>).children ?? []),
            ...extraSubitems,
        ].forEach(child => {
            child.parent = this;
            child.model = this.model;
        });
    }

    protected render(prefix: string, body: string | IRenderable[], name?: string, mustHaveBody: boolean = true) {
        if (prefix) prefix += " ";
        if (Array.isArray(body)) body = body.map(r => r.toSysl()).join("\n");
        if (mustHaveBody && !body && !this.annos.length) body = "...";
        const isExpanded = body || this.annos.length;
        let sysl = `${addTags(`${prefix}${name ?? this.safeName}`, this.tags)}${isExpanded ? ":" : ""}`;
        if (this.annos.length) sysl += `\n${indent(renderAnnos(this.annos))}`;
        return `${sysl}${body ? "\n" + indent(body) : ""}`;
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
