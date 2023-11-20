// This file is separate from ./type.ts to break a circular type reference.
// PbTypeDefList is imported to type.ts with an `import type`.

import { jsonObject, jsonMember } from "typedjson";
import { PbTypeDef } from "./type";
import { ElementRef, FieldValue } from "../model";

@jsonObject
export class PbTypeDefList {
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toValue(parentRef: ElementRef | undefined): FieldValue {
        return this.type.toValue(parentRef);
    }
}
