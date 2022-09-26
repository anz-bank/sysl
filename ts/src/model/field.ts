import { safeName } from "../common/format";
import { Element, IElementParams, setParentAndModelDeep } from "./element";
import { addTags, renderAnnos } from "./renderers";
import { Primitive } from "./primitive";
import { CollectionDecorator } from "./decorator";
import { ElementRef } from "./common";

export type FieldValue = Primitive |  ElementRef | CollectionDecorator;

export class Field extends Element {

    constructor(name: string, public value: FieldValue, public optional: boolean = false, p?: IElementParams) {
        super(name, p?.locations ?? [], p?.annos ?? [], p?.tags ?? [], p?.model, p?.parent);
        setParentAndModelDeep(this, this.annos, this.tags);
    }

    override toSysl(): string {
        const optStr = this.optional ? "?" : "";
        let sysl = `${safeName(this.name)}${this.value.toSysl()}${optStr}`;
        sysl = addTags(sysl, this.tags);
        if (this.annos.length) {
            sysl += `:\n${renderAnnos(this.annos)}`;
        }
        return sysl;
    }
}
