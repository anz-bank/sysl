import {
    Annotation,
    Element,
    Endpoint,
    Field,
    IChild,
    ILocational,
    Model,
    Param,
    Statement,
    Tag,
    Type,
} from "../model";
import { Application } from "../model/application";

export function allItems(model: Model): ILocational[] {
    const items: ILocational[] = [];
    const listener = new AnyWalkListener((item: ILocational) => items.push(item));
    walk(model, listener);
    return items;
}

/**
 * Receives a callback (if implemented) for each element of a matching type over the walk of a
 * {@link Model}.
 *
 * The callbacks {@link visitApp}, {@link visitType}, {@link visitField}, {@link visitEndpoint} and
 * {@link visitStatement} return a boolean that indicate whether the walking should continue into their children.
 * Return `true` to continue the walk as normal, or `false` to skip descending any further into the tree from that
 * element. The callback will not be called for any descendants of that element.
 */
export interface WalkListener {
    visitApp?: (app: Application) => boolean;
    visitAppAnno?: (anno: Annotation) => void;
    visitAppTag?: (tag: Tag) => void;
    visitType?: (type: Type) => boolean;
    visitTypeAnno?: (anno: Annotation) => void;
    visitTypeTag?: (tag: Tag) => void;
    visitField?: (field: Field) => void;
    visitFieldAnno?: (anno: Annotation) => void;
    visitFieldTag?: (tag: Tag) => void;
    visitEndpoint?: (endpoint: Endpoint) => boolean;
    visitEndpointAnno?: (anno: Annotation) => void;
    visitEndpointTag?: (tag: Tag) => void;
    visitParam?: (param: Param) => void;
    visitStatement?: (statement: Statement) => boolean;
    visitStatementAnno?: (anno: Annotation) => void;
    visitStatementTag?: (tag: Tag) => void;
}

export class AnyWalkListener implements WalkListener {
    constructor(private readonly visitAny: (item: ILocational) => void) {}

    visitApp(app: Application): boolean {
        this.visitAny(app);
        return true;
    }
    visitAppAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitAppTag(tag: Tag): void {
        this.visitAny(tag);
    }
    visitType(type: Type): boolean {
        this.visitAny(type);
        return true;
    }
    visitTypeAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitTypeTag(tag: Tag): void {
        this.visitAny(tag);
    }
    visitField(field: Field): void {
        this.visitAny(field);
    }
    visitFieldAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitFieldTag(tag: Tag): void {
        this.visitAny(tag);
    }
    visitEndpoint(endpoint: Endpoint): boolean {
        this.visitAny(endpoint);
        return true;
    }
    visitEndpointAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitEndpointTag(tag: Tag): void {
        this.visitAny(tag);
    }
    visitParam(param: Param): void {
        this.visitAny(param);
    }
    visitParamAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitParamTag(tag: Tag): void {
        this.visitAny(tag);
    }
    visitStatement(statement: Statement): boolean {
        this.visitAny(statement);
        return true;
    }
    visitStatementAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitStatementTag(tag: Tag): void {
        this.visitAny(tag);
    }
}

/** Walks a model tree and invokes callbacks on {@link listener} for each kind of content. */
export function walk(model: Model, listener: WalkListener) {
    model.apps.forEach(app => {
        fillModel(app, model);
        if (listener.visitApp?.(app) == false) return;
        visitMeta(app, listener.visitAppAnno?.bind(listener), listener.visitAppTag?.bind(listener));

        app.children.forEach(element => {
            fillParent(element, app);
            if (element instanceof Type) {
                if (listener.visitType?.(element) == false) return;
                visitMeta(element, listener.visitTypeAnno?.bind(listener), listener.visitTypeTag?.bind(listener));

                element.children.forEach(field => {
                    fillParent(field, element);
                    listener.visitField?.(field);
                    visitMeta(field, listener.visitFieldAnno?.bind(listener), listener.visitFieldTag?.bind(listener));
                });
            }
        });

        app.endpoints.forEach(ep => {
            fillParent(ep, app);
            if (listener.visitEndpoint?.(ep) == false) return;
            visitMeta(ep, listener.visitEndpointAnno?.bind(listener), listener.visitEndpointTag?.bind(listener));

            ep.params.forEach(p => {
                fillModel(p, model);
                listener.visitParam?.(p);
            });

            const visitStatements = (stmts: Statement[], parent: Element): void => {
                stmts.forEach(stmt => {
                    fillParent(stmt, parent);
                    if (listener.visitStatement?.(stmt) == false) return;
                    visitMeta(
                        stmt,
                        listener.visitStatementAnno?.bind(listener),
                        listener.visitStatementTag?.bind(listener)
                    );
                    visitStatements(stmt.children, stmt);
                });
            };
            visitStatements(ep.statements, ep);
        });
    });
}

function visitMeta(element: Element, annoVisitor?: (anno: Annotation) => void, tagVisitor?: (tag: Tag) => void) {
    for (const anno of element.annos) {
        fillParent(anno, element);
        annoVisitor?.(anno);
    }

    for (const tag of element.tags) {
        fillParent(tag, element);
        tagVisitor?.(tag);
    }
}

function fillParent(child: IChild, parent: Element): void {
    fillModel(child, parent.model!);
    if (child.parent == parent) return;
    if (!child.parent) {
        child.parent = parent;
    } else {
        let subject = `a fragment of type '${child.constructor.name}'`;
        if (child instanceof Element) subject = `an element '${child.toRef().toSysl(true)}'`;
        else if (child instanceof Annotation) subject = `an anno '${child?.name}'`;
        else if (child instanceof Tag) subject = `a tag '${child?.name}'`;
        throw new Error(
            `Detected ${subject} with an incorrect parent. It has the parent ` +
                `'${child.parent.toRef().toSysl(true)}' but is a child of '${parent.toRef().toSysl(true)}.`
        );
    }
}

function fillModel(locational: ILocational, model: Model): void {
    if (locational.model == model) return;
    if (!locational.model) {
        locational.model = model;
    } else {
        // TODO: Make filterByFile() not incorporate items from a different model, then uncomment below checks.
        //
        // if (locational instanceof Element) {
        //     throw new Error(`Detected an element with an incorrect model. Element ` +
        //     `'${locational.toRef().toSysl(true)}' has a different model to the one being walked.`);
        // } else {
        //     throw new Error(`Detected a fragment (${locational.constructor.name}) with an incorrect model. Fragment ` +
        //     `at location '${locational.locations[0].toString()}' has a different model to the one being walked.`);
        // }
    }
}
