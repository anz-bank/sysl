import "reflect-metadata";
import { Location } from "../location";
import { indent, joinedAppName } from "../format";
import { IDescribable, ILocational, IRenderable } from "./common";
import { renderAnnos, renderInlineSections, addTags } from "./renderers";
import { Primitive, Reference, Type, TypeValue } from "./type";
import { Annotation, Tag } from "./attribute";

export class Param {
    name: string;
    type: Primitive | Reference | undefined;

    constructor(name: string, type: Type | undefined) {
        this.name = name;
        this.type = type;
    }

    toSysl(): string {
        let sysl = `${this.name}`;
        if (this.type) {
            sysl += ` <: ${this.type.toSysl()}`;
        }
        return sysl;
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

export class CallArg {
    value: TypeValue | undefined;
    name: string;

    constructor(value: TypeValue | undefined, name: string) {
        this.value = value;
        this.name = name;
    }
}

export class Call {
    endpoint: string;
    arg: CallArg[];
    targetApp: string[];
    originApp: string[];

    constructor(
        endpoint: string,
        arg: CallArg[],
        targetApp: string[],
        originApp: string[]
    ) {
        this.endpoint = endpoint;
        this.arg = arg;
        this.targetApp = targetApp;
        this.originApp = originApp;
    }

    toSysl(): string {
        const joinedTarget = joinedAppName(this.targetApp);
        const appName =
            joinedAppName(this.originApp) === joinedTarget ? "." : joinedTarget;
        return `${appName} <- ${this.endpoint}${
            this.arg
                ? this.arg
                      .select(a => a.name)
                      .toArray()
                      .join(",")
                : ""
        }`;
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
        return (sysl += this.stmt
            .select(s => `\n${indent(s.toSysl())}`)
            .toArray()
            .join(""));
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
        return `${this.cond ? "" + this.cond : ""}:${this.stmt
            .select(s => `\n${indent(s.toSysl())}`)
            .toArray()
            .join("")}`;
    }
}

export class Alt {
    choice: AltChoice[];

    constructor(choice: AltChoice[]) {
        this.choice = choice;
    }

    toSysl(): string {
        let sysl = `one of:${this.choice
            .select(c => `\n${indent(c.toSysl())}`)
            .toArray()
            .join("")}`;
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
        return (sysl += this.stmt
            .select(s => `\n${indent(s.toSysl())}`)
            .toArray()
            .join(""));
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
        return (sysl += this.stmt
            .select(s => `\n${indent(s.toSysl())}`)
            .toArray()
            .join(""));
    }
}

export class Loop {
    mode: LoopMode;
    criterion: string | undefined;
    stmt: Statement[];

    constructor(
        mode: LoopMode,
        criterion: string | undefined,
        stmt: Statement[]
    ) {
        this.mode = mode;
        this.criterion = criterion;
        this.stmt = stmt;
    }

    toSysl(): string {
        let sysl = `${this.mode.toString().toLowerCase()}${
            this.criterion ? " " + this.criterion : ""
        }:`;
        return (sysl += this.stmt
            .select(s => `\n${indent(s.toSysl())}`)
            .toArray()
            .join(""));
    }
}

export enum LoopMode {
    NOMode = 0,
    WHILE = 1,
    UNTIL = 2,
    UNRECOGNIZED = -1,
}

export class Statement implements IDescribable, ILocational, IRenderable {
    value:
        | Action
        | Call
        | Cond
        | Loop
        | LoopN
        | Foreach
        | Alt
        | Group
        | Return
        | undefined;
    locations: Location[];
    tags: Tag[];
    annos: Annotation[];

    constructor(
        value:
            | Action
            | Call
            | Cond
            | Loop
            | LoopN
            | Foreach
            | Alt
            | Group
            | Return
            | undefined,
        locations: Location[],
        tags: Tag[],
        annos: Annotation[]
    ) {
        this.value = value;
        this.locations = locations;
        this.tags = tags;
        this.annos = annos;
    }

    toSysl(): string {
        return this.value?.toSysl() ?? "...";
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

export class RestParams {
    method: RestMethod;
    path: string;
    queryParam: Param[];
    urlParam: Param[];

    constructor(
        method: RestMethod,
        path: string,
        queryParam: Param[],
        urlParam: Param[]
    ) {
        this.method = method;
        this.path = path;
        this.queryParam = queryParam;
        this.urlParam = urlParam;
    }
}

export class Endpoint implements IDescribable, ILocational, IRenderable {
    name: string;
    longName: string | undefined;
    docstring: string | undefined;
    flag: string[] | undefined;
    isPubsub: boolean;
    param: Param[];
    statements: Statement[];
    restParams: RestParams | undefined;
    source: string[];
    locations: Location[];
    tags: Tag[];
    annos: Annotation[];

    constructor(
        name: string,
        longName: string | undefined,
        docstring: string | undefined,
        flag: string[] | undefined,
        isPubsub: boolean,
        param: Param[],
        statements: Statement[],
        restParams: RestParams | undefined,
        source: string[],
        locations: Location[],
        tags: Tag[],
        annos: Annotation[]
    ) {
        this.locations = locations;
        this.name = name;
        this.longName = longName;
        this.docstring = docstring;
        this.flag = flag;
        this.isPubsub = isPubsub;
        this.param = param;
        this.statements = statements;
        this.restParams = restParams;
        this.source = source;
        this.locations = locations;
        this.tags = tags;
        this.annos = annos;
    }

    private renderRestEndpoint(): string {
        const segments = this.name.split("/");
        return segments
            .select(s => {
                if (s.match(`^{(\\w+)}$`)?.any) {
                    let name = s.replace("{", "").replace("}", "");
                    return `{${this.restParams?.urlParam
                        ?.single(p => p.name === name)
                        .toSysl()}}`;
                }
                return s;
            })
            .toArray()
            .join("/");
    }

    private renderParams(): string {
        return this.param.any()
            ? `(${this.param
                  .select(p => p.toSysl())
                  .toArray()
                  .join(", ")})`
            : "";
    }

    private renderEndpoint(): string {
        if (this.name === "...") return this.name;
        const longName = this.longName ? `"${this.longName}"` : "";

        const sections = [this.name, longName, this.renderParams()];
        let sysl = `${addTags(renderInlineSections(sections), this.tags)}:`;
        if (this.annos.any()) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        if (this.statements) {
            sysl += this.statements
                .select(s => `\n${indent(s.toSysl())}`)
                .toArray()
                .join("");
        }
        return sysl;
    }

    toSysl(): string {
        return this.restParams
            ? this.renderRestEndpoint()
            : this.renderEndpoint();
    }
}
