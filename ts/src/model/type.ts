import "reflect-metadata";
import { indent, toSafeName } from "../common/format";
import { ElementRef } from "./common";
import {
    IElementParams,
    setParentAndModelDeep,
    ParentElement,
} from "./element";
import { Field } from "./field";
import { addTags, renderAnnos } from "./renderers";

export class Type extends ParentElement<Field> {
    constructor(
        name: string,
        public isTable: boolean = false,
        public children: Field[] = [],
        p?: IElementParams
    ) {
        super(
            name,
            p?.locations ?? [],
            p?.annos ?? [],
            p?.tags ?? [],
            p?.model,
            p?.parent
        );
        setParentAndModelDeep(this, this.children, this.annos, this.tags);
    }

    toRef(): ElementRef {
        return this.parent!.toRef().with({typeName: this.name});
    }

    toSysl(): string {
        const discriminator = this.isTable ? "!table" : "!type";

        let sysl = `${addTags(
            `${discriminator} ${toSafeName(this.name)}`,
            this.tags
        )}:`;

        if (this.annos.length) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }

        if (this.children.length) {
            // sysl += `\n${new Struct(this.children).toSysl()}`;
            sysl += `\n${indent(this.children.map(f => f.toSysl()).join("\n"))}`;
        }

        if (!this.annos.length && !this.children.length) {
            sysl += `\n${indent("...")}`;
        }

        return sysl;
    }
}
