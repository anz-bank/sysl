// This file is separate from ./type.ts to break a circular type reference.
// PbTypeDefList is imported to type.ts with an `import type`.

import { jsonObject, jsonMember } from "typedjson";
import { PbTypeDef } from "./type";
import { Type } from "../model";

@jsonObject
export class PbTypeDefList {
    @jsonMember(() => PbTypeDef) type!: PbTypeDef;

    toModel(): Type {
        return this.type.toModel();
    }
}
