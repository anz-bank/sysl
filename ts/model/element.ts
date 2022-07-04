import { indent } from "../util";
import { Annotation, Tag } from "./attribute";
import { ComplexType } from "./common";
import { renderAnnos, renderTags } from "./renderers";
import { Type } from "./type";

export class Element<T extends ComplexType> {
    tags: Tag[];
    annos: Annotation[];
    content: T

    constructor(tags: Tag[], annos: Annotation[], type: T) {
        this.tags = tags;
        this.annos = annos;
        this.content = type;
    }

    toSysl(): string {
        if (this.content instanceof Type && this.content.discriminator == '') {
            let sysl = `${this.content.name} <: ${this.content.toSysl()}${renderTags(this.tags)}`;
            if (this.annos.any()) {
                sysl += `:\n${indent(renderAnnos(this.annos))}`;
            }
            return sysl;
        }
        let sysl = `${this.content.discriminator.length > 0 ? `${this.content.discriminator} ` : ''}${this.content.name}${renderTags(this.tags)}:`;
        if (this.annos.any()) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        sysl += `\n${this.content.toSysl()}`;
        return sysl;
    }
}
