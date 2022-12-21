const defaultIndent = "    ";

export function indent(text: string): string {
    return text
        .split("\n")
        .map(l => `${defaultIndent}${l}`)
        .join("\n");
}

export function joinedAppName(name: string[], compact: boolean = false): string {
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
export function toSafeName(name: string): string {
    const percentEncode = (m: string) => `%${m.charCodeAt(0).toString(16).toUpperCase()}`;
    name = name.replaceAll(/([^-a-zA-Z0-9_])/g, percentEncode);
    if (keywords.includes(name)) {
        name += "_";
    }
    return name;
}

/** Unescapes characters in Sysl names that are unsafe to use in Sysl names. */
export function fromSafeName(name: string): string {
    return name.trim().replaceAll(/%([0-9A-Fa-f]{2})/g, m => String.fromCharCode(parseInt(m.slice(1), 16)));
}

/**
 * Realigns text to be correctly indented: Removes all common indentation so the least indented line starts at col 0,
 * detects indentation unit (number of spaces) and resizes further indentations according to the indentation unit size
 * requested (default to 4). Assumes the first non-whitespace indented line (after common indentation removal) is
 * indented by exactly one unit.
 * @param str The test to realign.
 * @param indentSize Optional. The desired indentation unit size. The default is 4.
 * @returns The realigned text.
 */
export function realign(str: string, indentSize = 4) {
    // Remove the first newline to flatten final output, split on the remaining new lines.
    const lines = str.replace("\n", "").split(/\r?\n/);
    const minIndent = Math.min(...lines.filter(l => !!l.trim()).map(l => l.search(/\S/)));

    // First non-whitespace indented line in a Sysl file should have only one indent unit.
    const detectedIndent = (lines.find(l => l.search(/\S/) > minIndent)?.search(/\S/) ?? 8) - minIndent;
    const getIndentUnits = (l: string) => Math.max(l.search(/\S|$/) - minIndent, 0) / detectedIndent;
    return lines.map(l => " ".repeat(getIndentUnits(l) * indentSize) + l.trimStart()).join("\n");
}
