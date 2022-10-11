import "reflect-metadata";
import { indent, safeName } from "../common/format";
import {
    IElementParams,
    setParentAndModelDeep,
    ParentElement,
} from "./element";
import { Field } from "./field";
import { addTags, renderAnnos } from "./renderers";
import { Struct } from "./struct";

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

    toSysl(): string {
        const discriminator = this.isTable ? "!table" : "!type";
        let sysl = `${addTags(
            `${discriminator} ${safeName(this.name)}`,
            this.tags
        )}:`;
        if (this.annos.length) {
            sysl += `\n${indent(renderAnnos(this.annos))}`;
        }
        return (sysl += `\n${new Struct(this.children).toSysl()}`);
    }
}
