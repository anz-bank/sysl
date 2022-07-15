import { Annotation, Tag } from "./attribute";

function renderSections(sections: string[], delimiter: string): string {
    return sections.filter(s => s.length > 0).join(delimiter);
}

export function renderInlineSections(sections: string[]): string {
    return renderSections(sections, " ");
}

export function addTags(existing: string, tags: Tag[]): string {
    return tags.length
        ? `${existing} [${tags.map(t => t.toSysl()).join(", ")}]`
        : existing;
}

export function renderAnnos(annos: Annotation[]): string {
    return annos.map(a => `@${a.toSysl()}`).join("\n");
}
