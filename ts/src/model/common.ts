import { Location } from "../common/location";
import { Element } from "./element";
import { Model } from "./model";

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
