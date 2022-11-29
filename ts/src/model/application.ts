import { indent } from "../common/format";
import { ElementRef } from "./common";
import { Element, IElementParams, ParentElement, setParentAndModelDeep } from "./element";
import { addTags, renderAnnos } from "./renderers";
import { Endpoint } from "./statement";
import { Type } from "./type";

export class Application extends ParentElement<Element> {
    namespace: string[];
    endpoints: Endpoint[];
    children: Element[];

    constructor({ namespace, name, endpoints, children: types, locations, annos, tags, model }: ApplicationParams) {
        if (!name) throw new Error("'name' must be specified, and optionally 'namespace'.");

        super(name, locations ?? [], annos ?? [], tags ?? [], model);
        this.namespace = namespace ?? [];
        this.endpoints = endpoints ?? [];
        this.children = types ?? [];

        setParentAndModelDeep(this, this.endpoints, this.children, this.annos, this.tags);
    }

    public get types(): Type[] {
        return this.children.filter(c => c instanceof Type) as Type[];
    }

    toSysl(): string {
        let sysl = `${addTags(this.toString(), this.tags)}:`;
        if (this.annos.length) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        if (this.endpoints.length) {
            sysl += `\n${this.endpoints
                .filter(e => !e.isPubsub)
                .map(e => indent(e.toSysl()))
                .join("\n\n")}`;
        }
        if (this.children.length) {
            sysl += `\n${this.children.map(t => indent(t.toSysl())).join("\n\n")}`;
        }
        if (!this.annos.length && !this.endpoints.length && !this.children.length) {
            sysl += `\n${indent("...")}`;
        }
        return sysl;
    }

    toRef(): ElementRef {
        return new ElementRef(this.namespace, this.name);
    }

    override toString(): string {
        return this.toRef().toSysl();
    }
}

export type ApplicationParams = IElementParams & {
    namespace?: string[];
    name?: string;
    endpoints?: Endpoint[];
    children?: Element[];
};
