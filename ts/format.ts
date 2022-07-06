export function indent(text: string): string {
    return `${text
        .split("\n")
        .select(l => `    ${l}`)
        .toArray()
        .join("\n")}`;
}

export function joinedAppName(
    name: string[],
    compact: boolean = false
): string {
    return name.join(compact ? "::" : " :: ");
}

export function realign(str: string) {
    // Remove the first newline to flatten final output, split on the remaining new lines
    const lines = str.replace("\n", "").split(/\r?\n/);
    const minIndent = lines
        .where(l => !!l.trim())
        .min(l => l.takeWhile(c => c == " ").count());
    return lines
        .select(l => l.substring(minIndent))
        .toArray()
        .join("\n");
}
