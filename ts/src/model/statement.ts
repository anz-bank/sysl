import "reflect-metadata";
import { Location } from "../common";
import { indent, joinedAppName } from "../common/format";
import { ElementRef, IChild, ILocational, IRenderable } from "./common";
import { CollectionDecorator } from "./decorator";
import { Element, IElementParams, ParentElement } from "./element";
import { Field } from "./field";
import { GenericElement } from "./genericElement";
import { Model } from "./model";
import { addTags, renderAnnos, renderInlineSections } from "./renderers";

/** Name of an endpoint that represents the absence of endpoints in an {@link Application}. */
const placeholder = "...";

export type ParamParams = IElementParams & {
    name: string;
    element?: Element;
};

export class Param implements ILocational, IRenderable {
    constructor(public name: string, public locations: Location[], public element?: Element, public model?: Model) {}

    toSysl(): string {
        return `${this.name}${this.element ? ` <: ${this.element.toSysl()}` : ""}`;
    }
}

export class Action {
    action: string;

    constructor(action: string) {
        this.action = action;
    }

    toSysl(): string {
        return this.action;
    }
}

export type ValueType = boolean | number | string;

export class CallArg {
    name: string;
    value: ValueType | undefined;

    constructor(name: string, value?: ValueType) {
        this.value = value;
        this.name = name;
    }
}

export type CallParams = {
    endpoint: string;
    arg?: CallArg[];
    targetApp: string[];
    originApp: string[];
};

export class Call {
    endpoint: string;
    arg: CallArg[];
    targetApp: string[];
    originApp: string[];

    constructor({ endpoint, arg, targetApp, originApp }: CallParams) {
        this.endpoint = endpoint;
        this.arg = arg ?? [];
        this.targetApp = targetApp;
        this.originApp = originApp;
    }

    toSysl(): string {
        const joinedTarget = joinedAppName(this.targetApp);
        const appName = joinedAppName(this.originApp) === joinedTarget ? "." : joinedTarget;
        return `${appName} <- ${this.endpoint}${this.arg ? this.arg.map(a => a.name).join(",") : ""}`;
    }
}

export class LoopN {
    count: number;
    stmt: Statement[];

    constructor(count: number, stmt: Statement[]) {
        this.count = count;
        this.stmt = stmt;
    }

    toSysl(): string {
        // TODO: implement or remove LoopN
        return "";
    }
}

export class Foreach {
    collection: string;
    stmt: Statement[];

    constructor(collection: string, stmt: Statement[]) {
        this.collection = collection;
        this.stmt = stmt;
    }

    toSysl(): string {
        let sysl = `for each ${this.collection}:`;
        return (sysl += this.stmt.map(s => `\n${indent(s.toSysl())}`).join(""));
    }
}

export class AltChoice {
    cond: string;
    stmt: Statement[];

    constructor(cond: string, stmt: Statement[]) {
        this.cond = cond;
        this.stmt = stmt;
    }

    toSysl(): string {
        return `${this.cond ? "" + this.cond : ""}:${this.stmt.map(s => `\n${indent(s.toSysl())}`).join("")}`;
    }
}

export class Alt {
    choice: AltChoice[];

    constructor(choice: AltChoice[]) {
        this.choice = choice;
    }

    toSysl(): string {
        let sysl = `one of:${this.choice.map(c => `\n${indent(c.toSysl())}`).join("")}`;
        return sysl;
    }
}

export class Group {
    title: string;
    stmt: Statement[];

    constructor(title: string, stmt: Statement[]) {
        this.title = title;
        this.stmt = stmt;
    }

    toSysl(): string {
        let sysl = `${this.title}:`;
        return (sysl += this.stmt.map(s => `\n${indent(s.toSysl())}`).join(""));
    }
}

export class Return {
    payload: string;

    constructor(payload: string) {
        this.payload = payload;
    }

    toSysl(): string {
        return `return ${this.payload}`;
    }
}

export class Cond {
    test: string;
    stmt: Statement[];

    constructor(test: string, stmt: Statement[]) {
        this.test = test;
        this.stmt = stmt;
    }

    toSysl(): string {
        let sysl = `if ${this.test}:`;
        return (sysl += this.stmt.map(s => `\n${indent(s.toSysl())}`).join(""));
    }
}

export class Loop {
    mode: LoopMode;
    criterion: string | undefined;
    stmt: Statement[];

    constructor(mode: LoopMode, criterion: string | undefined, stmt: Statement[]) {
        this.mode = mode;
        this.criterion = criterion;
        this.stmt = stmt;
    }

    toSysl(): string {
        let sysl = `${this.mode.toString().toLowerCase()}${this.criterion ? " " + this.criterion : ""}:`;
        return (sysl += this.stmt.map(s => `\n${indent(s.toSysl())}`).join(""));
    }
}

export enum LoopMode {
    NOMode = 0,
    WHILE = 1,
    UNTIL = 2,
    UNRECOGNIZED = -1,
}

export type StatementValue = Action | Call | Cond | Loop | LoopN | Foreach | Alt | Group | Return | undefined;

export type StatementParams = IElementParams & { value: StatementValue };

export class Statement extends ParentElement<Statement> {
    value: StatementValue;

    constructor({ value, annos, tags, locations, parent, model }: StatementParams) {
        super(value?.constructor.name ?? "", locations ?? [], annos ?? [], tags ?? [], model, parent);
        this.value = value;
        this.attachSubitems();
    }

    /** Returns a statement with the given action text. */
    static action(action: string): Statement {
        return new Statement({ value: new Action(action) });
    }

    /** Returns an array of child statements nested in this statement's {@link value}. */
    get children(): Statement[] {
        if (this.value && "stmt" in this.value) {
            return this.value.stmt;
        }
        return [];
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    toSysl(): string {
        return this.value?.toSysl() ?? placeholder;
    }
}

export enum RestMethod {
    NOMethod = "NOMethod",
    GET = "GET",
    PUT = "PUT",
    POST = "POST",
    DELETE = "DELETE",
    PATCH = "PATCH",
}

export type RestParamsParams = {
    method: RestMethod;
    path: string;
    queryParam?: Param[];
    urlParam?: Param[];
};

export class RestParams {
    method: RestMethod;
    path: string;
    queryParam: Param[];
    urlParam: Param[];

    constructor({ method, path, queryParam, urlParam }: RestParamsParams) {
        this.method = method;
        this.path = path;
        this.queryParam = queryParam ?? [];
        this.urlParam = urlParam ?? [];
    }
}

export type EndpointParams = IElementParams & {
    name: string;
    longName?: string | undefined;
    docstring?: string | undefined;
    params?: Param[];
    restParams?: RestParams | undefined;
    statements?: Statement[];
    isPubsub?: boolean;
    pubsubSource?: string[];
};

export class Endpoint extends ParentElement<Statement> {
    longName: string | undefined;
    docstring: string | undefined;
    params: Param[];
    statements: Statement[];
    restParams: RestParams | undefined;
    isPubsub: boolean;
    pubsubSource: string[];

    constructor({
        name,
        longName,
        docstring,
        isPubsub,
        params,
        statements,
        restParams,
        pubsubSource,
        annos,
        tags,
        locations,
        parent,
        model,
    }: EndpointParams) {
        super(name, locations ?? [], annos ?? [], tags ?? [], model, parent);
        this.longName = longName;
        this.docstring = docstring;
        this.params = params ?? [];
        this.statements = statements ?? [];
        this.restParams = restParams;
        this.isPubsub = isPubsub ?? false;
        this.pubsubSource = pubsubSource ?? [];

        this.attachSubitems(this.params);
    }

    protected override attachSubitems(extraSubitems: IChild[] = []): void {
        super.attachSubitems([...this.params, ...extraSubitems]);
    }

    private renderQueryParams(): string {
        if (!this.restParams!.queryParam.length) {
            return "";
        }

        return `?${this.restParams!.queryParam.map(p => {
            let s = `${p.name}=`;
            if (
                (p.element instanceof GenericElement || p.element instanceof Field) &&
                p.element?.value instanceof CollectionDecorator
            ) {
                return (s += `{${p.element.toSysl()}}`);
            }
            return (s += `${p.element?.toSysl() ?? "Type"}`);
        }).join("&")}`;
    }

    private renderParams(): string {
        return this.params.length ? `(${this.params.map(p => p.toSysl()).join(", ")})` : "";
    }

    private renderRestEndpoint(): string {
        const segments = this.restParams!.path.split("/");
        const path = segments
            .map(s => {
                const match = s.match(`^{(\\w+)}$`);
                if (match) {
                    let name = match[1];
                    let param = this.restParams?.urlParam.find(p => p.name === name)?.toSysl();
                    return `{${param}}`;
                }
                return s;
            })
            .join("/");
        let sysl = `${path}:`;
        let content = `${this.restParams!.method}${this.renderQueryParams()}${
            this.params.length ? ` ${this.renderParams()}` : ""
        }:`;
        content += this.statements.map(s => `\n${indent(s.toSysl())}`).join("");
        return (sysl += `\n${indent(content)}`);
    }

    private renderGRPCEndpoint(): string {
        if (this.name === placeholder) {
            return this.name;
        }
        const longName = this.longName ? `"${this.longName}"` : "";

        const sections = [this.name, longName, this.renderParams()];
        let sysl = `${addTags(renderInlineSections(sections), this.tags)}:`;
        if (this.annos.length) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        sysl += this.statements.map(s => `\n${indent(s.toSysl())}`).join("");
        return sysl;
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    toSysl(): string {
        return this.restParams ? this.renderRestEndpoint() : this.renderGRPCEndpoint();
    }

    get children(): Statement[] {
        return this.statements;
    }
}
