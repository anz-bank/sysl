import { toSafeName, indent } from "../common/format";
import { ElementRef } from "./common";
import { CollectionDecorator } from "./decorator";
import { Element, IElementParams } from "./element";
import { Enum } from "./enum";
import { Primitive } from "./primitive";
import { addTags, renderAnnos } from "./renderers";
import { Struct } from "./struct";
import { Union } from "./union";

/** Represents anything that doesn't have it's own special class yet. */
// TODO: Put everything in it's own class and remove this one.
export class GenericElement extends Element {
    discriminator: string;
    value: GenericValue;

    constructor({ discriminator, name, value, annos, tags, locations, parent, model }: GenericElementParams) {
        if (!discriminator) throw new Error("You must specify a discriminator");

        super(name, locations ?? [], annos ?? [], tags ?? [], model, parent);
        this.value = value;
        this.discriminator = discriminator;
        this.attachSubitems();
    }

    toRef(): ElementRef {
        throw new Error("Method not implemented.");
    }

    toSysl(): string {
        let sysl = `${addTags(`${this.discriminator} ${toSafeName(this.name)}`, this.tags)}:`;
        if (this.annos.length) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }

        // TODO: Make Alias class be independent of GenericElement to remove this special-casing
        if (this.discriminator == "!alias") return (sysl += `\n${indent(this.value.toSysl())}`);

        return (sysl += `\n${this.value.toSysl()}`);
    }
}

export type GenericElementParams = IElementParams & {
    discriminator?: string;
    name: string;
    value: GenericValue;
    optional?: boolean;
};

export type GenericValue = Primitive | Struct | CollectionDecorator | Enum | Union | ElementRef;
