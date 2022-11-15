import {
    Annotation,
    Element,
    Endpoint,
    Field,
    ILocational,
    Model,
    Param,
    ParentElement,
    Statement,
    Tag,
    Type,
} from "../model";
import { Application } from "../model/application";

export function allItems(model: Model): ILocational[] {
    const items: ILocational[] = [];
    const listener = new AnyWalkListener((item: ILocational) =>
        items.push(item)
    );
    walk(model, listener);
    return items;
}

/**
 * Receives a callback (if implemented) for each element of a matching type over the walk of a
 * {@link Model}.
 */
export interface WalkListener {
    visitApp?: (app: Application) => void;
    visitAppAnno?: (anno: Annotation) => void;
    visitAppTag?: (tag: Tag) => void;
    visitType?: (type: Type) => void;
    visitTypeAnno?: (anno: Annotation) => void;
    visitTypeTag?: (tag: Tag) => void;
    visitField?: (field: Field) => void;
    visitFieldAnno?: (anno: Annotation) => void;
    visitFieldTag?: (tag: Tag) => void;
    visitEndpoint?: (endpoint: Endpoint) => void;
    visitEndpointAnno?: (anno: Annotation) => void;
    visitEndpointTag?: (tag: Tag) => void;
    visitParam?: (param: Param) => void;
    visitStatement?: (statement: Statement) => void;
    visitStatementAnno?: (anno: Annotation) => void;
    visitStatementTag?: (tag: Tag) => void;
}

export class AnyWalkListener implements WalkListener {
    constructor(private readonly visitAny: (item: ILocational) => void) {}

    visitApp(app: Application): void {
        this.visitAny(app);
    }
    visitAppAnno(anno: Annotation): void {
        this.visitAny(anno);
    }
    visitAppTag(tag: Tag): void {
        this.visitAny(tag);
    }
    visitType(type: Type): void {
        this.visitAny(type);
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
    visitEndpoint(endpoint: Endpoint): void {
        this.visitAny(endpoint);
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
    visitStatement(statement: Statement): void {
        this.visitAny(statement);
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
        listener.visitApp?.(app);
        listener.visitAppAnno &&
            app.annos.forEach(listener.visitAppAnno.bind(listener));
        listener.visitAppTag &&
            app.tags.forEach(listener.visitAppTag.bind(listener));

        app.children.forEach(element => {
            if (element instanceof Type) {
                listener.visitType?.(element);
                listener.visitTypeAnno &&
                    element.annos.forEach(
                        listener.visitTypeAnno.bind(listener)
                    );
                listener.visitTypeTag &&
                    element.tags.forEach(listener.visitTypeTag.bind(listener));
            }
            if (element instanceof ParentElement) {
                (element.children as Element[]).forEach(childElement => {
                    if (childElement instanceof Field) {
                        listener.visitField?.(childElement);
                        listener.visitFieldAnno &&
                            childElement.annos.forEach(
                                listener.visitFieldAnno.bind(listener)
                            );
                        listener.visitFieldTag &&
                            childElement.tags.forEach(
                                listener.visitFieldTag.bind(listener)
                            );
                    }
                });
            }
        });

        app.endpoints.forEach(ep => {
            listener.visitEndpoint?.(ep);
            listener.visitEndpointAnno &&
                ep.annos.forEach(listener.visitEndpointAnno.bind(listener));
            listener.visitEndpointTag &&
                ep.tags.forEach(listener.visitEndpointTag.bind(listener));

            ep.params.forEach(p => {
                listener.visitParam?.(p);
            });

            const visitStatements = (stmts: Statement[]): void => {
                stmts.forEach(stmt => {
                    listener.visitStatement?.(stmt);
                    listener.visitStatementAnno &&
                        stmt.annos.forEach(
                            listener.visitStatementAnno.bind(listener)
                        );
                    listener.visitStatementTag &&
                        stmt.tags.forEach(
                            listener.visitStatementTag.bind(listener)
                        );
                    visitStatements(stmt.children);
                });
            };
            visitStatements(ep.statements);
        });
    });
}
