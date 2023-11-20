import { toSafeName, fromSafeName, isSafeName } from "../common/format";
import { ElementKind } from "./elementKind";
import { CloneContext } from "./clone";
import { IRenderable } from "./common";

/** A flexible way to specify the identity of an element. If a string is specified, it will be parsed. */
export type ElementID = ElementRef | string;

/**
 * Represents a reference to an Element. This element may or may not exist in a model, and does not point to any
 * actual instance of an element. Rather, it is used to reference an element by it's unique path in the hierarchy, in a
 * similar way to how XPath is used. This includes specifying the app's namespace and name, and optionally the type and
 * field names or endpoint and statement indices.
 */

export class ElementRef implements IRenderable {
    /** The kind element this this instance references. */
    readonly kind: ElementKind = ElementKind.App;

    /** Represents an reference to the current app */
    static readonly CurrentApp = new ElementRef([], ".");

    /**
     * Creates a new instance of {@link ElementRef}, with the {@link kind} determined using the supplied arguments.
     * @param namespace The namespace parts of the Element's application, or an empty array if the application is in the
     * root namespace. All parts of the namespace must be non-empty.
     * @param appName The application's name.
     * @param typeName Optional. The type's name. Must be specified if {@link fieldName} is specified.
     * @param fieldName Optional. The field's name.
     * @param endpointName Optional. The endpoint's name. Must be specified if {@link statementIndices} is specified.
     * @param statementIndices Optional. The statements positional zero-based indices for each level of nesting.
     * @returns An instance of the {@link ElementRef} class.
     * @throws `Error` if any part of {@link namespace} is empty.
     * @throws `Error` if an empty {@link appName} was specified.
     * @throws `Error` if an empty {@link fieldName} was specified without specifying {@link typeName}.
     * @throws `Error` if an empty {@link statementIndices} was specified without specifying {@link endpointName}.
     * @throws `Error` if specifying both data and behavior element parameters.
     */
    constructor(
        public readonly namespace: readonly string[],
        public readonly appName: string,
        public readonly typeName: string = "",
        public readonly fieldName: string = "",
        public readonly endpointName: string = "",
        public readonly statementIndices: readonly number[] = []
    ) {
        if (namespace.some(n => !n)) throw new Error("All namespace parts must be non-empty");
        if (!appName) throw new Error("appName must not be empty.");

        if (fieldName) this.kind = ElementKind.Field;
        else if (statementIndices.length) this.kind = ElementKind.Statement;
        else if (typeName) this.kind = ElementKind.Type;
        else if (endpointName) this.kind = ElementKind.Endpoint;

        switch (this.kind) {
            case ElementKind.Field:
                if (!typeName) throw new Error("typeName must be specified if fieldName is specified");
                if (endpointName || statementIndices.length)
                    throw new Error(
                        "endpointName or statementIndices must not be specified if fieldName is specified."
                    );
                break;
            case ElementKind.Statement:
                if (typeName || fieldName)
                    throw new Error(
                        "endpointName or statementIndices must not be specified if typeName or fieldName is specified."
                    );
                if (!endpointName) throw new Error("endpointName must be specified if statementIndices is specified.");
                if (statementIndices.some(i => i < 0 || isNaN(i)))
                    throw new Error("All values in statementIndices must be positive.");
                break;
            case ElementKind.Type:
                if (endpointName) throw new Error("endpointName must not be specified if typeName is specified.");
                break;
        }
    }

    /** True if this reference's {@link kind} is {@link ElementKind.App}, otherwise false. */
    public get isApp(): boolean {
        return this.kind == ElementKind.App;
    }

    /** True if this reference's {@link kind} is {@link ElementKind.Type}, otherwise false. */
    public get isType(): boolean {
        return this.kind == ElementKind.Type;
    }

    /** True if this reference's {@link kind} is {@link ElementKind.Field}, otherwise false. */
    public get isField(): boolean {
        return this.kind == ElementKind.Field;
    }

    /** True if this reference's {@link kind} is {@link ElementKind.Endpoint}, otherwise false. */
    public get isEndpoint(): boolean {
        return this.kind == ElementKind.Endpoint;
    }

    /** True if this reference's {@link kind} is {@link ElementKind.Statement}, otherwise false. */
    public get isStatement(): boolean {
        return this.kind == ElementKind.Statement;
    }

    /**
     * True if this reference's {@link kind} is {@link ElementKind.App}, {@link ElementKind.Type} or
     * {@link ElementKind.Field}, otherwise false.
     */
    public get isData(): boolean {
        return this.isApp || this.isType || this.isField;
    }

    /**
     * True if this reference's {@link kind} is {@link ElementKind.App}, {@link ElementKind.Endpoint} or
     * {@link ElementKind.Statement}, otherwise false.
     */
    public get isBehavior(): boolean {
        return this.isApp || this.isEndpoint || this.isStatement;
    }

    /** The references kind represented numerically. @see {@link ElementKind.toNumeric()} */
    public get numericKind(): number {
        return ElementKind.toNumeric(this.kind);
    }

    /**
     * Converts this element reference to a it's sysl representation, typically used for fields types or annotation
     * values that refer to other elements. This method doesn't support behavior types since they're not supported by
     * the Sysl syntax.
     * @param compact Optional. True to remove spaces around the namespace delimiter (`::`), otherwise keep them.
     * Defaults to false.
     * @param parentRef Optional. If specified and if it has the same app as the current instance, will produce a string
     * of a relative reference (omitting the app part), unless it's an app reference.
     * @returns The formatted string representing the element reference. Unsafe characters are escaped.
     * @throws `Error` if the current instance describes a behavior element (Endpoint or Statement).
     */
    toSysl(compact: boolean = false, parentRef?: ElementRef): string {
        if (!this.isData) throw new Error("Behavior element references cannot be represented in Sysl.");
        if (this === ElementRef.CurrentApp) return ".";
        const fullAppName = this.toAppPartsSafe().join(compact ? "::" : " :: ");
        const typeName = toSafeName(this.typeName);
        const fieldName = toSafeName(this.fieldName);
        let parts = [fullAppName, typeName, fieldName];
        if (!this.isApp && parentRef?.appsEqual(this)) parts.shift();
        return parts.filter(n => n).join(".");
    }

    /** Returns a string representation of the reference, with compact formatting. */
    toString(): string {
        if (this.isData) return this.toSysl(true);

        const fullAppName = this.toAppPartsSafe().join("::");
        const endpointName = `[${this.endpointName.replaceAll(/\]/g, "%5D")}]`;
        const indices = this.statementIndices.length ? `[${this.statementIndices.join(",")}]` : "";
        const parts = [fullAppName, endpointName, indices];
        return parts.filter(n => n).join(".");
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
        const parent = this.toParentOrUndefined();
        if (!parent) throw new Error("Cannot get parent reference for a root-level Application");
        return parent;
    }

    /**
     * Returns an reference to the parent of the current reference, or `undefined` if it's an application in
     * the root namespace. All namespace parts are considered intermediate applications, so the last part of a namespace
     * will be returned as the parent when an application reference is provided.
     * @returns An {@link ElementRef} to the parent/self of the current reference.
     */
    toParentOrUndefined(): ElementRef | undefined {
        if (this.isApp) return this.namespace.length ? this.popApp() : undefined;
        else if (this.statementIndices.length > 1)
            return this.with({ statementIndices: this.statementIndices.slice(0, -1) });

        return this.truncate(ElementKind.fromNumeric(this.numericKind - 1, this.isBehavior));
    }

    /**
     * Parses a string into an {@link ElementRef}. The string can refer to an app, type or field. The element isn't
     * guaranteed to exist in any given model. Decodes escaped characters.
     * @param refStr The element reference string to parse.
     * @returns An {@link ElementRef} that references the element specified in the supplied string.
     * @throws `Error` Thrown if the syntax of the supplied string is invalid.
     */
    static parse(refStr: string): ElementRef {
        refStr = refStr.trim();
        const firstDot = refStr.indexOf(".");
        const appNameParts = refStr
            .substring(0, firstDot == -1 ? undefined : firstDot)
            .trimEnd()
            .split(/\s*::\s*/);
        if (appNameParts.some(p => !isSafeName(p))) throw new Error("Disallowed characters in app name/namespace.");

        const namespace = appNameParts.slice(0, -1).map(fromSafeName);
        const appName = fromSafeName(appNameParts.at(-1)!);
        const endpointBracket = this.parseNextBracket(refStr, firstDot + 1);

        if (!endpointBracket) {
            let parts: string[] = [];
            if (firstDot != -1) parts = refStr.substring(firstDot + 1).split(/\s*\.\s*/);
            if (parts.length > 2) throw new Error("Too many dots.");
            if (parts.some(p => !isSafeName(p))) throw new Error("Disallowed characters in type/field name.");
            return new ElementRef(namespace, appName, fromSafeName(parts[0] ?? ""), fromSafeName(parts[1] ?? ""));
        } else {
            const endpointName = fromSafeName(endpointBracket.content);
            let indices: number[] = [];

            const indicesBracket = this.parseNextBracket(refStr, endpointBracket.end + 1);
            if (indicesBracket) {
                indices = indicesBracket.content.split(",").map(n => (n ? Number(n) : NaN));
                if (indices.some(i => isNaN(i) || i < 0)) throw new Error("Invalid statement indices.");
            } else if (refStr[endpointBracket.end] == ".") {
                throw new Error("Expected open brackets after dot.");
            }

            return new ElementRef(namespace, appName, "", "", endpointName, indices);
        }
    }

    /**
     * Tries to parse a string into an {@link ElementRef}. The string can refer to an app, type or field. The element
     * isn't guaranteed to exist in any given model. Decodes escaped characters.
     * @param refStr The element reference string to parse.
     * @returns An {@link ElementRef} that references the element specified in the supplied string, or `undefined` if
     * the the syntax of the supplied string is invalid.
     */
    static tryParse(refStr: string): ElementRef | undefined {
        try {
            return this.parse(refStr);
        } catch {
            return undefined;
        }
    }

    /** Parses the next bracketed part. Returns undefined if no bracket is found. `end` is exclusive. */
    private static parseNextBracket(refStr: string, start: number) {
        const open = refStr.indexOf("[", start);
        if (open == -1) return undefined;

        const close = refStr.indexOf("]", open + 1);
        if (close == -1) throw new Error(`Missing closing bracket.`);

        let end = refStr.indexOf(".", close + 1);
        if (end == -1) end = refStr.length;
        if (refStr.substring(start, open).trim() || refStr.substring(close + 1, end).trim())
            throw new Error("Unexpected characters outside brackets.");

        const content = refStr.substring(open + 1, close).trim();
        if (!content) throw new Error("No content found inside brackets.");

        return { content, end };
    }

    /**
     * Returns true if the supplied element reference refers to the same element as the current instance, otherwise
     * false.
     */
    equals(other: ElementRef): boolean {
        if (this.isData) return this.typesEqual(other) && this.fieldName == other.fieldName;
        return this.endpointsEqual(other) && ElementRef.arraysEqual(this.statementIndices, other.statementIndices);
    }

    /**
     * Returns true if the supplied element reference refers to the same app and type as the current instance
     * (ignoring the field), otherwise false. If both references are to applications, the types are ignored.
     */
    typesEqual(other: ElementRef): boolean {
        return this.appsEqual(other) && this.isData && other.isData && this.typeName == other.typeName;
    }

    /**
     * Returns true if the supplied element reference refers to the same app and type as the current instance
     * (ignoring the field), otherwise false. If both references are to applications, the types are ignored.
     */
    endpointsEqual(other: ElementRef): boolean {
        return this.appsEqual(other) && this.isBehavior && other.isBehavior && this.endpointName == other.endpointName;
    }

    /**
     * Returns true if the supplied element reference refers to the same app (ignoring the type or field), otherwise
     * false.
     */
    appsEqual(other: ElementRef): boolean {
        return (
            this.namespace.length == other.namespace.length &&
            this.namespace.every((_, i) => this.namespace[i] == other.namespace[i]) &&
            ElementRef.arraysEqual(this.namespace, other.namespace) &&
            this.appName == other.appName
        );
    }

    private static arraysEqual(first: readonly any[], second: readonly any[]): boolean {
        return first.length == second.length && first.every((_, i) => first[i] == second[i]);
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
            if (ancestor.isData != this.isData) return false;
            if (this.isData) return isDeeper && this.typesEqual(ancestor);
            if (ancestor.isStatement && this.isStatement)
                return (
                    ancestor.statementIndices.length < this.statementIndices.length &&
                    ancestor.statementIndices.every((p, i) => this.statementIndices[i] == p)
                );
            return isDeeper && this.endpointsEqual(ancestor);
        }
    }

    /**
     * Clones the current instance with additional modifications to its parts.
     * @param parts An objects with optional properties to modify parts of the new element reference. If any property
     * is not specified (or specified with `undefined`), it will not be modified. If an empty string/array is specified,
     * it will remove that part of the reference.
     * @returns An {@link ElementRef} combining the current instances and the part changes requested.
     */
    with(parts: {
        namespace?: string[];
        appName?: string;
        typeName?: string;
        fieldName?: string;
        endpointName?: string;
        statementIndices?: readonly number[];
    }): ElementRef {
        return new ElementRef(
            parts.namespace ?? this.namespace,
            parts.appName ?? this.appName,
            parts.typeName ?? this.typeName,
            parts.fieldName ?? this.fieldName,
            parts.endpointName ?? this.endpointName,
            parts.statementIndices ?? this.statementIndices
        );
    }

    /**
     * Truncates an ElementKind to a certain {@link ElementKind} by removing the parts that aren't in that kind. If
     * a truncation is requested to a kind that contains the same part or to one that has more parts than the current
     * reference has, the same reference is returned.
     * @param kind The {@link ElementKind} to truncate the current reference into.
     * @returns A truncated reference, or the same reference if there is nothing to truncate.
     * @throws `Error` if an invalid kind is specified.
     */
    truncate(kind: ElementKind): ElementRef {
        switch (kind) {
            case ElementKind.App:
                return new ElementRef(this.namespace, this.appName);
            case ElementKind.Type:
                if (!this.isData) throw new Error("Cannot truncate non-Data ref into type ref.");
                return new ElementRef(this.namespace, this.appName, this.typeName);
            case ElementKind.Endpoint:
                if (!this.isBehavior) throw new Error("Cannot truncate non-behavior ref into endpoint ref.");
                return new ElementRef(this.namespace, this.appName, "", "", this.endpointName);
            case ElementKind.Field:
                if (!this.isData) throw new Error("Cannot truncate non-data ref into field ref.");
                return this;
            case ElementKind.Statement:
                if (!this.isBehavior) throw new Error("Cannot non-behavior data ref into statement ref.");
                return this;
            default:
                throw new Error("Invalid kind specified.");
        }
    }

    /**
     * Returns a new {@link ElementRef} where the app is set to the provided argument and the existing app name is
     * added to the namespace.
     */
    pushApp(appName: string): ElementRef {
        return this.with({ namespace: [...this.namespace, this.appName], appName });
    }

    /**
     * Return a new {@link ElementRef} where the last part of the namespace is moved to the app name. Throws if there is
     * no namespace.
     */
    popApp(): ElementRef {
        if (!this.namespace.length) throw new Error("Cannot pop app from reference with no namespace.");
        return this.with({ namespace: this.namespace.slice(0, -1), appName: this.namespace.at(-1) });
    }

    /**
     * Returns the current instance, since ElementRef is immutable and doesn't need cloning.
     * @param _context Unused.
     * @returns The current instance.
     */
    clone(_context?: CloneContext): ElementRef {
        return this;
    }

    /**
     * Deep clones the provided object, substituting all instances of {@link ElementRef} with their
     * {@link ElementRef.toString()} representation. Mainly used for debugging output and test assertions.
     */
    static stringSubstitute(obj: any): any {
        if (!obj || typeof obj != "object") return obj;
        if (obj instanceof ElementRef) return obj.toString();
        if (Array.isArray(obj)) return obj.map(x => this.stringSubstitute(x));
        return Object.keys(obj)
            .filter(k => obj.hasOwnProperty(k))
            .reduce((o, k) => ({ ...o, [k]: this.stringSubstitute(obj[k]) }), {});
    }
}
