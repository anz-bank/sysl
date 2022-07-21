import "reflect-metadata";
import { jsonArrayMember, jsonMember, jsonObject } from "typedjson";
import { Location } from "../common/location";
import { Annotation, AnnoValue, Tag } from "../model/attribute";

@jsonObject
export class PbAttributeArray {
    @jsonArrayMember(() => PbAttribute) elt!: PbAttribute[];
}

@jsonObject
export class PbAttribute {
    @jsonMember s?: string;
    @jsonMember n?: number;
    @jsonMember(PbAttributeArray) a?: PbAttributeArray;
    @jsonArrayMember(Location) sourceContexts!: Location[];

    getAnnoValue(): AnnoValue {
        if (this.s != null) {
            return this.s;
        } else if (this.n != null) {
            return this.n;
        } else if (this.a != null) {
            return this.a.elt?.map(i => i.getAnnoValue()) ?? [];
        } else {
            throw new Error("Missing attribute value: " + JSON.stringify(this));
        }
    }

    toTag(): Tag {
        return new Tag({
            value: this.getAnnoValue(),
            locations: this.sourceContexts ?? [],
        });
    }

    toAnno(name: string): Annotation {
        return new Annotation({
            name,
            locations: this.sourceContexts ?? [],
            value: this.getAnnoValue(),
        });
    }
}
