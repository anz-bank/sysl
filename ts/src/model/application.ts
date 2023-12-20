import { IChild } from "./common";
import { ElementID, ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams, IParentElement } from "./element";
import { Endpoint } from "./endpoint";
import { Type } from "./type";
import { Enum } from "./enum";
import { Union } from "./union";
import { Alias } from "./alias";

export type AppChild = Type | Endpoint | Enum | Union | Alias;

export class Application extends Element implements IParentElement<AppChild> {
    namespace: readonly string[];
    children: AppChild[];

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

    override get safeName(): string {
        return this.toRef().toSysl();
    }

    override attachSubitems(extraSubitems: IChild[] = []): void {
        super.attachSubitems([...this.endpoints, ...extraSubitems]);
    }

    get types(): readonly Type[] {
        return this.children.filter(c => c instanceof Type) as Type[];
    }

    get endpoints(): readonly Endpoint[] {
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

    override toDto() {
        return {
            ...super.toDto(),
            namespace: this.namespace,
            children: this.children.map(e => e.toDto()),
        };
    }

    static fromDto(dto: ReturnType<Application["toDto"]>): Application {
        return new Application(new ElementRef(dto.namespace, dto.name), {
            ...Element.paramsFromDto(dto),
            children: dto.children.map(Application.fromChildDto),
        });
    }

    private static fromChildDto(dto: ReturnType<Element["toDto"]>): AppChild {
        // prettier-ignore
        switch (dto.kind) {
            case "Type":     return Type.fromDto(dto as ReturnType<Type["toDto"]>);
            case "Endpoint": return Endpoint.fromDto(dto as ReturnType<Endpoint["toDto"]>);
            case "Enum":     return Enum.fromDto(dto as ReturnType<Enum["toDto"]>);
            case "Union":    return Union.fromDto(dto as ReturnType<Union["toDto"]>);
            case "Alias":    return Alias.fromDto(dto as ReturnType<Alias["toDto"]>);
            default:         throw new Error(`Unknown app child kind '${dto.kind}'.`);
        }
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

export type ApplicationParams = IElementParams<void> & {
    namespace?: string[];
    children?: AppChild[];
};
