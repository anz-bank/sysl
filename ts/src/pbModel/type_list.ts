// This file is separate from ./type.ts to break a circular type reference.
// PbTypeDefList is imported to type.ts with an `import type`.

import { jsonObject, jsonMember } from "typedjson";
import { PbTypeDef } from "./type";
import { Element } from "../model";

@jsonObject
export class PbTypeDefList {
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toModel(): Element {
        return this.type.toModel();
    }
}
