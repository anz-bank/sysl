import { indent, toSafeName } from "../common/format";
import { IChild } from "./common";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams, IParentElement } from "./element";
import { Field } from "./field";
import { RestParams, Statement, ellipsis } from "./statement";
import { Application } from "./application";

export type EndpointParams = IElementParams<Application> & {
    longName?: string | undefined;
    params?: Field[];
    restParams?: RestParams | undefined;
    children?: Statement[];
    isPubsub?: boolean;
};

export class Endpoint extends Element implements IParentElement<Statement> {
    longName: string | undefined;
    params: Field[];
    children: Statement[];
    restParams: RestParams | undefined;
    isPubsub: boolean;

    public override get parent(): Application | undefined {
        return super.parent as Application;
    }
    public override set parent(app: Application | undefined) {
        super.parent = app;
    }

    constructor(name: string, p: EndpointParams = {}) {
        super(name, p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
        this.longName = p.longName;
        this.params = p.params ?? [];
        this.children = p.children ?? [];
        this.restParams = p.restParams;
        this.isPubsub = p.isPubsub ?? false;
        this.attachSubitems();
    }

    protected override attachSubitems(extraSubitems: IChild[] = []): void {
        super.attachSubitems([...this.params, ...extraSubitems]);
    }

    private renderQueryParams(): string {
        if (!this.restParams!.queryParams.length) return "";
        return `?${this.restParams!.queryParams.map(p => p.toRestParam()).join("&")}`;
    }

    private renderParams(): string {
        const renderParam = (p: Field) => {
            if (p.annos.length)
                throw new Error(
                    `Inline annotation rendering is not supported, so cannot render annotations` +
                        ` for param '${p.name}' on endpoint '${this.toRef()}'.`
                );
            return p.toSysl(true);
        };
        return this.params.length ? `(${this.params.map(renderParam).join(", ")})` : "";
    }

    private renderRestPath(): string {
        return this.restParams!.path.split("/")
            .map(s => {
                const match = s.match(`^{(\\w+)}$`);
                if (match) {
                    let name = match[1];
                    let param = this.restParams?.urlParams.find(p => p.name === name)?.toSysl();
                    return `{${param}}`;
                }
                return toSafeName(s);
            })
            .join("/");
    }

    private renderRestEndpoint(): string {
        const params = this.params.length ? ` ${this.renderParams()}` : "";
        const statements = this.children.length ? this.children.map(s => s.toSysl()).join("\n") : "...";
        let body = `${this.restParams!.method}${params}${this.renderQueryParams()}:\n${indent(statements)}`;
        return this.render("", body, this.renderRestPath());
    }

    private renderRpcEndpoint(): string {
        if (this.name === ellipsis) return this.name;
        const longName = this.longName ? `"${this.longName}"` : "";
        return this.render("", this.children, [this.name, longName, this.renderParams()].filter(s => s).join(" "));
    }

    toRef(): ElementRef {
        return this.parent!.toRef().with({ endpointName: this.name });
    }

    toSysl(): string {
        return this.restParams ? this.renderRestEndpoint() : this.renderRpcEndpoint();
    }

    override toDto() {
        return {
            ...super.toDto(),
            longName: this.longName,
            isPubsub: this.isPubsub,
            children: this.children.map(e => e.toDto()),
            params: this.params.map(p => p.toDto()),
            restParams: this.restParams?.toDto(),
        };
    }

    static fromDto(dto: ReturnType<Endpoint["toDto"]>): Endpoint {
        return new Endpoint(dto.name, {
            longName: dto.longName,
            isPubsub: dto.isPubsub,
            children: dto.children.map(e => Statement.fromDto(e)),
            params: dto.params.map(p => Field.fromDto(p)),
            restParams: dto.restParams ? RestParams.fromDto(dto.restParams) : undefined,
            ...Endpoint.paramsFromDto(dto),
        });
    }

    override toString(): string {
        return this.restParams ? `[REST] ${this.restParams.path}` : `[RPC] ${this.name}`;
    }

    clone(context = new CloneContext(this.model)): Endpoint {
        return new Endpoint(this.name, {
            longName: this.longName,
            isPubsub: this.isPubsub,
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            params: context.recurse(this.params),
            restParams: context.applyUnder(this.restParams),
            children: context.recurse(this.children),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        });
    }
}
