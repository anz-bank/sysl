const defaultIndent = "    ";

export function indent(text: string): string {
    return text
        .split("\n")
        .map(l => `${defaultIndent}${l}`)
        .join("\n");
}

export function joinedAppName(
    name: string[],
    compact: boolean = false
): string {
    return name.join(compact ? "::" : " :: ");
}

/** Sysl reserved keywords that cannot be used in names. */
const keywords = [
    "int",
    "int32",
    "int64",
    "float",
    "float32",
    "float64",
    "decimal",
    "string",
    "date",
    "datetime",
    "bool",
    "bytes",
    "any",
];

/** Escapes characters that are unsafe to use in Sysl names. */
export function safeName(name: string): string {
    const percentEncode = (m: string) => `%${m.charCodeAt(0).toString(16)}`;
    name = name
        .replaceAll(/[/\\{} ]+/g, "_")
        .replaceAll(/([^-a-zA-Z0-9_])/g, percentEncode);
    if (keywords.includes(name)) {
        name += "_";
    }
    return name;
}

/** Unescapes characters in Sysl names that are unsafe to use in Sysl names. */
export function unescapeName(name: string): string {
    return name.replaceAll(/%([0-9A-Fa-f]{2})/g, m =>
        String.fromCharCode(parseInt(m.slice(1), 16))
    );
}

export function realign(str: string) {
    // Remove the first newline to flatten final output, split on the remaining new lines.
    const lines = str.replace("\n", "").split(/\r?\n/);
    const minIndent = Math.min(
        ...lines.filter(l => !!l.trim()).map((l: string) => l.search(/[^ ]/))
    );
    return lines.map(l => l.substring(minIndent)).join("\n");
}
