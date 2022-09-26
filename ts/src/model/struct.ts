import { indent, safeName } from "../common/format";
import { IRenderable } from "./common";
import { Field } from "./field";
import { addTags, renderAnnos } from "./renderers";

// The type of `!type`, `!table` and endpoints.
export class Struct implements IRenderable {
    fields: Field[];

    constructor(elements: Field[]) {
        this.fields = elements ?? [];
    }

    toSysl(): string {
        return indent(
            this.fields
                .map(e => {
                    let sysl = addTags(
                        `${safeName(e.name)} <: ${e.value.toSysl()}${
                            e.optional ? "?" : ""
                        }`,
                        e.tags
                    );
                    if (e.annos.length) {
                        sysl += `:\n${indent(renderAnnos(e.annos))}`;
                    }
                    return sysl;
                })
                .join("\n")
        );
    }
}
