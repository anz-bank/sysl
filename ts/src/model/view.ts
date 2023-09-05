import { Application } from "./application";
import { Element } from "./element";
import { Type } from "./type";
import { Field } from "./field";
import { Endpoint, Statement } from "./statement";
import { Lazy, flatMapDeep } from "../common/util";

/**
 * Provides easy access to categories of elements regardless of their position in the element hierarchy. Every property
 * is lazily evaluated and computed when first accessed, then cached. Once cached, data will not be re-evaluated, even
 * if the underlying model has been changed. If modifications are made to the original model, create a new instance of
 * the view to access the latest data. The returned arrays cannot be modified since they are a read-only view of a
 * model, but each {@link Element} can be read and written to, and writes will update the {@link Element} in the
 * original model.
 */
//prettier-ignore
export class FlatView {
    #types = new Lazy(() => this.apps.flatMap((a) => a.types));
    #fields = new Lazy(() => this.types.flatMap((t) => t.children));
    #endpoints = new Lazy(() => this.apps.flatMap((a) => a.endpoints));
    #statements = new Lazy(() => flatMapDeep(this.endpoints.flatMap((e) => e.statements), (e) => e.children));
    #dataElements = new Lazy(() => [...this.apps, ...this.types, ...this.fields]);
    #behaviorElements = new Lazy(() => [...this.apps, ...this.endpoints, ...this.statements]);
    #allElements = new Lazy(() => [...this.apps, ...this.types, ...this.fields, ...this.endpoints, ...this.statements]);

    /** Creates a new instance of {@link FlatView} from the provided array of {@link Application}s. */
    constructor(public apps: readonly Application[]) { }

    /** Lazily evaluated. Returns all types, from all applications. */
    get types(): readonly Type[] { return this.#types.value; }
    /** Lazily evaluated. Returns all fields, from all types in all applications. */
    get fields(): readonly Field[] { return this.#fields.value; }
    /** Lazily evaluated. Returns all endpoints, from all applications. */
    get endpoints(): readonly Endpoint[] { return this.#endpoints.value; }
    /** Lazily evaluated. Returns all statements, from all endpoints in all applications. */
    get statements(): readonly Statement[] { return this.#statements.value; }
    /** Lazily evaluated. Returns all apps, types and fields. */
    get dataElements(): readonly Element[] { return this.#dataElements.value; }
    /** Lazily evaluated. Returns all apps, endpoints and statements. */
    get behaviorElements(): readonly Element[] { return this.#behaviorElements.value; }
    /** Lazily evaluated. Returns all elements. */
    get allElements(): readonly Element[] { return this.#allElements.value; }
}
