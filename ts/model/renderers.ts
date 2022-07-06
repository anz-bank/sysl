import { Annotation, Tag } from "./attribute";

function renderSections(sections: string[], delimiter: string): string {
    return sections
        .where(s => s.length > 0)
        .toArray()
        .join(delimiter);
}

export function renderInlineSections(sections: string[]): string {
    return renderSections(sections, " ");
}

export function addTags(existing: string, tags: Tag[]): string {
    return tags.any()
        ? `${existing} [${tags
              .select(t => t.toSysl())
              .toArray()
              .join(", ")}]`
        : existing;
}

export function renderAnnos(annos: Annotation[]): string {
    return annos
        .select(a => `@${a.toSysl()}`)
        .toArray()
        .join("\n");
}
