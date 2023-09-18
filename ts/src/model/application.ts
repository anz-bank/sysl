import { ElementID, ElementKind, ElementRef, IChild } from "./common";
import { CloneContext } from "./clone";
import { Element, IElementParams, ParentElement } from "./element";
import { Endpoint } from "./statement";
import { Type } from "./type";

export class Application extends ParentElement<Element> {
    namespace: readonly string[];
    children: Element[];

    constructor(name: ElementID, p: ApplicationParams = {}) {
        if (name instanceof ElementRef) {
            if (p.namespace) throw new Error("If namespace is specified, it must be a simple string, not ElementRef.");
        }

        const parsed = typeof name == "string" ? ElementRef.parse(name) : name;
        super(parsed.appName, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model);

        if (!parsed.isApp) throw Error(`Expected name to be of kind 'App' but got '${parsed.kind}'.`);
        if (parsed.namespace.length && p.namespace) throw Error("Found namespace in both 'name' and 'namespace'.");

        this.namespace = [...(p.namespace ?? parsed.namespace)];
        this.children = p.children ?? [];
        this.attachSubitems();
    }

    public override get safeName(): string {
        return this.toRef().toSysl();
    }

    public override attachSubitems(extraSubitems: IChild[] = []): void {
        super.attachSubitems([...this.endpoints, ...extraSubitems]);
    }

    public get types(): readonly Type[] {
        return this.children.filter(c => c instanceof Type) as Type[];
    }

    public get endpoints(): readonly Endpoint[] {
        return this.children.filter(c => c instanceof Endpoint) as Endpoint[];
    }

    toSysl(): string {
        return this.render(
            "",
            this.children
                .filter(x => !(x as Endpoint).isPubsub) // Pubsub endpoint rendering is not yet implemented
                .map(t => t.toSysl())
                .join("\n\n")
        );
    }

    public override toDto() {
        return { ...super.toDto(), namespace: this.namespace };
    }

    toRef(): ElementRef {
        return new ElementRef(this.namespace, this.name);
    }

    override toString(): string {
        return this.toRef().toString();
    }

    clone(context = new CloneContext(this.model)): Application {
        return new Application(this.toRef(), {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            children: context.recurse(this.children),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}

export type ApplicationParams = IElementParams & {
    namespace?: string[];
    children?: Element[];
};
