import "reflect-metadata";
import { Location } from "../common";
import { indent, joinedAppName, toSafeName } from "../common/format";
import { ElementRef, IChild, ILocational, IRenderable } from "./common";
import { CloneContext } from "./clone";
import { CollectionDecorator } from "./decorator";
import { Element, IElementParams, ParentElement } from "./element";
import { Field } from "./field";
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

    toString(): string {
        return this.name;
    }

    clone(context = new CloneContext()): Param {
        return new Param(this.name, [], context.applyUnder(this.element), context.model ?? this.model);
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

    toString(): string {
        return this.toSysl();
    }

    clone(_context: CloneContext): Action {
        return new Action(this.action);
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

    clone(_context: CloneContext): CallArg {
        return new CallArg(this.name, this.value);
    }
}

export type CallParams = {
    endpoint: string;
    args?: CallArg[];
    targetApp: string[];
    originApp: string[];
};

export class Call {
    endpoint: string;
    args: CallArg[];
    targetApp: string[];
    originApp: string[];

    constructor({ endpoint, args: arg, targetApp, originApp }: CallParams) {
        this.endpoint = endpoint;
        this.args = arg ?? [];
        this.targetApp = targetApp;
        this.originApp = originApp;
    }

    toSysl(): string {
        const joinedTarget = joinedAppName(this.targetApp);
        const appName = joinedAppName(this.originApp) === joinedTarget ? "." : joinedTarget;
        return `${appName} <- ${this.endpoint}${this.args ? this.args.map(a => a.name).join(",") : ""}`;
    }

    toString(): string {
        return `${joinedAppName(this.targetApp, true)} <-- ${this.endpoint}`;
    }

    clone(context = new CloneContext()): Call {
        return new Call({
            endpoint: this.endpoint,
            args: context.recurse(this.args),
            targetApp: [...this.targetApp],
            originApp: [...this.originApp],
        });
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

    toString(): string {
        return `count=${this.count}`;
    }

    clone(context = new CloneContext()): LoopN {
        return new LoopN(this.count, context.recurse(this.stmt));
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

    toString(): string {
        return this.collection;
    }

    clone(context = new CloneContext()): Foreach {
        return new Foreach(this.collection, context.recurse(this.stmt));
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

    toString(): string {
        return this.cond;
    }

    clone(context = new CloneContext()): AltChoice {
        return new AltChoice(this.cond, context.recurse(this.stmt));
    }
}

export class Alt {
    choices: AltChoice[];

    constructor(choice: AltChoice[]) {
        this.choices = choice;
    }

    toSysl(): string {
        let sysl = `one of:${this.choices.map(c => `\n${indent(c.toSysl())}`).join("")}`;
        return sysl;
    }

    toString(): string {
        return `[${this.choices} choices]`;
    }

    clone(context = new CloneContext()): Alt {
        return new Alt(context.recurse(this.choices));
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

    toString(): string {
        return this.title;
    }

    clone(context = new CloneContext()): Group {
        return new Group(this.title, context.recurse(this.stmt));
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

    toString(): string {
        return this.toSysl();
    }

    clone(_context: CloneContext): Return {
        return new Return(this.payload);
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

    toString(): string {
        return this.test;
    }

    clone(context = new CloneContext()): Cond {
        return new Cond(this.test, context.recurse(this.stmt));
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

    toString(): string {
        return `[${this.mode}] ${this.criterion}`;
    }

    clone(context = new CloneContext()): Loop {
        return new Loop(this.mode, this.criterion, context.recurse(this.stmt));
    }
}

export enum LoopMode {
    NOMode = 0,
    WHILE = 1,
    UNTIL = 2,
    UNRECOGNIZED = -1,
}

export type StatementValue = Action | Call | Cond | Loop | LoopN | Foreach | Alt | Group | Return | undefined;

export class Statement extends ParentElement<Statement> {
    constructor(public value: StatementValue, { annos, tags, locations, parent, model }: IElementParams = {}) {
        super(value?.constructor.name ?? "", locations ?? [], annos ?? [], tags ?? [], model, parent);
        this.attachSubitems();
    }

    /** Returns a statement with the given action text. */
    static action(action: string): Statement {
        return new Statement(new Action(action));
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

    toString(): string {
        return this.value?.constructor.name ?? placeholder;
    }

    clone(context = new CloneContext(this.model)): Statement {
        return new Statement(context.applyUnder(this.value), {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
        });
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
    queryParams?: Param[];
    urlParams?: Param[];
};

export class RestParams {
    method: RestMethod;
    path: string;
    queryParams: Param[];
    urlParams: Param[];

    constructor({ method, path, queryParams: queryParam, urlParams: urlParam }: RestParamsParams) {
        this.method = method;
        this.path = path;
        this.queryParams = queryParam ?? [];
        this.urlParams = urlParam ?? [];
    }

    toString() {
        return this.method;
    }

    clone(context = new CloneContext()): RestParams {
        return new RestParams({
            method: this.method,
            path: this.path,
            queryParams: context.recurse(this.queryParams),
            urlParams: context.recurse(this.urlParams),
        });
    }
}

export type EndpointParams = IElementParams & {
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

    constructor(name: string, {
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
    }: EndpointParams = {}) {
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
        if (!this.restParams!.queryParams.length) {
            return "";
        }

        return `?${this.restParams!.queryParams.map(p => {
            let s = `${p.name}=`;
            if (p.element instanceof Field && p.element?.value instanceof CollectionDecorator) {
                return (s += `{${p.element.toSysl()}}`);
            }
            return (s += `${p.element?.toSysl() ?? "Type"}`);
        }).join("&")}`;
    }

    private renderParams(): string {
        return this.params.length ? `(${this.params.map(p => p.toSysl()).join(", ")})` : "";
    }

    private renderRestPath(): string {
        return this.restParams!.path
            .split("/")
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
        let sysl = this.renderRestPath();
        let content = "";
        content += this.restParams!.method;
        content += this.renderQueryParams();
        content += this.params.length ? ` ${this.renderParams()}` : "";
        content += ":";
        content += this.statements.map(s => `\n${indent(s.toSysl())}`).join("");
        return (sysl += `:\n${indent(content)}`);
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

    override toString(): string {
        return this.restParams ? `[REST] ${this.restParams.path}` : `[gRPC] ${this.name}`;
    }

    get children(): Statement[] {
        return this.statements;
    }

    clone(context = new CloneContext(this.model)): Endpoint {
        return new Endpoint(this.name, {
            longName: this.longName,
            docstring: this.docstring,
            isPubsub: this.isPubsub,
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            params: context.recurse(this.params),
            restParams: context.applyUnder(this.restParams),
            statements: context.recurse(this.statements),
            model: context.model ?? this.model,
        });
    }
}
