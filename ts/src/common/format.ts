const defaultIndent = "    ";

export function indent(text: string): string {
    return text
        .split("\n")
        .map(l => `${l ? defaultIndent : ""}${l}`) // Don't add indent to empty lines
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

/**
 * Converts a given identifier to a safe name which can be used in a Sysl file. Any characters that is disallowed in
 * safe names will be escaped. Additionally, if the resulting name is still disallowed due to it being a reserved
 * keyword or due to it starting with a digit, the first character will be escaped as well.
 * @param name A string containing the identifier to convert.
 * @returns A string containing a safe name.
 */
export function toSafeName(name: string): string {
    const encode = (m: string) => `%${m.charCodeAt(0).toString(16).toUpperCase()}`;
    name = name.replaceAll(/(^[0-9])|([^-a-zA-Z0-9_])/g, encode);
    return keywords.includes(name.toLowerCase()) ? encode(name[0]) + name.slice(1) : name;
}

/**
 * Converts a given safe name to it's original identifier by decoding escape sequences in the given string.
 * @param name A string containing the safe name to convert.
 * @returns A string containing the identifier.
 */
export function fromSafeName(name: string): string {
    return decodeURIComponent(name);
}

/**
 * Determine if a given string is a safe name which can be directly used in a Sysl file.
 * @param name The string to test for safety.
 * @returns True if the string is a safe name, otherwise false.
 */
export function isSafeName(name: string): boolean {
    return /^([-A-Za-z_]|%[0-9A-Fa-f]{2})([-A-Za-z0-9_]|%[0-9A-Fa-f]{2})*$/.test(name) && !keywords.includes(name);
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
