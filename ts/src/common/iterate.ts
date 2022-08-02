import {
    Annotation,
    Application,
    Model,
    Struct,
    Tag,
    Type,
    TypeValue,
} from "../model";

function isStruct(value: TypeValue): value is Struct {
    return value.hasOwnProperty("elements");
}

export function forEachTypeField(type: Type, callback: (field: Type) => any) {
    if (!isStruct(type.value)) {
        return;
    }
    type.value.elements.forEach(callback);
}

export function forEachAppType(
    app: Application,
    callback: (type: Type) => any
) {
    app.types.forEach(callback);
}

export function forEachAppField(
    app: Application,
    callback: (field: Type) => any
) {
    forEachAppType(app, type => forEachTypeField(type, callback));
}

export function forEachFieldAnno(
    field: Type,
    callback: (anno: Annotation) => any
) {
    field.annos.forEach(callback);
}

// type Listener<T> = (element: T) => any;

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
    visitField?: (field: Type) => void;
    visitFieldAnno?: (anno: Annotation) => void;
    visitFieldTag?: (tag: Tag) => void;
}

export function walk(model: Model, listener: WalkListener) {
    model.apps.forEach(app => {
        listener.visitApp?.(app);
        listener.visitAppAnno &&
            app.annos.forEach(listener.visitAppAnno.bind(listener));
        listener.visitAppTag &&
            app.tags.forEach(listener.visitAppTag.bind(listener));

        app.types.forEach(type => {
            listener.visitType?.(type);
            listener.visitTypeAnno &&
                type.annos.forEach(listener.visitTypeAnno.bind(listener));
            listener.visitTypeTag &&
                type.tags.forEach(listener.visitTypeTag.bind(listener));

            if (isStruct(type.value)) {
                type.value.elements.forEach(field => {
                    listener.visitField?.(field);
                    listener.visitFieldAnno &&
                        field.annos.forEach(
                            listener.visitFieldAnno.bind(listener)
                        );
                    listener.visitFieldTag &&
                        field.tags.forEach(
                            listener.visitFieldTag.bind(listener)
                        );
                });
            }
        });
    });
}
