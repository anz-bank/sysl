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

export interface IStatementParams extends IElementParams<Endpoint | ParentStatement> {}
export type StatementKind = "ActionStatement" | "ReturnStatement" | "CallStatement" | "GroupStatement" |
                            "OneOfStatement" | "CondStatement" | "ForEachStatement" | "LoopStatement";

export abstract class Statement extends Element {
    public override get parent(): Endpoint | ParentStatement | undefined {
        return super.parent as Endpoint | ParentStatement;
    }
    public override set parent(epOrStatement: Endpoint | ParentStatement | undefined) {
        super.parent = epOrStatement;
    }

    constructor(p: IStatementParams) {
        super("", p.locations ?? [], p.annos ?? [], p.tags ?? [], p.model, p.parent);
        this.name = this.constructor.name;
        this.attachSubitems();
    }

    toRef(): ElementRef {
        const parentRef = this.parent!.toRef();
        const index = this.parent!.children.indexOf(this);
        if (index < 0) throw new Error(`Statement '${this.name}' not found in parent '${parentRef}'.`);
        return parentRef.with({ statementIndices: [...parentRef.statementIndices, index] });
    }

    toSysl(): string {
        return this.toString();
    }

    static fromDto(dto: ReturnType<Element["toDto"]>): Statement {
        switch (dto.kind) {
            case "ActionStatement":
                return Statement.create<ActionStatement>(dto, (dto, p) => new ActionStatement(dto.action, p));
            case "ReturnStatement":
                return Statement.create<ReturnStatement>(dto, (dto, p) => new ReturnStatement(dto.payload, p));
            case "CallStatement":
                return Statement.create<CallStatement>(dto, (dto, p) => new CallStatement(
                    ElementRef.parse(dto.targetEndpoint), dto.args, ElementRef.parse(dto.sourceApp), p));
            case "GroupStatement":
                return ParentStatement.create<GroupStatement>(dto, (dto, p) => new GroupStatement(dto.title, p));
            case "OneOfStatement":
                return ParentStatement.create<OneOfStatement>(dto, (_dto, p) => new OneOfStatement(p));
            case "CondStatement":
                return ParentStatement.create<CondStatement>(dto, (dto, p) => new CondStatement(dto.title, p));
            case "ForEachStatement":
                return ParentStatement.create<ForEachStatement>(dto, (dto, p) => new ForEachStatement(dto.title, p));
            case "LoopStatement":
                return ParentStatement.create<LoopStatement>(dto, (dto, p) => new LoopStatement(
                    dto.title, dto.prefix, p));
            default:
                throw new Error(`Unknown statement kind: ${dto.kind}`);
        }
    }

    static create<T extends Statement>(
        dto: ReturnType<Element["toDto"]>,
        factory: (dto: ReturnType<T["toDto"]>, p: ReturnType<typeof Element["paramsFromDto"]>) => T
    ): T {
        const params = this.getParams(dto);
        return factory(dto as ReturnType<T["toDto"]>, params);
    }

    static getParams(dto: ReturnType<Element["toDto"]>): ReturnType<typeof Element["paramsFromDto"]> {
        return Element.paramsFromDto(dto);
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
        public args: string[],
        public sourceApp: ElementRef,
        p: IStatementParams = {}) {
        super(p);
    }

    override toSysl(): string {
        const app = this.sourceApp.appsEqual(this.targetEndpoint)
            ? ElementRef.CurrentApp
            : this.targetEndpoint.truncate(ElementKind.App);
        return `${app.toSysl(false)} <- ${this.targetEndpoint.endpointName}` +
            `${this.args.length ? " (" + this.args.join(", ") + ")" : ""}`;
    }

    override toString(): string {
        return `${this.targetEndpoint.truncate(ElementKind.App)} <- ${this.targetEndpoint.endpointName}`;
    }

    override toDto() {
        return {
            ...super.toDto(),
            targetEndpoint: this.targetEndpoint.toString(),
            args: this.args,
            sourceApp: this.sourceApp.toString(),
        };
    }

    clone(context = new CloneContext(this.model)): CallStatement {
        return new CallStatement(this.targetEndpoint, [...this.args], this.sourceApp, this.cloneParams(context));
    }
}

export interface IParentStatementParams extends IStatementParams { children?: Statement[] };
export type ParamsWithChildren = ReturnType<typeof Element["paramsFromDto"]> & { children: Statement[] };

export abstract class ParentStatement extends Statement implements IStatementParams {
    private _children: Statement[];
    public get children(): Statement[] { return this._children ?? []; }
    public set children(value: Statement[]) { this._children = value; }

    constructor(public prefix: string, public title: string, p: IParentStatementParams = {}) {
        super(p);
        this._children = p.children ?? [];
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
            children: this.children.map(s => s.toDto()),
        };
    }

    static override getParams(dto: ReturnType<ParentStatement["toDto"]>): ParamsWithChildren {
        return { ...Element.paramsFromDto(dto), children: dto.children.map(s => Statement.fromDto(s)) };
    }

    protected override cloneParams(context: CloneContext): IParentStatementParams {
        return { children: context.recurse(this.children), ...super.cloneParams(context) };
    }
}

export class GroupStatement extends ParentStatement {
    constructor(title: string, p: IParentStatementParams = {}) { super("", title, p); }
    override toDto() { return { ...super.toDto(), title: this.title }; }
    clone(context = new CloneContext(this.model)): GroupStatement {
        return new GroupStatement(this.title, this.cloneParams(context));
    }
}

export class OneOfStatement extends ParentStatement {
    override get children(): GroupStatement[] {
        return super.children.filter(c => c instanceof GroupStatement) as GroupStatement[];
    }

    override set children(value: GroupStatement[]) {
        if (value.every(c => c instanceof GroupStatement)) super.children = value;
        throw new Error("All children of OneOfStatements must be a GroupStatement.");
    }

    constructor(p: IParentStatementParams = {}) {
        super("one of", "", p);
    }

    clone(context = new CloneContext(this.model)): OneOfStatement {
        return new OneOfStatement(super.cloneParams(context));
    }
}

export class CondStatement extends ParentStatement {
    // TODO: Also parse `else` and `else if` statements (which are currently treated as group statements).
    constructor(predicate: string, p: IParentStatementParams = {}) { super("", predicate, p); }
    clone(context = new CloneContext(this.model)): CondStatement {
        return new CondStatement(this.title, this.cloneParams(context));
    }
}

export class ForEachStatement extends ParentStatement {
    constructor(collection: string, p: IParentStatementParams = {}) { super("for each", collection, p); }
    clone(context = new CloneContext(this.model)): ForEachStatement {
        return new ForEachStatement(this.title, this.cloneParams(context));
    }
}

export class LoopStatement extends ParentStatement {
    constructor(criterion: string, mode: string, p: IParentStatementParams = {}) { super(mode, criterion, p); }
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
            method: this.method.toString() as keyof typeof RestMethod,
            path: this.path,
            queryParams: this.queryParams.map(p => p.toDto()),
            urlParams: this.urlParams.map(p => p.toDto()),
        };
    }

    static fromDto(dto: ReturnType<RestParams["toDto"]>): RestParams {
        return new RestParams({
            method: RestMethod[dto.method],
            path: dto.path,
            queryParams: dto.queryParams.map(p => Field.fromDto(p)),
            urlParams: dto.urlParams.map(p => Field.fromDto(p)),
        });
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
