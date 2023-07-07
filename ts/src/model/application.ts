import { ElementID, ElementKind, ElementRef, IChild } from "./common";
import { CloneContext } from "./clone";
import { Element, IElementParams, ParentElement } from "./element";
import { Endpoint } from "./statement";
import { Type } from "./type";

export class Application extends ParentElement<Element> {
    namespace: readonly string[];
    endpoints: Endpoint[];
    children: Element[];

    constructor(name: ElementID, p: ApplicationParams = {}) {
        if (name instanceof ElementRef) {
            if (p.namespace) throw new Error("If namespace is specified, it must be a simple string, not ElementRef.");
        }

        const parsed = typeof name == "string" ? ElementRef.parse(name) : name;
        super(parsed.appName, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model);

        if (parsed.kind != ElementKind.App) throw Error(`Expected name to be of kind 'App' but got '${parsed.kind}'.`);
        if (parsed.namespace.length && p.namespace) throw Error("Found namespace in both 'name' and 'namespace'.");

        this.namespace = [...(p.namespace ?? parsed.namespace)];
        this.endpoints = p.endpoints ?? [];
        this.children = p.types ?? [];
        this.attachSubitems();
    }

    public override get safeName(): string {
        return this.toRef().toSysl();
    }

    protected override attachSubitems(extraSubitems: IChild[] = []): void {
        super.attachSubitems([...this.endpoints, ...extraSubitems]);
    }

    public get types(): Type[] {
        return this.children.filter(c => c instanceof Type) as Type[];
    }

    toSysl(): string {
        const endpoints = `${this.endpoints.filter(e => !e.isPubsub).map(e => e.toSysl()).join("\n\n")}`;
        const children = `${this.children.map(t => t.toSysl()).join("\n\n")}`;
        return this.render("", [endpoints, children].filter(x => x).join("\n"));
    }

    toRef(): ElementRef {
        return new ElementRef(this.namespace, this.name);
    }

    override toString(): string {
        return this.toRef().toSysl();
    }

    clone(context = new CloneContext(this.model)): Application {
        return new Application(this.toRef(), {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            endpoints: context.recurse(this.endpoints),
            types: context.recurse(this.children),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}

export type ApplicationParams = IElementParams & {
    namespace?: string[];
    endpoints?: Endpoint[];
    types?: Element[];
};
