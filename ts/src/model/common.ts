import { toSafeName, fromSafeName } from "../common/format";
import { Location } from "../common/location";
import { Element } from "./element";
import { Model } from "./model";

/** An element that can be serialized to Sysl source. */
export interface IRenderable {
    toSysl(): string;
}

/**
 * An element that has source context locating it in one or more source files.
 *
 * All objects with locations in Sysl source exist in a Sysl model. The {@link model} property is
 * a reference to that model.
 *
 * {@link model} may be undefined if an object has been created detached from a model.
 */
export interface ILocational {
    locations: Location[];
    model?: Model;
}

/**
 * An object that has a parent in the Sysl model.
 *
 * Follow the chain of {@link parent} properties should lead to an {@link Application} which will
 * have a falsey parent.
 *
 * {@link parent} can also be undefined if an object has been created detached from a model.
 */
export interface IChild extends ILocational {
    parent?: Element;
}

export class ElementRef implements IRenderable {
    readonly kind: ElementKind;

    constructor(
        public readonly namespace: readonly string[],
        public readonly appName: string,
        public readonly typeName: string = "",
        public readonly fieldName: string = ""
    ) {
        if (fieldName && !typeName)
            throw new Error("Cannot specify fieldName but omit typeName");

        if (fieldName) this.kind = ElementKind.Field;
        else if (typeName) this.kind = ElementKind.Type;
        else this.kind = ElementKind.App;
    }

    /**
     * Converts this element reference to a string that can be used as a value for annotations that reference other
     * elements.
     * @param compact Optional. True to remove spaces around the namespace delimiter (`::`), otherwise keep them.
     * Defaults to false.
     * @returns The formatted string representing the element reference. Unsafe characters are escaped.
     */
    toSysl(compact: boolean = false): string {
        const fullAppName = [...this.namespace, this.appName]
            .map(toSafeName)
            .join(compact ? "::" : " :: ");
        const typeName = toSafeName(this.typeName);
        const fieldName = toSafeName(this.fieldName);

        return [fullAppName, typeName, fieldName].filter(n => n).join(".");
    }

    /**
     * Parses a string into an {@link ElementRef}. The string can refer to an app, type or field. The element isn't guaranteed
     * to exist in any given model. Decodes escaped characters.
     * @param refStr The element reference string to parse.
     * @returns An {@link ElementRef} that references the element specified in the supplied string.
     * @throws {@link Error}
     * Thrown if the syntax of the supplied string is invalid.
     */
    static parse(refStr: string): ElementRef {
        const ref = this.tryParse(refStr);
        if (!ref) throw new Error(`Invalid element reference string: ${refStr}`);
        return ref;
    }

    /**
     * Tries to parse a string into an {@link ElementRef}. The string can refer to an app, type or field. The element isn't
     * guaranteed to exist in any given model. Decodes escaped characters.
     * @param refStr The element reference string to parse.
     * @returns An {@link ElementRef} that references the element specified in the supplied string, or `undefined` if the
     * the syntax of the supplied string is invalid.
     */
    static tryParse(refStr: string): ElementRef | undefined {
        const parts = refStr.split(".");

        if (parts.length > 3 || parts.some(p => !p))
            return undefined;

        const appNameParts: string[] = parts[0].split(/\s*::\s*/);
        return new ElementRef(
            appNameParts.slice(0, -1).map(fromSafeName),
            fromSafeName(appNameParts.at(-1)!),
            parts[1] ? fromSafeName(parts[1]) : undefined,
            parts[2] ? fromSafeName(parts[2]) : undefined
        );
    }

    /** 
     * Returns true if the supplied element reference refers to the same element as the current instance, otherwise
     * false.
     */
    equals(other: ElementRef): boolean {
        return (
            this.typesEqual(other) &&
            this.fieldName == other.fieldName
        );
    }

    /**
     * Returns true if the supplied element reference refers to the same app or type as the current instance
     * (ignoring the field), otherwise false.
     */
    typesEqual(other: ElementRef): boolean {
        return this.appsEqual(other) && this.typeName == other.typeName;
    }

    /**
     * Returns true if the supplied element reference refers to the same app (ignoring the type or field),
     * otherwise false.
     */
    appsEqual(other: ElementRef): boolean {
        return this.namespace.length == other.namespace.length &&
               this.namespace.every((_, i) => this.namespace[i] == other.namespace[i]) &&
               this.appName == other.appName;
    }

    /**
     * Clones the current instance with additional modifications to its parts.
     * @param parts An objects with optional properties to modify parts of the new element reference. If any property
     * is not specified, it will not be modified.
     * @returns A modified clone of the current instance.
     */
    with(parts: { namespace?: string[], appName?: string, typeName?: string, fieldName?: string }): ElementRef {
        return new ElementRef(parts.namespace ?? this.namespace, parts.appName ?? this.appName, parts.typeName ?? this.typeName, parts.fieldName ?? this.fieldName);
    }
}

/**
 * The kinds of Sysl element that an annotation can be applied to.
 */
export enum ElementKind {
    App = "app",
    Type = "type",
    Field = "field",
    Endpoint = "ep",
    Parameter = "param",
    Statement = "stmt",
}
