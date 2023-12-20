import path from "path";
import "reflect-metadata";
import { jsonMember, jsonObject } from "typedjson";
import { ILocational, Model } from "../model";

/** Describes the offset location within a file by zero-based line and column. */
@jsonObject
export class Offset {
    /** The zero-based line number of the {@link Offset}. */
    @jsonMember line: number;
    /** The zero-based column number of the {@link Offset}. */
    @jsonMember col: number;

    /**
     * Create a new {@link Offset} object with the specified zero-based line and column.
     * @param line Optional. The zero-based line number of the {@link Offset}. If not specified, defaults to `0`.
     * @param col Optional. The zero-based column number of the {@link Offset}. If not specified, defaults to `0`.
     */
    constructor(line: number = 0, col: number = 0) {
        if (isNaN(line) || isNaN(col) || line < 0 || col < 0)
            throw new Error(`NaN or negative values are not allowed for line and column.`);
        this.line = line;
        this.col = col;
    }
}

/** Describes the source location of a Sysl document fragment. */
@jsonObject
export class Location {
    /** The Sysl path or URL of the file where the fragment is located, relative to {@link Model.syslRoot}. */
    @jsonMember file!: string;
    /** The {@link Offset} inside the {@link file} where the fragment begins. */
    @jsonMember start!: Offset;
    /** The {@link Offset} inside the {@link file} where the fragment ends. */
    @jsonMember end!: Offset;

    /**
     * Creates a new {@link Location} object with the specified file and offset range.
     * @param file The file name where the fragment is located.
     * @param start The {@link Offset} inside the {@link file} where the fragment begins.
     * @param end The {@link Offset} inside the {@link file} where the fragment ends.
     * @throw `Error` when the file name is not specified, or the offsets are invalid.
     */
    constructor(file: string, start: Offset, end: Offset) {
        if (file == undefined) return; // Parameterless constructor used by TypedJson for deserialization.

        if ((end.line == 0 || end.line == start.line) && end.col == 0)
            end = new Offset(end.line || start.line, start.col);

        if (start.line > end.line) throw new Error("End line must be greater than or equal to start line.");
        if (start.line == end.line && start.col > end.col) {
            console.dir({ start, end });
            throw new Error("When on the same line, end column must be greater than or equal to start column.");
        }

        this.file = file;
        this.start = start;
        this.end = end;
    }

    /**
     * Gets the file extension part of {@link file} path. This is the part of the file name beginning at the first
     * period.
     * @example
     * If the {@link file} is
     * ```//github.com/someone/myproject/program.edg.sysl@master```
     * the returned value will be `.edg.sysl`
     */
    public get extension(): string {
        return "." + path.basename(this.file).split("@")[0].split(".").slice(1).join(".");
    }

    /**
     * Gets the git reference (branch, tag of hash) of the {@link file} path. This is the part after the '@' symbol,
     * or an empty string if none is specified.
     */
    public get gitRef(): string {
        return path.basename(this.file).split("@")[1] ?? "";
    }

    /**
     * True if the {@link file} path point to a remote file (starts with `//`), otherwise false.
     */
    public get isRemote(): boolean {
        return this.file.startsWith("//");
    }

    /**
     * Formats the {@link Location} object as a string with five parts joined by the colon character:
     * - File name
     * - 1-based start line
     * - 1-based start column (omitted if 0 and end line/column are omitted)
     * - 1-based end line (omitted if identical to start line)
     * - 1-based end column (omitted if not greater than start column and end line is omitted)
     * This format is backwards-compatible with the file linking in VS Code (for file, and start line/col).
     * @param [startLineOnly=false] Optional. If true, only the file name and start line are included in the string.
     * @returns A string representation of the location. Examples of strings returned:
     * - `test.sysl:1` - Starts at the first line, no columns specified. Assumed to mean the entire line.
     * - `test.sysl:1:5` - Starts at the first line at column 5. Assumed to mean until the end of the line.
     * - `test.sysl:1:5::10` - Columns 5 to 10 of the first line. End line is omitted when identical to start line.
     * - `test.sysl:1:5:2:10` - From column 5 of the first line to column 10 of the second line.
     */
    public toString(startLineOnly: boolean = false): string {
        const parts = [this.file, this.start.line + 1];

        if (!startLineOnly) {
            const hasEndLine = this.end.line > this.start.line;
            const hasEnd = hasEndLine || this.end.col > this.start.col;
            if (this.start.col > 0 || hasEnd) {
                parts.push(this.start.col + 1);
                if (hasEnd) {
                    parts.push(hasEndLine ? this.end.line + 1 : "");
                    parts.push(this.end.col + 1);
                }
            }
        }

        return parts.join(":");
    }

    /**
     * Parses a location string into a {@link Location} object. See {@link Location.toString} for the format.
     * @param locStr The location string to parse.
     * @returns A {@link Location} object corresponding to the location string.
     * @throws `Error` if the location string is invalid.
     */
    public static parse(locStr: string): Location {
        const parts = locStr.split(":");

        if (parts.length != 2 && parts.length != 3 && parts.length != 5)
            throw new Error(`Invalid number of parts in location string (must be 2, 3 or 5): ${locStr}`);

        const file = parts[0].trim();
        const startLine = Number(parts[1]) - 1;
        const startCol = Number(parts[2] || 1) - 1;
        const endLine = Number(parts[3] || parts[1]) - 1;
        const endCol = Number(parts[4] || parts[2] || 1) - 1;

        if (!file) throw new Error(`File name is missing in location string: ${locStr}`);

        return new Location(file, new Offset(startLine, startCol), new Offset(endLine, endCol));
    }

    /**
     * Creates a deep copy of the {@link Location} object.
     * @returns A new instance of a {@link Location} with the equal values to this instance.
     */
    public clone(): Location {
        return new Location(
            this.file,
            new Offset(this.start.line, this.start.col),
            new Offset(this.end.line, this.end.col)
        );
    }

    /**
     * Compares two {@link Location} objects.
     * @param a The first {@link Location} to compare.
     * @param b The second {@link Location} to compare.
     * @returns A negative number if `a` is before `b`, a positive number if `a` is after `b`, or 0 if they are equal.
     */
    public static compare(a: Location, b: Location): number {
        return (
            a?.file.localeCompare(b?.file) ||
            a?.start.line - b?.start.line ||
            a?.start.col - b?.start.col ||
            a?.end.line - b?.end.line ||
            a?.end.col - b?.end.col
        );
    }

    /**
     * Compares the first location of two {@link ILocational} objects.
     * @param a The first {@link ILocational} to compare.
     * @param b The second {@link ILocational} to compare.
     * @returns A negative number if `a` is before `b`, a positive number if `a` is after `b`, or 0 if they are equal.
     */
    public static compareFirst(a: ILocational, b: ILocational): number {
        return Location.compare(a.locations[0], b.locations[0]);
    }
}
