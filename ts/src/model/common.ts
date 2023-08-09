import { toSafeName, fromSafeName, isSafeName } from "../common/format";
import { Location } from "../common/location";
import { CloneContext } from "./clone";
import { Element } from "./element";
import { Model } from "./model";
import * as util from "util";

// TODO: Move below three interfaces to src/common/interfaces.ts
/** An element that can be serialized to Sysl source. */
export interface IRenderable {
    toSysl(): string;
}

/**
 * An element that has source context locating it in one or more source files.
 *
 * All objects with locations in Sysl source exist in a Sysl model. The {@link model} property is a reference to that
 * model.
 *
 * {@link model} may be undefined if an object has been created detached from a model. After adding detached objects
 * to a model, call {@link Model.attachSubitems()} to populate this field.
 */
export interface ILocational {
    locations: Location[];
    model?: Model;
}

/**
 * An object that has a parent in the Sysl model.
 *
 * Following the chain of {@link parent} properties should eventually lead to an {@link Application}, which will have a
 * parent value of `undefined`.
 *
 * {@link parent} on non-applications may be undefined if an object has been created detached from a model. After adding
 * detached objects to a model, call {@link Model.attachSubitems()} to populate this field.
 */
export interface IChild extends ILocational {
    parent?: Element;
}

// TODO: Move to src/common/elementRef.ts
/**
 * Represents a reference to an Element. This element may or may not exist in a model, and does not point to any
 * actual instance of an element. Rather, it is used to reference an element by it's unique path in the hierarchy, in a
 * similar way to how XPath is used. This includes specifying the app's namespace and name, and optionally the type and
 * field names. It currently only supports referencing apps, types and fields.
 */
export class ElementRef implements IRenderable {
    readonly kind: ElementKind = ElementKind.App;

    /**
     * Creates a new instance of {@link ElementRef}, with the {@link kind} determined using the supplied arguments.
     * @param namespace The namespace parts of the Element's application, or an empty array if the application is in the
     * root namespace. All parts of the namespace must be non-empty.
     * @param appName The application's name.
     * @param typeName Optional. The type's name. Must be specified if {@link fieldName} is specified.
     * @param fieldName Optional. The field's name.
     * @returns An instance of the {@link ElementRef} class.
     * @throws `Error` if any part of {@link namespace} is empty.
     * @throws `Error` if an empty {@link appName} was specified.
     * @throws `Error` if an empty {@link fieldName} was specified without specifying {@link typeName}.
     */
    constructor(
        public readonly namespace: readonly string[],
        public readonly appName: string,
        public readonly typeName: string = "",
        public readonly fieldName: string = ""
    ) {
        if (namespace.some(n => !n)) throw new Error("Empty namespace parts are invalid: " + util.inspect(namespace));
        if (!appName) throw new Error("appName must not be empty.");
        if (fieldName && !typeName) throw new Error("typeName must be specified if fieldName is specified");

        if (fieldName) this.kind = ElementKind.Field;
        else if (typeName) this.kind = ElementKind.Type;
    }

    /** True is this reference's {@link kind} is {@link ElementKind.App}, otherwise false. */
    public get isApp(): boolean {
        return this.kind == ElementKind.App;
    }

    /** True is this reference's {@link kind} is {@link ElementKind.Type}, otherwise false. */
    public get isType(): boolean {
        return this.kind == ElementKind.Type;
    }

    /** True is this reference's {@link kind} is {@link ElementKind.Field}, otherwise false. */
    public get isField(): boolean {
        return this.kind == ElementKind.Field;
    }

    /** The references kind represented numerically. @see {@link ElementKind.toNumeric()} */
    public get numericKind(): number {
        return ElementKind.toNumeric(this.kind);
    }

    /**
     * Converts this element reference to a string that can be used as a value for annotations that reference other
     * elements.
     * @param compact Optional. True to remove spaces around the namespace delimiter (`::`), otherwise keep them.
     * Defaults to false.
     * @param parentRef Optional. If specified and if it has the same app as the current instance, will produce a string
     * of a relative reference (omitting the app part), unless it's an app reference.
     * @returns The formatted string representing the element reference. Unsafe characters are escaped.
     */
    toSysl(compact: boolean = false, parentRef?: ElementRef): string {
        const fullAppName = this.toAppPartsSafe().join(compact ? "::" : " :: ");
        const typeName = toSafeName(this.typeName);
        const fieldName = toSafeName(this.fieldName);
        let parts = [fullAppName, typeName, fieldName];
        if (!this.isApp && parentRef?.appsEqual(this)) parts.shift();
        return parts.filter(n => n).join(".");
    }

    /** Returns a string representation of the reference, with compact formatting. */
    toString(): string {
        return this.toSysl(true);
    }

    /**
     * Converts this element reference to an array of string representing the **unescaped** application name parts,
     * usually delimited by `::`. For example, the `ElementRef` referencing the field
     * `Company::Billing%26Invoice::Service.Customers.City` will return the string array
     * `[ "Company", "Billing&Invoice", "Service" ]`. If you are producing text that will go in a Sysl file, use
     * {@link toAppPartsSafe()} instead.
     * @returns An array of strings having the name parts of the app portion this instance refers to.
     */
    toAppParts(): string[] {
        return [...this.namespace, this.appName];
    }

    /**
     * Converts this element reference to an array of string representing the **escaped** application name parts,
     * usually delimited by `::`. For example, the `ElementRef` referencing the field
     * `Company::Billing%26Invoice::Service.Customers.City` will return the string array
     * `[ "Company", "Billing%26Invoice", "Service" ]`. If you require the text for display purposes, use
     * {@link toAppParts()} instead.
     * @returns An array of strings having the name parts of the app portion this instance refers to.
     */
    toAppPartsSafe(): string[] {
        return [...this.namespace, this.appName].map(toSafeName);
    }

    /**
     * Create an application reference to an application from array of string representing the **unescaped** application
     * name parts.
     * @param parts An array of strings having the name parts (namespace and appName) of the app.
     * @returns A new `ElementRef` instance with the namespace and appName filled in according to the provided parts.
     */
    static fromAppParts(parts: string[]): ElementRef {
        if (!parts.length) throw new Error("At least one part must be specified to create an ElementRef.");
        return new ElementRef(parts.slice(0, -1), parts.at(-1)!);
    }

    /**
     * Create an application reference to an application from array of string representing the **escaped** application
     * name parts.
     * @param parts An array of strings having the name parts (namespace and appName) of the app.
     * @returns A new `ElementRef` instance with the namespace and appName filled in according to the provided parts.
     */
    static fromAppPartsSafe(parts: string[]): ElementRef {
        return this.fromAppParts(parts.map(fromSafeName));
    }

    /**
     * Returns an reference to the parent of the current reference. All namespace parts are considered intermediate
     * applications, so the last part of a namespace will be returned as the parent when an application reference is
     * provided.
     * @returns An {@link ElementRef} to the parent of the current reference.
     * @throws `Error` if the current reference is an Application in the root namespace, which does not have a parent.
     */
    toParent(): ElementRef {
        const parent = this.toParentOrSelf();
        if (parent === this) throw new Error("Cannot get parent reference for an Application");
        return parent;
    }

    /**
     * Returns an reference to the parent of the current reference, or the current reference if it's an application in
     * the root namespace. All namespace parts are considered intermediate applications, so the last part of a namespace
     * will be returned as the parent when an application reference is provided.
     * @returns An {@link ElementRef} to the parent/self of the current reference.
     */
    toParentOrSelf(): ElementRef {
        if (this.isApp) return this.namespace.length ? ElementRef.fromAppParts(this.toAppParts().slice(0, -1)) : this;

        return this.truncate(ElementKind.fromNumeric(this.numericKind - 1));
    }

    /**
     * Parses a string into an {@link ElementRef}. The string can refer to an app, type or field. The element isn't
     * guaranteed to exist in any given model. Decodes escaped characters.
     * @param refStr The element reference string to parse.
     * @returns An {@link ElementRef} that references the element specified in the supplied string.
     * @throws `Error` Thrown if the syntax of the supplied string is invalid.
     */
    static parse(refStr: string): ElementRef {
        const ref = this.tryParse(refStr);
        if (!ref) throw new Error(`Invalid element reference string: ${refStr}`);
        return ref;
    }

    /**
     * Tries to parse a string into an {@link ElementRef}. The string can refer to an app, type or field. The element
     * isn't guaranteed to exist in any given model. Decodes escaped characters.
     * @param refStr The element reference string to parse.
     * @returns An {@link ElementRef} that references the element specified in the supplied string, or `undefined` if
     * the the syntax of the supplied string is invalid.
     */
    static tryParse(refStr: string): ElementRef | undefined {
        const parts = refStr.trim().split(/\s*\.\s*/);
        const appNameParts: string[] = parts[0].split(/\s*::\s*/);

        if (parts.length > 3 || [...appNameParts, ...parts.slice(1)].some(p => !isSafeName(p))) return undefined;

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
        return this.typesEqual(other) && this.fieldName == other.fieldName;
    }

    /**
     * Returns true if the supplied element reference refers to the same app and type as the current instance
     * (ignoring the field), otherwise false. If both references are to applications, the types are ignored.
     */
    typesEqual(other: ElementRef): boolean {
        return this.appsEqual(other) && this.typeName == other.typeName;
    }

    /**
     * Returns true if the supplied element reference refers to the same app (ignoring the type or field), otherwise
     * false.
     */
    appsEqual(other: ElementRef): boolean {
        return (
            this.namespace.length == other.namespace.length &&
            this.namespace.every((_, i) => this.namespace[i] == other.namespace[i]) &&
            this.appName == other.appName
        );
    }

    /**
     * Determines if the current reference is a descendant of the supplied reference. All namespace parts are considered
     * intermediate applications that qualify as ancestors.
     * @param ancestor The ancestor reference to test the descendancy of the current reference against.
     * @returns True if the current reference is a descendant of the supplied ancestor reference, otherwise false.
     */
    isDescendantOf(ancestor: ElementRef): boolean {
        let isDeeper = this.numericKind > ancestor.numericKind;
        if (ancestor.isApp) {
            const descendantParts = this.toAppParts();
            const ancestorParts = ancestor.toAppParts();
            if (this.isApp) isDeeper = descendantParts.length > ancestorParts.length;
            return isDeeper && ancestorParts.every((p, i) => descendantParts[i] == p); // Check app starts the same
        } else {
            return isDeeper && this.typesEqual(ancestor);
        }
    }

    /**
     * Clones the current instance with additional modifications to its parts.
     * @param parts An objects with optional properties to modify parts of the new element reference. If any property
     * is not specified (or specified with `undefined`), it will not be modified. If an empty string/array is specified,
     * it will remove that part of the reference.
     * @returns An {@link ElementRef} combining the current instances and the part changes requested.
     */
    with(parts: { namespace?: string[]; appName?: string; typeName?: string; fieldName?: string }): ElementRef {
        return new ElementRef(
            parts.namespace ?? this.namespace,
            parts.appName ?? this.appName,
            parts.typeName ?? this.typeName,
            parts.fieldName ?? this.fieldName
        );
    }

    /**
     * Truncates an ElementKind to a certain {@link ElementKind} by removing the parts that aren't in that kind. If
     * a truncation is requested to a kind that contains the same part or to one that has more parts than the current
     * reference has, the same reference is returned.
     * @param kind The {@link ElementKind} to truncate the current reference into.
     * @returns A truncated reference, or the same reference if there is nothing to truncate.
     * @throws `Error` if an unsupported request kind is specified.
     */
    truncate(kind: ElementKind): ElementRef {
        switch (kind) {
            case ElementKind.App:
                return new ElementRef(this.namespace, this.appName);
            case ElementKind.Type:
                return new ElementRef(this.namespace, this.appName, this.typeName);
            case ElementKind.Field:
                return this;
            default:
                throw new Error("Requested kind is not supported.");
        }
    }

    /**
     * Returns the current instance, since ElementRef is immutable and doesn't need cloning.
     * @param _context Unused.
     * @returns The current instance.
     */
    clone(_context?: CloneContext): ElementRef {
        return this;
    }
}

// TODO: Move rest of the file to src/common/elementKind.ts
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

export namespace ElementKind {
    /**
     * Converts an {@link ElementKind} to a numerical form. Applications are 1, types and endpoints are 2, fields and
     * statements are 3.
     * @param kind The {@link ElementKind} for which to return the numeric form.
     * @returns The numerical form for the specified kind.
     */
    export function toNumeric(kind: ElementKind): number {
        if (kind == ElementKind.App) return 1;
        if (kind == ElementKind.Type || kind == ElementKind.Endpoint) return 2;
        return 3;
    }

    const numericToKind: readonly ElementKind[] = [ElementKind.App, ElementKind.Type, ElementKind.Field];

    /**
     * Converts the numerical form of an {@link ElementKind} to their corresponding data-related instance: application
     * for 1, type for 2 and field for 3.
     * @param numericKind The numerical form of an {@link ElementKind}.
     * @returns The data-related {@link ElementKind} for the corresponding numerical form.
     * @throws `Error` if an invalid numerical form is specified.
     */
    export function fromNumeric(numericKind: number): ElementKind {
        const kind = numericToKind[numericKind - 1];
        if (!kind) throw new Error("Invalid depth. Must be between 1 and 3.");
        return kind;
    }
}

export type ElementID = ElementRef | string;
