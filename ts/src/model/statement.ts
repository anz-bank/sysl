import "reflect-metadata";
import { indent } from "../common/format";
import { ElementRef } from "./elementRef";
import { CloneContext } from "./clone";
import { Element, IElementParams } from "./element";
import { Endpoint } from "./endpoint";
import { Field } from "./field";
import { ElementKind } from "./elementKind";

export const ellipsis = "...";
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

export interface IStatementParams extends IElementParams<Endpoint | ParentStatement> { };

export abstract class Statement extends Element {
    public override get parent(): Endpoint | ParentStatement | undefined { return super.parent as Endpoint | ParentStatement; }
    public override set parent(epOrStatement: Endpoint | ParentStatement | undefined) { super.parent = epOrStatement; }

    constructor(p: IStatementParams) {
        super("", p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
        this.name = this.constructor.name;
        this.attachSubitems();
    }

    toRef(): ElementRef {
        const parentRef = this.parent!.toRef();
        const index = this.parent!.children.indexOf(this);
        if (index < 0) throw new Error(`Statement '${this.name}' not found in parent '${parentRef}'.`)
        return parentRef.with({ statementIndices: [...parentRef.statementIndices, index] });
    }

    toSysl(): string {
        return this.toString();
    }

    protected cloneParams(context: CloneContext): IStatementParams {
        return {
            tags: context.recurse(this.tags),
            annos: context.recurse(this.annos),
            model: context.model ?? this.model,
            locations: context.keepLocation ? context.recurse(this.locations) : [],
        };
    }
}

export class ActionStatement extends Statement {
    constructor(public action: string, p: IStatementParams = {}) { super(p); }
    toString(): string { return this.action; }
    override toDto() { return { ...super.toDto(), action: this.action }; }
    clone(context = new CloneContext(this.model)): ActionStatement {
        return new ActionStatement(this.action, this.cloneParams(context));
    }
}

export class ReturnStatement extends Statement {
    constructor(public payload: string, p: IStatementParams = {}) { super(p); }
    toString(): string { return `return ${this.payload}`; }
    override toDto() { return { ...super.toDto(), payload: this.payload }; }
    clone(context = new CloneContext(this.model)): ReturnStatement {
        return new ReturnStatement(this.payload, this.cloneParams(context));
    }
}

export class CallStatement extends Statement {
    constructor(
        public targetEndpoint: ElementRef,
        public args: CallArg[],
        public sourceApp: ElementRef,
        p: IStatementParams = {}) {
        super(p);
    }

    override toSysl(): string {
        const app = this.sourceApp.appsEqual(this.targetEndpoint) ? ElementRef.CurrentApp : this.targetEndpoint.truncate(ElementKind.App);
        return `${app} <- ${this.targetEndpoint.endpointName}${this.args ? this.args.map(a => a.name).join(",") : ""}`;
    }

    override toString(): string {
        return `${this.targetEndpoint.truncate(ElementKind.App)} <- ${this.targetEndpoint.endpointName}`;
    }

    override toDto()
    {
        return {
            ...super.toDto(),
            action: this.targetEndpoint.toString(),
            args: Object.fromEntries(this.args.map(a => [a.name, a.value])),
            sourceApp: this.sourceApp.toString(),
        };
    }

    clone(context = new CloneContext(this.model)): CallStatement {
        return new CallStatement(
            this.targetEndpoint,
            context.recurse(this.args),
            this.sourceApp,
            this.cloneParams(context)
        );
    }
}

export interface IParentStatementParams extends IStatementParams { children?: Statement[] };

export abstract class ParentStatement extends Statement implements IStatementParams {
    private _children: Statement[];
    public get children(): Statement[] { return this._children ?? []; }
    public set children(value: Statement[]) { this._children = value; }

    constructor(public prefix: string, public title: string, children: Statement[] = [], p: IParentStatementParams = {}) {
        super(p);
        this._children = children;
    }

    override toString() {
        return [this.prefix, this.title].filter(s => s).join(" ");
    }
    
    override toSysl() {
        const childrenBlock = this.children.map(s => s.toSysl()).join("\n");
        return `${this.toString()}:\n${indent(childrenBlock || ellipsis)}`;
    }

    override toDto() {
        return {
            ...super.toDto(),
            prefix: this.prefix,
            title: this.title,
            children: this.children.map(s => s.toDto())
        };
    }

    protected override cloneParams(context: CloneContext): IParentStatementParams {
        return { children: context.recurse(this.children), ...super.cloneParams(context) };
    }
}

export class GroupStatement extends ParentStatement {
    constructor(title: string, p: IParentStatementParams = {}) { super("", title, p.children, p); }
    override toDto() { return { ...super.toDto(), title: this.title }; }
    clone(context = new CloneContext(this.model)): GroupStatement {
        return new GroupStatement(this.title, this.cloneParams(context));
    }
}

export class OneOfStatement extends ParentStatement {
    override get children() : GroupStatement[] {
        return super.children.filter(c => c instanceof GroupStatement) as GroupStatement[];
    }

    override set children(value: GroupStatement[]) {
        if (value.every(c => c instanceof GroupStatement)) super.children = value;
        throw new Error("All children of OneOfStatements must be a GroupStatement.");
    }

    constructor(cases: GroupStatement[], p: IParentStatementParams = {})
    {
        super("one of", "", cases, p);
    }

    clone(context = new CloneContext(this.model)): OneOfStatement {
        return new OneOfStatement(this.children, super.cloneParams(context));
    }
}

export class CondStatement extends ParentStatement {
    // TODO: Also parse `else` and `else if` statements (which are currently treated as group statements).
    constructor(predicate: string, p: IParentStatementParams = {}) { super("", predicate, p.children, p); }
    clone(context = new CloneContext(this.model)): CondStatement {
        return new CondStatement(this.title, this.cloneParams(context));
    }
}

export class ForEachStatement extends ParentStatement {
    constructor(collection: string, p: IParentStatementParams = {}) { super("for each", collection, p.children, p); }
    clone(context = new CloneContext(this.model)): ForEachStatement {
        return new ForEachStatement(this.title, this.cloneParams(context));
    }
}

export class LoopStatement extends ParentStatement {
    constructor(criterion: string, mode: string, p: IParentStatementParams = {}) { super(mode, criterion, p.children, p); }
    clone(context = new CloneContext(this.model)): LoopStatement {
        return new LoopStatement(this.title, this.prefix, this.cloneParams(context));
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

export type IRestParams = {
    method: RestMethod;
    path: string;
    queryParams?: Field[];
    urlParams?: Field[];
};

export class RestParams {
    method: RestMethod;
    path: string;
    queryParams: Field[];
    urlParams: Field[];

    constructor({ method, path, queryParams: queryParam, urlParams: urlParam }: IRestParams) {
        this.method = method;
        this.path = path;
        this.queryParams = queryParam ?? [];
        this.urlParams = urlParam ?? [];
    }

    toString() {
        return this.method;
    }

    toDto() {
        return {
            method: this.method.toString(),
            path: this.path,
            queryParams: this.queryParams.map(p => p.toDto()),
            urlParams: this.urlParams.map(p => p.toDto()),
        }
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

