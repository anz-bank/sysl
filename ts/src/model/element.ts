import { Location } from "../common/location";
import { indent, toSafeName } from "../common/format";
import { Annotation, AnnoValue, Tag } from "./attribute";
import { IChild, ILocational, IRenderable } from "./common";
import { ElementRef } from "./elementRef";
import { CloneContext, ICloneable } from "./clone";
import { Model } from "./model";
import { StatementKind } from "./statement";

type DtoKind = "Application" | "Type" | "Field" | "Endpoint" | "Enum" | "Union" | "Alias" | StatementKind;

/** An object in a Sysl model that can have metadata (annotations and tags). */
export abstract class Element implements ILocational, IRenderable, IChild, ICloneable {
    public get parent(): Element | undefined {
        return this._parent;
    }

    public set parent(p: Element | undefined) {
        this._parent = p;
    }

    constructor(
        public name: string,
        public locations: Location[],
        public annos: Annotation[],
        public tags: Tag[],
        public model?: Model,
        private _parent?: Element
    ) {}

    abstract toSysl(): string;
    abstract toString(): string;

    public get safeName(): string {
        return toSafeName(this.name);
    }

    /**
     * Converts this object to a simple Data Transfer Object (DTO) that makes serializing to other formats easy. The DTO
     * will not contain any circular references, will simplify certain key-based arrays into objects, and will represent
     * certain complex structures as strings.
     * @returns A simple DTO object that can be serialized.
     */
    public toDto() {
        return {
            kind: this.constructor.name as DtoKind,
            name: this.name,
            locations: Object.fromEntries([
                ...this.locations.map((l, i) => [i, l.toString()]),
                ...this.tags.map(t => [t.name, t.locations[0]?.toString()]),
                ...this.annos.map(a => [a.name, a.locations[0]?.toString()]),
            ]) as { [index: number]: string; [name: string]: string },
            metadata: Object.fromEntries([
                ...this.tags.map(t => [t.name, undefined]),
                ...this.annos.map(a => [a.name, a.value]),
            ]) as { [name: string]: AnnoValue | undefined },
        };
    }

    public static paramsFromDto(dto: ReturnType<Element["toDto"]>) {
        const locations: Location[] = [];
        let locStr: string;

        // Rehydrate all number-indexed locations that are found, which belong to the element itself (not to metadata).
        for (let i = 0; (locStr = dto.locations[i]); i++) locations.push(Location.parse(locStr));

        const metadata = Object.entries(dto.metadata).map(([name, value]) => {
            const loc = dto.locations[name];
            const params = { locations: loc ? [Location.parse(loc)] : [] };
            return value == undefined ? new Tag(name, params) : new Annotation(name, value, params);
        });
        const tags = metadata.filter(a => a instanceof Tag) as Tag[];
        const annos = metadata.filter(a => a instanceof Annotation) as Annotation[];

        return { locations, annos, tags };
    }

    /**
     * Returns an {@link ElementRef} that references this element.
     */
    abstract toRef(): ElementRef;

    /**
     * Tries to find the specified annotation from this element. If it's not found, `undefined` is returned. If multiple
     * annotations of the same name exist, the first is returned.
     * @param name The name of the annotation to retrieve.
     * @returns The {@link Annotation} of the specified name, or `undefined` if no annotation with that name exists.
     */
    public findAnno(name: string): Annotation | undefined {
        return this.annos.find(a => a.name == name);
    }

    /**
     * Retrieves the specified annotation from this element. If it's not found, an error is thrown. If multiple
     * annotations of the same name exist, the first is returned.
     * @param name The name of the annotation to retrieve.
     * @returns The {@link Annotation} of the specified name.
     * @throws {@link Error} Thrown if no annotation matching the specified name is found.
     */
    public getAnno(name: string): Annotation {
        const anno = this.findAnno(name);
        if (!anno) throw new Error(`No annotation named '${name}' was found on element '${this.toString()}'.`);
        return anno;
    }

    /**
     * Sets the specified annotation to the specified value. If the annotation doesn't already exist on this element,
     * a new one will be inserted with an best effort to place it alphabetically. If `undefined` is specified as the
     * {@link value}, all annotations with that name will be removed.
     * @param name The name of the annotation to set.
     * @param value The value of the annotation to set, or `undefined` to remove the annotation.
     * @returns The first {@link Annotation} that was updated, inserted or removed, or `undefined` if attempting to
     * remove an annotation that doesn't exist.
     */
    public setAnno(name: string, value: AnnoValue | undefined): Annotation | undefined {
        const anno = this.findAnno(name);

        if (value == undefined) {
            if (anno) {
                this.annos = this.annos.filter(a => a.name != name);
                anno.parent = undefined;
            }
        } else {
            if (anno) {
                anno.value = value;
            } else {
                const newAnno = new Annotation(name, value);
                this.insertAnnoOrdered(newAnno);
                return newAnno;
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
     * Determines if an element has a tag with one of the names specified.
     * @param names Names of the tag to look for.
     * @returns `true` if a tag with one of the specified names is present on the element, otherwise `false`.
     */
    public hasTag(...names: string[]): boolean {
        return this.tags.some(t => names.includes(t.name));
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
     * Retrieves the specified tag from this element. If it's not found, an error is thrown. If multiple tags with the
     * same name exist, the first will be returned.
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
            tag = new Tag(name, { parent: this, model: this.model });
            this.tags.push(tag);
        }
        return tag;
    }

    /**
     * Removes the all tags with the specified name from the element, if it exists.
     * @param name The name of the tag to remove.
     * @returns The first {@link Tag} that was removed, or `undefined` if no tag with that name was found.
     */
    public removeTag(name: string): Tag | undefined {
        const tag = this.findTag(name);
        if (tag) {
            this.tags = this.tags.filter(t => t.name != name);
            tag.parent = undefined;
        }
        return tag;
    }

    public abstract clone(context?: CloneContext): Element;

    /**
     * Ensures the `.parent` and `.model` properties of descendants are set. This includes all items in
     * {@link annos}, {@link tags} and {@link children}, recursively, and any {@link extraSubitems} specified. It does
     * not set properties on this instance itself, as they are expected to already be set correctly. The `.model`
     * property is set to the value of this instance, even if set to `undefined`.
     * @param extraSubitems Additional children to ensure attachment.
     */
    protected attachSubitems(extraSubitems: IChild[] = []) {
        const children = (this as unknown as IParentElement<Element>).children?.filter(c => c instanceof Element) ?? [];
        [...this.annos, ...this.tags, ...children, ...extraSubitems].forEach(child => {
            child.parent = this;
            child.model = this.model;
        });

        [...children, ...(extraSubitems.filter(i => i instanceof Element) as Element[])].forEach(c =>
            c.attachSubitems()
        );
    }

    protected render(prefix: string, body: string | IRenderable[], name?: string, mustHaveBody: boolean = true) {
        const annos = indent(this.annos.map(a => `@${a.toSysl()}`).join("\n"));
        if (prefix) prefix += " ";
        if (Array.isArray(body)) body = body.map(r => r.toSysl()).join("\n");
        if (mustHaveBody && !body && !annos) body = "...";
        const isExpanded = body || this.annos.length;
        const title = this.renderWithTags(`${prefix}${name ?? this.safeName}`, !!isExpanded);
        return [title, annos, indent(body)].filter(p => p).join("\n");
    }

    protected renderWithTags(part: string, colonSuffix: boolean = false): string {
        const hiddenTags = ["rest"];
        const tagList = this.tags
            .filter(t => !hiddenTags.includes(t.name))
            .map(t => t.toSysl())
            .join(", ");
        const withTags = tagList ? `${part} [${tagList}]` : part;
        return withTags + (colonSuffix ? ":" : "");
    }
}

/** An {@link Element} that also has child elements of type {@link TChild}. */
export interface IParentElement<TChild extends Element> extends Element {
    get children(): TChild[];
}

/** Common set of properties that are received by all {@link IElement} constructors. */
export type IElementParams<TParent extends Element | void> = {
    annos?: Annotation[];
    tags?: Tag[];
    parent?: TParent;
    locations?: Location[];
    model?: Model;
};
