import "reflect-metadata";
import { Location } from "../location";
import { indent, joinedAppName } from "../util";
import { ComplexType } from "./common";
import { Element } from "./element";
import { renderInlineSections } from "./renderers";
import { Type, TypeValue } from "./type";

export class Param {
    name: string;
    type: Type | undefined;

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

    constructor(endpoint: string, arg: CallArg[], targetApp: string[], originApp: string[]) {
        this.endpoint = endpoint;
        this.arg = arg;
        this.targetApp = targetApp;
        this.originApp = originApp;
    }

    toSysl(): string {
        const joinedTarget = joinedAppName(this.targetApp)
        const appName = joinedAppName(this.originApp) === joinedTarget ? '.' : joinedTarget;
        return `${appName} <- ${this.endpoint}${this.arg ? this.arg.select(a => a.name).toArray().join(",") : ''}`;
    }
}

export class LoopN {
    count: number;
    stmt: Statement[];

    constructor(count: number, stmt: Statement[]) {
        this.count = count;
        this.stmt = stmt;
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
        return sysl += this.stmt.select(s => `\n${indent(s.toSysl())}`).toArray().join("");
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
        return `${this.cond ? "" + this.cond : ''}:${this.stmt.select(s => `\n${indent(s.toSysl())}`).toArray().join("")}`
    }
}

export class Alt {
    choice: AltChoice[];

    constructor(choice: AltChoice[]) {
        this.choice = choice;
    }

    toSysl(): string {
        let sysl = `one of:${this.choice.select(
            c => `\n${indent(c.toSysl())}`)
            .toArray().join("")}`;
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
        return sysl += this.stmt.select(s => `\n${indent(s.toSysl())}`).toArray().join("");
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
        return sysl += this.stmt.select(s => `\n${indent(s.toSysl())}`).toArray().join("");
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
        let sysl = `${this.mode.toString().toLowerCase()}${this.criterion ? ' ' + this.criterion : ''}:`;
        return sysl += this.stmt.select(s => `\n${indent(s.toSysl())}`).toArray().join("");
    }
}

export enum LoopMode {
    NOMode = 0,
    WHILE = 1,
    UNTIL = 2,
    UNRECOGNIZED = -1,
}

export class Statement extends ComplexType {
    // action: Action | undefined;
    // call: Call | undefined;
    // cond: Cond | undefined;
    // loop: Loop | undefined;
    // loopN: LoopN | undefined;
    // foreach: Foreach | undefined;
    // alt: Alt | undefined;
    // group: Group | undefined;
    // ret: Return | undefined;
    // locations: Location[];
    // tags: Tag[];
    // annotations: Annotation[];

    // constructor(action: Action | undefined,
    //     call: Call | undefined,
    //     cond: Cond | undefined,
    //     loop: Loop | undefined,
    //     loopN: LoopN | undefined,
    //     foreach: Foreach | undefined,
    //     alt: Alt | undefined,
    //     group: Group | undefined,
    //     ret: Return | undefined,
    //     locations: Location[],
    //     tags: Tag[],
    //     annotations: Annotation[]) {
    //     this.action = action
    //     this.call = call
    //     this.cond = cond
    //     this.loop = loop
    //     this.loopN = loopN
    //     this.foreach = foreach
    //     this.alt = alt
    //     this.group = group
    //     this.ret = ret
    //     this.locations = locations
    //     this.tags = tags
    //     this.annotations = annotations
    // }

    override toSysl(): string {
        let sysl = ``;
        // if (this.action) {
        //     sysl += this.action.toSysl();
        // }
        // else if (this.call) {
        //     sysl += this.call.toSysl();
        // }
        // else if (this.cond) {
        //     sysl += this.cond.toSysl();
        // }
        // else if (this.loop) {
        //     sysl += this.loop.toSysl();
        // }
        // else if (this.foreach) {
        //     sysl += this.foreach.toSysl()
        // }
        // else if (this.alt) {
        //     sysl += this.alt.toSysl();
        // }
        // else if (this.group) {
        //     sysl += this.group.toSysl();
        // }
        // else if (this.ret) {
        //     sysl += this.ret.toSysl();
        // }
        // else {
        //     sysl += "...";
        // }
        return sysl;
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

    constructor(method: RestMethod, path: string, queryParam: Param[], urlParam: Param[]) {
        this.method = method;
        this.path = path;
        this.queryParam = queryParam;
        this.urlParam = urlParam;
    }
}

export class Endpoint extends ComplexType {
    longName: string | undefined;
    docstring: string | undefined;
    flag: string[] | undefined;
    isPubsub: boolean;
    param: Param[];
    statements: Element<Statement>[];
    restParams: RestParams | undefined;
    source: string[];

    constructor(name: string,
        longName: string | undefined,
        docstring: string | undefined,
        flag: string[] | undefined,
        isPubsub: boolean,
        param: Param[],
        statements: Element<Statement>[],
        restParams: RestParams | undefined,
        locations: Location[],
        source: string[]) {
        super("", locations, name)
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
    }

    private renderRestEndpoint(): string {
        const segments = this.name.split("/");
        return segments.select(s => {
            if (s.match(`^{(\\w+)}$`)?.any) {
                let name = s.replace("{", "").replace("}", "");
                return `{${this.restParams?.urlParam?.single(p => p.name === name).toSysl()}}`;
            }
            return s;
        }).toArray().join("/")
    }

    private renderParams(): string {
        return this.param.any() ? `(${this.param.select(p => p.toSysl()).toArray().join(", ")})` : "";
    }

    private renderEndpoint(): string {
        if (this.name === "...") return this.name;
        const longName = this.longName ? `"${this.longName}"` : '';

        const sections = [this.name, longName, this.renderParams()];
        let sysl = `${renderInlineSections(sections)}:`;

        if (this.statements) {
            sysl += this.statements.select(s => `\n${indent(s.toSysl())}`).toArray().join("");
        }
        return sysl;
    }

    override toSysl(): string {
        return this.restParams ? this.renderRestEndpoint() : this.renderEndpoint();
    }
}
