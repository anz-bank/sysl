import { Annotation, Tag } from "./attribute";

function renderSections(sections: string[], delimiter: string): string {
    return sections.where(s => s.length > 0).toArray().join(delimiter);
}

export function renderInlineSections(sections: string[]): string {
    return renderSections(sections, ' ');
}

function renderMultilineSections(sections: string[]): string {
    return renderSections(sections, '\n\n');
}

export function renderTags(tags: Tag[]): string {
    return tags.any() ? ` [${tags.select(t => t.toSysl()).toArray().join(", ")}]` : '';

}

export function renderAnnos(annos: Annotation[]): string {
    return annos.select(a => `@${a.toSysl()}`).toArray().join("\n");
}
