import "reflect-metadata";
import { jsonArrayMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../location";
import { Annotation, AnnoValue, Tag } from "../model/attribute";

@jsonObject
export class PbAttributeArray {
    @jsonArrayMember(() => PbAttribute) elt!: PbAttribute[];
}

@jsonObject
export class PbAttribute {
    @jsonMember s?: string
    @jsonMember n?: number
    @jsonArrayMember(Location) sourceContexts!: Location[];
    @jsonMember(PbAttributeArray) a?: PbAttributeArray;

    getAnnoValue(): AnnoValue {
        if (this.s) return this.s;
        else if (this.n) return this.n;
        else if (this.a) return this.a.elt.select(i => i.getAnnoValue()).toArray();
        else throw new Error("Missing attribute value")
    }

    toTag(): Tag {
        return new Tag(this.getAnnoValue(), this.sourceContexts ?? [])
    }

    toAnno(name: string): Annotation {
        return new Annotation(name, this.sourceContexts ?? [], this.getAnnoValue());
    }
}
