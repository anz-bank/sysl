import { joinedAppName, safeName } from "../common/format";
import { Location } from "../common/location";
import { Element } from "./element";
import { Model } from "./model";

/** An element that can be serialized to Sysl source. */
export interface IRenderable {
    toSysl(): string;
}

/**
 * An element that has source context locating it in one or more sourcefiles.
 *
 * All objects with locations in Sysl source exist in a Sysl model. The {@code model} property is
 * a reference to that model.
 *
 * {@code model} may be undefined if an object has been created detached from a model.
 */
export interface ILocational {
    locations: Location[];
    model?: Model;
}

/**
 * An object that has a parent in the Sysl model.
 *
 * Follow the chain of {@code parent} properties should lead to an {@link Application} which will
 * have a falsey parent.
 *
 * {@code parent} can also be undefined if an object has been created detached from a model.
 */
export interface IChild extends ILocational {
    parent?: Element;
}

export class ElementRef implements IRenderable {
    readonly kind: ElementKind;
    
    constructor(public readonly namespace: string[], public readonly appName: string, public readonly typeName: string = "", public readonly fieldName: string = "") {
        if (fieldName && !typeName)
            throw new Error("Cannot specify fieldName but omit typeName");
        
        if (fieldName)
            this.kind = ElementKind.Field;
        else if (typeName)
            this.kind = ElementKind.Type;
        else
            this.kind = ElementKind.App;
    }

    toSysl(compact: boolean = false): string {
        const fullAppName = [...this.namespace, this.appName].map(safeName).join(compact ? "::" : " :: ");
        const typeName = safeName(this.typeName);
        const fieldName = safeName(this.fieldName);

        return [fullAppName, typeName, fieldName].filter(n => n).join(".");
    }

    static parse(refStr: string): ElementRef {
        const parts = refStr.split(".", 3);

        if (!parts[0])
            throw new Error(`Invalid string element reference: ${refStr}`)

        const appNameParts: string[] = parts[0].split(/\s*::\s*/);
        return new ElementRef(appNameParts.slice(0, -1), appNameParts.at(-1)!, parts[1] ?? undefined, parts[2] ?? undefined);
    }

    equals(other: ElementRef): boolean {
        return this.namespace.length == other.namespace.length &&
            this.namespace.every((_, i) => this.namespace[i] == other.namespace[i]) &&
            this.appName == other.appName &&
            this.typeName == other.typeName &&
            this.fieldName == other.fieldName;
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
